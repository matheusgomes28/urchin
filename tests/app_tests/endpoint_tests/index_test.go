package endpoint_tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/common"
	"github.com/matheusgomes28/urchin/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexSuccess(t *testing.T) {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3006,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
	}

	database_mock := mocks.DatabaseMock{
		GetPostsHandler: func(limit int, offset int) ([]common.Post, error) {
			return []common.Post{
				{
					Title:   "TestPost",
					Content: "TestContent",
					Excerpt: "TestExcerpt",
					Id:      0,
				},
			}, nil
		},
	}

	r := app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/page/0", nil)
	require.Nil(t, err)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
	assert.Contains(t, w.Body.String(), "TestExcerpt")
}

func TestIndexFailToGetPosts(t *testing.T) {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3006,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
	}

	error_msg := "could not find posts"
	database_mock := mocks.DatabaseMock{
		GetPostsHandler: func(limit int, offset int) ([]common.Post, error) {
			return []common.Post{}, fmt.Errorf("%s", error_msg)
		},
	}

	r := app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	require.Nil(t, err)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), error_msg)
}
