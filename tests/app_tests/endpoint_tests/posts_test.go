package endpoint_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPostSuccess(t *testing.T) {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3006,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
	}

	database_mock := mocks.DatabaseMock{
		GetPostHandler: func(post_id int) (common.Post, error) {
			return common.Post{
				Title:   "TestPost",
				Content: "TestContent",
				Excerpt: "TestExcerpt",
				Id:      post_id,
			}, nil
		},
	}

	r := app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/post/0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
}

func TestPostFailure(t *testing.T) {

	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3006,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
	}

	database_mock := mocks.DatabaseMock{
		GetPostHandler: func(post_id int) (common.Post, error) {
			return common.Post{
				Title:   "TestPost",
				Content: "TestContent",
				Excerpt: "TestExcerpt",
				Id:      post_id,
			}, nil
		},
	}

	router := app.SetupRoutes(app_settings, database_mock)
	responseRecorder := httptest.NewRecorder()

	request, err := http.NewRequest("GET", "/post/sampleString", nil)

	if err != nil {
		t.Errorf("could not create request: %v", err)
	}

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)

}
