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

func TestGalleryMetadataFileDoesntExist(t *testing.T) {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3006,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
		CacheEnabled:     false,
		Galleries: map[string]common.Gallery{
			"test": {
				Name:        "test",
				Description: "my test gallery",
				Link:        "my-gallery-link",
				Thumbnail:   "thumbnail.jpg",
				Images:      []string{"some_manifest.toml"},
			},
		},
	}

	database_mock := mocks.DatabaseMock{}

	r := app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/gallery/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No images uploaded")
}

// func TestPageFailureInvalidKey(t *testing.T) {

// 	app_settings := common.AppSettings{
// 		DatabaseAddress:  "localhost",
// 		DatabasePort:     3006,
// 		DatabaseUser:     "root",
// 		DatabasePassword: "root",
// 		DatabaseName:     "urchin",
// 		WebserverPort:    8080,
// 	}

// 	database_mock := mocks.DatabaseMock{
// 		GetPageHandler: func(link string) (common.Page, error) {
// 			return common.Page{}, fmt.Errorf("invalid page")
// 		},
// 	}

// 	router := app.SetupRoutes(app_settings, database_mock)
// 	responseRecorder := httptest.NewRecorder()

// 	request, err := http.NewRequest("GET", "/pages/test", nil)

// 	require.Nil(t, err)

// 	router.ServeHTTP(responseRecorder, request)

// 	require.Equal(t, http.StatusNotFound, responseRecorder.Code)
// 	require.Contains(t, responseRecorder.Body.String(), "invalid page")

// }
