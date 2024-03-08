package endpoint_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/tests/mocks"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestAddPageHappyPath(t *testing.T) {
	databaseMock := mocks.DatabaseMock{
		AddPageHandler: func(string, string, string) (int, error) {
			return 0, nil
		},
	}

	page_data := admin_app.AddPageRequest{
		Title:   "Title",
		Content: "Content",
		Link:    "Link",
	}

	router := admin_app.SetupRoutes(app_settings, databaseMock, make(map[string]*lua.LState))
	responseRecorder := httptest.NewRecorder()

	body, _ := json.Marshal(page_data)
	request, _ := http.NewRequest(http.MethodPost, "/pages", bytes.NewBuffer(body))

	router.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	var response admin_app.PageResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)

	assert.NotNil(t, response.Id)
	assert.NotEmpty(t, response.Link)
	assert.Equal(t, page_data.Link, response.Link)
}
