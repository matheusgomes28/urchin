package admin_post_tests

import (
	"bytes"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"

	admin_app "github.com/matheusgomes28/urchin/admin-app"
	"github.com/matheusgomes28/urchin/tests/system_tests/helpers"
	"github.com/pressly/goose/v3"
)

type AddPostRequest struct {
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type AddPostResponse struct {
	Id int `json:"id"`
}

type GetPostResponse struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

func TestPostExists(t *testing.T) {

	// This is gonna be the in-memory mysql
	app_settings := helpers.GetAppSettings(20)
	go helpers.RunDatabaseServer(app_settings)
	database, err := helpers.WaitForDb(app_settings)
	require.Nil(t, err)
	goose.SetBaseFS(helpers.EmbedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		require.Nil(t, err)
	}

	if err := goose.Up(database.Connection, "migrations"); err != nil {
		require.Nil(t, err)
	}

	// Send the post in
	add_post_request := admin_app.AddPostRequest{
		Title:   "Test Post Title",
		Excerpt: "test post excerpt",
		Content: "test post content",
	}

	w := httptest.NewRecorder()
	r := admin_app.SetupRoutes(app_settings, database)
	request_bytes, err := json.Marshal(add_post_request)
	require.Nil(t, err)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(request_bytes))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Make sure it's an expected response
	var add_post_response AddPostResponse
	err = json.Unmarshal(w.Body.Bytes(), &add_post_response)
	require.Nil(t, err)

	// Get the post
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/posts/%d", add_post_response.Id), nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Make sure it's the expected content
	var get_post_response GetPostResponse

	err = json.Unmarshal(w.Body.Bytes(), &get_post_response)
	require.Nil(t, err)

	require.Equal(t, get_post_response.Content, add_post_request.Content)
	require.Equal(t, get_post_response.Excerpt, add_post_request.Excerpt)
	require.Equal(t, get_post_response.Title, add_post_request.Title)
	require.Equal(t, get_post_response.Id, add_post_response.Id)
}
