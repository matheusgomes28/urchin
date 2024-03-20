package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/tests/mocks"
	"github.com/stretchr/testify/assert"
)

type imageResponse struct {
	Id string `json:"id"`
}

func TestPostImage(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string) error {
			return nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := createTestForm(writer, "file", "test.png", "image/png")
		if err != nil {
			t.Error(err)
		}

		img := createImage()
		if err != nil {
			t.Error(err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response imageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.NotNil(t, response.Id)
}

func TestPostImageNotAnImageFile(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string) error {
			return nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := createTestForm(writer, "file", "test.png", "image/png")
		if err != nil {
			t.Error(err)
		}

		text := bytes.NewBufferString("This is some dumy text to check the content test")
		_, err = io.Copy(part, text)
		if err != nil {
			t.Error(err)
		}
	}()

	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestPostImageWrongFileContentType(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string) error {
			return nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := createTestForm(writer, "file", "test.png", "application/json")
		if err != nil {
			t.Error(err)
		}

		img := createImage()

		if err != nil {
			t.Error(err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestPostImageFailedToCreateDatabaseEntry(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string) error {
			return errors.New("test error")
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := createTestForm(writer, "file", "test.png", "image/png")
		if err != nil {
			t.Error(err)
		}

		img := createImage()
		if err != nil {
			t.Error(err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()

	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)

}

func createTestForm(writer *multipart.Writer, fieldname string, filename string, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
	h.Set("Content-Type", contentType)
	return writer.CreatePart(h)
}

// Creating an image in memory for testing: https://yourbasic.org/golang/create-image/
func createImage() image.Image {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
			}
		}
	}

	return img
}
