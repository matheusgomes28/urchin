package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/fossoreslp/go-uuid-v4"
	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type imageResponse struct {
	Id string `json:"id"`
}

func TestPostImage(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string, ext string) error {
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
		AddImageHandler: func(id string, file_name string, alt_text string, ext string) error {
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
		AddImageHandler: func(id string, file_name string, alt_text string, ext string) error {
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
		AddImageHandler: func(id string, file_name string, alt_text string, ext string) error {
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

func TestGetImage(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string, ext string) error {
			return nil
		},

		GetImageHandler: func(id string) (common.Image, error) {
			return common.Image{
				Uuid:    id,
				Name:    "test",
				AltText: "default",
				Ext:     ".png",
			}, nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	post_recorder := httptest.NewRecorder()

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

	// TODO: We have to create the image first. Maybe there's a better way to do this?
	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(post_recorder, req)

	var response imageResponse
	_ = json.Unmarshal(post_recorder.Body.Bytes(), &response)

	get_recorder := httptest.NewRecorder()
	uri := fmt.Sprintf("/images/%s", response.Id)
	req, _ = http.NewRequest("GET", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(get_recorder, req)

	assert.Equal(t, 200, get_recorder.Code)

	assert.Equal(t, "image/png", get_recorder.Result().Header.Get("Content-Type"))
}

func TestGetImageNoDatabaseEntry(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		GetImageHandler: func(id string) (common.Image, error) {
			return common.Image{
				Uuid:    id,
				Name:    "test",
				AltText: "default",
				Ext:     ".png",
			}, nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)

	get_recorder := httptest.NewRecorder()
	uuid, _ := uuid.New()

	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("GET", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(get_recorder, req)

	assert.Equal(t, 404, get_recorder.Code)
}

func TestGetImageNoImageFile(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		GetImageHandler: func(id string) (common.Image, error) {
			return common.Image{}, errors.New("failed")
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)

	get_recorder := httptest.NewRecorder()
	uuid, _ := uuid.New()

	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("GET", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(get_recorder, req)

	assert.Equal(t, 404, get_recorder.Code)
}

func TestDeleteImage(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		AddImageHandler: func(id string, file_name string, alt_text string, ext string) error {
			return nil
		},

		GetImageHandler: func(id string) (common.Image, error) {
			return common.Image{
				Uuid:    id,
				Name:    "test",
				AltText: "default",
				Ext:     ".png",
			}, nil
		},

		DeleteImageHandler: func(string) error {
			return nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	post_recorder := httptest.NewRecorder()

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

	// TODO: We have to create the image first. Maybe there's a better way to do this?
	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(post_recorder, req)

	var response imageResponse
	_ = json.Unmarshal(post_recorder.Body.Bytes(), &response)

	delete_recorder := httptest.NewRecorder()
	uri := fmt.Sprintf("/images/%s", response.Id)
	req, _ = http.NewRequest("DELETE", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(delete_recorder, req)

	assert.Equal(t, 200, delete_recorder.Code)

	var image_id_response admin_app.ImageIdResponse

	err := json.Unmarshal(delete_recorder.Body.Bytes(), &image_id_response)
	require.Nil(t, err)

	require.Equal(t, image_id_response.Id, response.Id)
}

func TestDeleteImageNoDatabaseEntry(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		GetImageHandler: func(id string) (common.Image, error) {
			return common.Image{}, errors.New("failed")
		},

		DeleteImageHandler: func(string) error {
			return nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	delete_recorder := httptest.NewRecorder()

	uuid, _ := uuid.New()
	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("DELETE", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(delete_recorder, req)

	assert.Equal(t, 404, delete_recorder.Code)
}

func TestDeleteImageNoImageFile(t *testing.T) {
	database_mock := mocks.DatabaseMock{
		GetImageHandler: func(id string) (common.Image, error) {
			return common.Image{
				Uuid:    id,
				Name:    "test",
				AltText: "default",
				Ext:     ".png",
			}, nil
		},

		DeleteImageHandler: func(string) error {
			return nil
		},
	}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	delete_recorder := httptest.NewRecorder()

	uuid, _ := uuid.New()
	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("DELETE", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(delete_recorder, req)

	assert.Equal(t, 200, delete_recorder.Code)
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
	width := 1
	height := 1

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	return image.NewRGBA(image.Rectangle{upLeft, lowRight})
}
