package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossoreslp/go-uuid-v4"
	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/tests/helpers"
	"github.com/matheusgomes28/urchin/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type imageResponse struct {
	Id string `json:"id"`
}

func TestPostImage(t *testing.T) {
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := helpers.CreateFormImagePart(writer, "file", "test.png", "image/png")
		require.Nil(t, err)

		err = png.Encode(part, helpers.CreateImage())
		require.Nil(t, err)
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
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := helpers.CreateFormImagePart(writer, "file", "test.png", "image/png")
		require.Nil(t, err)

		text := bytes.NewBufferString("This is some dumy text to check the content test")
		_, err = io.Copy(part, text)
		require.Nil(t, err)
	}()

	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestPostImageWrongFileContentType(t *testing.T) {
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := helpers.CreateFormImagePart(writer, "file", "test.png", "application/json")
		require.Nil(t, err)
		require.Nil(t, err)

		img := helpers.CreateImage()

		err = png.Encode(part, img)
		require.Nil(t, err)
	}()

	req, _ := http.NewRequest("POST", "/images", pr)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestGetImageNoDatabaseEntry(t *testing.T) {
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)

	get_recorder := httptest.NewRecorder()
	uuid, _ := uuid.New()

	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("GET", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(get_recorder, req)

	assert.Equal(t, 404, get_recorder.Code)
}

func TestGetImageNoImageFile(t *testing.T) {
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)

	get_recorder := httptest.NewRecorder()
	uuid, _ := uuid.New()

	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("GET", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(get_recorder, req)

	assert.Equal(t, 404, get_recorder.Code)
}

func TestDeleteImage(t *testing.T) {
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	post_recorder := httptest.NewRecorder()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		part, err := helpers.CreateFormImagePart(writer, "file", "test.png", "image/png")
		require.Nil(t, err)

		img := helpers.CreateImage()

		err = png.Encode(part, img)
		require.Nil(t, err)
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

func TestDeleteImageNoImageFile(t *testing.T) {
	database_mock := mocks.DatabaseMock{}

	r := admin_app.SetupRoutes(app_settings, database_mock)
	delete_recorder := httptest.NewRecorder()

	uuid, _ := uuid.New()
	uri := fmt.Sprintf("/images/%s", uuid.String())
	req, _ := http.NewRequest("DELETE", uri, bytes.NewBuffer([]byte{}))
	r.ServeHTTP(delete_recorder, req)

	assert.Equal(t, 200, delete_recorder.Code)
}
