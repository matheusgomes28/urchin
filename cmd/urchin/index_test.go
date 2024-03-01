package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matheusgomes28/urchin/app"
	"github.com/matheusgomes28/urchin/common"

	"github.com/stretchr/testify/assert"
)

type DatabaseMock struct{}

func (db DatabaseMock) GetPosts() ([]common.Post, error) {
	return []common.Post{
		{
			Title:   "TestPost",
			Content: "TestContent",
			Excerpt: "TestExcerpt",
			Id:      0,
		},
	}, nil
}

func (db DatabaseMock) GetPost(post_id int) (common.Post, error) {
	return common.Post{}, fmt.Errorf("not implemented")
}

func (db DatabaseMock) AddPost(title string, excerpt string, content string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (db DatabaseMock) ChangePost(id int, title string, excerpt string, content string) error {
	return nil
}

func (db DatabaseMock) DeletePost(id int) error {
	return fmt.Errorf("not implemented")
}

func (db DatabaseMock) AddImage(string, string, string) error {
	return fmt.Errorf("not implemented")
}

func (db DatabaseMock) AddCard(string, string, string, string) error {
	return fmt.Errorf("not implemented")
}

func TestIndexPing(t *testing.T) {
	app_settings := common.AppSettings{
		DatabaseAddress:  "localhost",
		DatabasePort:     3006,
		DatabaseUser:     "root",
		DatabasePassword: "root",
		DatabaseName:     "urchin",
		WebserverPort:    8080,
	}

	database_mock := DatabaseMock{}

	r := app.SetupRoutes(app_settings, database_mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "TestPost")
	assert.Contains(t, w.Body.String(), "TestExcerpt")
}
