package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
	"github.com/matheusgomes28/urchin/common"
	"github.com/rs/zerolog/log"
)

// Estructura para la respuesta de Nominatim
type NominatimResponse struct {
	PlaceID     int64  `json:"place_id"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Address     struct {
		Historic     string `json:"historic"`
		Road         string `json:"road"`
		Town         string `json:"town"`
		Municipality string `json:"municipality"`
		County       string `json:"county"`
		State        string `json:"state"`
		Country      string `json:"country"`
		CountryCode  string `json:"country_code"`
		Postcode     string `json:"postcode"`
	} `json:"address"`
}

// Estructura para el resultado final
type PhotoMetadata struct {
	Filename  string   `json:"filename"`
	FilenameS string   `json:"filename_small"`
	FilenameM string   `json:"filename_medium"`
	FilenameL string   `json:"filename_large"`
	Name      string   `json:"name"`
	Excerpt   string   `json:"excerpt"`
	Date      string   `json:"date"`
	Location  Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

// Extrae metadata de la imagen (GPS y fecha)
func extractImageMetadata(filepath string) (lat, lon float64, date time.Time, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	defer f.Close()
	var meta exif2.Exif
	meta, err = imagemeta.Decode(f)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	fmt.Printf("metadata: %v", meta)

	// Obtener coordenadas GPS
	lat = meta.GPS.Latitude()
	lon = meta.GPS.Longitude()

	// Obtener fecha (preferir DateOriginal, luego DateCreated, luego DateModified)
	date = meta.CreateDate()

	return lat, lon, date, nil
}

// Obtiene información de ubicación desde Nominatim
func getLocationFromNominatim(lat, lon float64) (*NominatimResponse, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=en", lat, lon)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PhotoMetadataExtractor/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result NominatimResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Función principal que extrae toda la metadata
func extractPhotoMetadata(filepath, filename, filename_s, filename_m, filename_l, name, excerpt string) (*PhotoMetadata, error) {
	// Extraer metadata de la imagen
	lat, lon, date, err := extractImageMetadata(filepath)
	if err != nil {
		return nil, fmt.Errorf("error extracting image metadata: %v", err)
	}
	// Declarar location fuera del scope del if
	var location *NominatimResponse
	var locationName string

	// Si no hay coordenadas GPS, crear una ubicación vacía
	if lat == 0 && lon == 0 {
		// Crear una estructura vacía
		location = &NominatimResponse{}
	} else {
		// Get info from location
		var err error
		location, err = getLocationFromNominatim(lat, lon)
		if err != nil {
			return nil, fmt.Errorf("error getting location: %v", err)
		}
		// Build location.name as Name, Country
		locationName = fmt.Sprintf("%s, %s", location.Name, location.Address.Country)
	}

	// Formatear la fecha
	dateStr := date.Format("2006-01-02")

	result := &PhotoMetadata{
		Filename:  filename,
		FilenameS: filename_s,
		FilenameM: filename_m,
		FilenameL: filename_l,
		Name:      name,
		Excerpt:   excerpt,
		Date:      dateStr,
		Location: Location{
			Latitude:  lat,
			Longitude: lon,
			Name:      locationName,
		},
	}

	return result, nil
}

// Escribe el JSON en un archivo
func writeJSONToFile(data interface{}, outputPath string) error {
	// Crear el directorio si no existe
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Convertir a JSON con formato bonito
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Escribir el archivo
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

func GenerateJson(filename, filename_s, filename_m, filename_l, name, excerpt string, app_settings common.AppSettings) {
	// Ejemplo de uso
	imagePath := path.Join(app_settings.ImageDirectory, filename)

	// Extraer metadata
	metadata, err := extractPhotoMetadata(imagePath, filename, filename_s, filename_m, filename_l, name, excerpt)
	if err != nil {
		log.Error().Msgf("Error: %v\n", err)
		return
	}

	// Mostrar en consola
	jsonData, _ := json.MarshalIndent(metadata, "", "    ")
	fmt.Println("Metadata extraída:")
	fmt.Println(string(jsonData))

	// Escribir a archivo
	outputFilename := fmt.Sprintf("%s.json",
		strings.TrimSuffix(filename, filepath.Ext(filename)))
	outputPath := path.Join(app_settings.ImageDirectory, outputFilename)

	if err := writeJSONToFile(metadata, outputPath); err != nil {
		log.Error().Msgf("Error escribiendo archivo: %v\n", err)
		return
	}

	fmt.Printf("\nArchivo JSON guardado en: %s\n", outputPath)
}
