package mocks

import (
	"github.com/matheusgomes28/urchin/common"
)

type DatabaseMock struct {
	GetPostHandler        func(int) (common.Post, error)
	GetPostsHandler       func(int, int) ([]common.Post, error)
	AddPostHandler        func(string, string, string) (int, error)
	DeletePostHandler     func(int) (int, error)
	AddPageHandler        func(string, string, string) (int, error)
	GetPageHandler        func(string) (common.Page, error)
	AddCardHandler        func(string, string, string, string) (string, error)
	AddChardSchemaHandler func(string, string, string, string) (string, error)
}

func (db DatabaseMock) GetPosts(limit int, offset int) ([]common.Post, error) {
	return db.GetPostsHandler(limit, offset)
}

func (db DatabaseMock) GetPost(post_id int) (common.Post, error) {
	return db.GetPostHandler(post_id)
}

func (db DatabaseMock) AddPost(title string, excerpt string, content string) (int, error) {
	return db.AddPostHandler(title, excerpt, content)
}

func (db DatabaseMock) ChangePost(id int, title string, excerpt string, content string) error {
	return nil
}

func (db DatabaseMock) DeletePost(id int) (int, error) {
	return db.DeletePostHandler(id)
}

func (db DatabaseMock) AddPage(title string, content string, link string) (int, error) {
	return db.AddPageHandler(title, content, link)
}

func (db DatabaseMock) GetPage(link string) (common.Page, error) {
	return db.GetPageHandler(link)
}

func (db DatabaseMock) AddCard(title string, image string, schema string, content string) (string, error) {
	return db.AddCardHandler(title, image, schema, content)
}

func (db DatabaseMock) AddCardSchema(json_id string, json_schema string, json_title string, schema string) (string, error) {
	return db.AddChardSchemaHandler(json_id, json_schema, json_title, schema)
}
