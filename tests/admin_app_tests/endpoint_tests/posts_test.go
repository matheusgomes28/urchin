package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/tests/mocks"
	"github.com/stretchr/testify/assert"
)

var app_settings = common.AppSettings{
	DatabaseAddress:  "localhost",
	DatabasePort:     3006,
	DatabaseUser:     "root",
	DatabasePassword: "root",
	DatabaseName:     "urchin",
	WebserverPort:    8080,
	ImageDirectory:   "../../../images",
}

func TestPostPostSuccess(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPostHandler: func(string, string, string) (int, error) {
			return 0, nil
		},
	}

	postData := admin_app.AddPostRequest{
		Title:   "Title",
		Excerpt: "Excerpt",
		Content: "Content",
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(postData)
	request, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	var response admin_app.PostIdResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.NotNil(t, response.Id)
}

func TestPostPostWithoutTitle(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPostHandler: func(string, string, string) (int, error) {
			return 0, nil
		},
	}

	postData := admin_app.AddPostRequest{
		Excerpt: "Excerpt",
		Content: "Content",
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(postData)
	request, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "missing required data")
	assert.NotNil(t, response.Err)
	assert.Contains(t, response.Err, "'Title'")
}

func TestPostPostWithoutExcerpt(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPostHandler: func(string, string, string) (int, error) {
			return 0, nil
		},
	}

	postData := admin_app.AddPostRequest{
		Title:   "Title",
		Content: "Content",
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(postData)
	request, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "missing required data")
	assert.NotNil(t, response.Err)
	assert.Contains(t, response.Err, "'Excerpt'")
}

func TestPostPostWithoutContent(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPostHandler: func(string, string, string) (int, error) {
			return 0, nil
		},
	}

	postData := admin_app.AddPostRequest{
		Title:   "Title",
		Excerpt: "Excerpt",
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(postData)
	request, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "missing required data")
	assert.NotNil(t, response.Err)
	assert.Contains(t, response.Err, "'Content'")
}

func TestPostPostNoBody(t *testing.T) {
	databaseMock := mocks.DatabaseMock{}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	request, _ := http.NewRequest(http.MethodPost, "/posts", nil)

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "no request body provided")
}

func TestPostPostInvalidPostRequest(t *testing.T) {
	databaseMock := mocks.DatabaseMock{}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal("{\"test\": \"Something\"}")
	request, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "invalid request body")
	assert.NotNil(t, response.Err)
}

func TestPostPostFailedSave(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPostHandler: func(string, string, string) (int, error) {
			return 0, fmt.Errorf("saving post failed")
		},
	}

	postData := admin_app.AddPostRequest{
		Title:   "Title",
		Excerpt: "Excerpt",
		Content: "Content",
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(postData)
	request, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "could not add post")
	assert.NotNil(t, response.Err)
}

func TestDeletePostSuccess(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		DeletePostHandler: func(int) (int, error) {
			return 1, nil
		},
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	request, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	var response admin_app.PostIdResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Id, 1)
}

func TestDeletePostFailedDelete(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		DeletePostHandler: func(int) (int, error) {
			return 0, fmt.Errorf("delete post failed")
		},
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	request, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "could not delete post")
	assert.NotNil(t, response.Err)
}

func TestDeletePostNotFound(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		DeletePostHandler: func(int) (int, error) {
			return 0, nil
		},
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock)
	responseRecorder := httptest.NewRecorder()

	request, _ := http.NewRequest(http.MethodDelete, "/posts/1", nil)

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, 404, responseRecorder.Code)
	var response common.ErrorResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Msg, "no post found")
}
