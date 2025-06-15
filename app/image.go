package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/database"
	"github.com/matheusgomes28/urchin/views"
	"github.com/rs/zerolog/log"
)

// This function assumes that the image file
// exists and is a valid image
func populateImageMetadata(filename string, app_settings common.AppSettings) common.Image {

	ext := path.Ext(filename)
	name := strings.TrimSuffix(filename, ext)
	metadata_file := path.Join(app_settings.ImageDirectory, fmt.Sprintf("%s.json", name))
	file_path := path.Join("/images/data", filename)

	// Check if a json metadata file exists
	metadata_contents, err := os.ReadFile(metadata_file)
	if err == nil {
		var image common.Image
		err = json.Unmarshal(metadata_contents, &image)

		if err == nil {
			image.Ext = ext
			image.Uuid = name
			image.Filename = filename
			image.Filepath = file_path
			return image
		}
	}

	return common.Image{
		Uuid:     filename[:len(filename)-len(ext)],
		Name:     filename,
		Filename: filename,
		Filepath: file_path,
		Ext:      ext,
	}
}

func imagesHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	pageNum := 1 // Default to page 0
	if pageNumQuery := c.Param("num"); pageNumQuery != "" {
		num, err := strconv.Atoi(pageNumQuery)
		if err == nil && num > 0 {
			pageNum = num
		} else {
			log.Error().Msgf("Invalid page number: %s", pageNumQuery)
		}
	}

	// Get all the files inside the image directory
	files, err := os.ReadDir(app_settings.ImageDirectory)
	if err != nil {
		log.Error().Msgf("could not read files in image directory: %v", err)
		return []byte{}, err
	}
	filepaths := common.Map(files, func(file os.DirEntry) string {
		return path.Join(app_settings.ImageDirectory, file.Name())
	})

	// Filter all the non-images out of the list
	image_files := common.FilterStrings(filepaths, func(filepath string) bool {
		ext := path.Ext(filepath)
		_, contains := common.ValidImageExtensions[ext]
		return contains
	})

	valid_images, err := common.GetImages(image_files, 10, pageNum, app_settings)
	if err != nil {
		return []byte{}, err
	}

	index_view := views.MakeImagesPage(valid_images, app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns)
	html_buffer := bytes.NewBuffer(nil)

	err = index_view.Render(c, html_buffer)
	if err != nil {
		log.Error().Msgf("Could not render index: %v", err)
		return []byte{}, err
	}

	return html_buffer.Bytes(), nil
}

func imageHandler(c *gin.Context, app_settings common.AppSettings, database database.Database) ([]byte, error) {
	var get_image_binding common.ImageIdBinding
	if err := c.ShouldBindUri(&get_image_binding); err != nil {
		return nil, err
	}

	// if not cached, create the cache
	filename := get_image_binding.Filename
	ext := path.Ext(get_image_binding.Filename)
	name := filename[:len(filename)-len(ext)]

	image := common.Image{
		Uuid: name,
		Name: filename,
		Ext:  ext,
	}

	return renderHtml(c, views.MakeImagePage(image, app_settings.AppNavbar.Links, app_settings.AppNavbar.Dropdowns))
}
