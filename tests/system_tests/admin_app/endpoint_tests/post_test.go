package admin_endpoint_tests

import (
	"bytes"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/tests/helpers"
	"github.com/test-go/testify/require"
)

func TestPostExists(t *testing.T) {

	// Usual database setup
	app_settings := helpers.GetAppSettings()
	cleanup, db, err := helpers.SetupDb(app_settings)
	require.Nil(t, err)
	defer func() { require.Nil(t, cleanup()) }()

	// Send the post in
	add_post_request := admin_app.AddPostRequest{
		Title:   "Test Post Title",
		Excerpt: "test post excerpt",
		Content: "test post content",
	}

	w := httptest.NewRecorder()
	r := admin_app.SetupRoutes(app_settings, db)
	request_bytes, err := json.Marshal(add_post_request)
	require.Nil(t, err)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(request_bytes))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Make sure it's an expected response
	var add_post_response admin_app.PostIdResponse
	err = json.Unmarshal(w.Body.Bytes(), &add_post_response)
	require.Nil(t, err)

	// Get the post
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/posts/%d", add_post_response.Id), nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Make sure it's the expected content
	var get_post_response admin_app.GetPostResponse

	err = json.Unmarshal(w.Body.Bytes(), &get_post_response)
	require.Nil(t, err)

	require.Equal(t, get_post_response.Content, add_post_request.Content)
	require.Equal(t, get_post_response.Excerpt, add_post_request.Excerpt)
	require.Equal(t, get_post_response.Title, add_post_request.Title)
	require.Equal(t, get_post_response.Id, add_post_response.Id)
}
