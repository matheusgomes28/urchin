package mocks

import (
	"fmt"

	"github.com/matheusgomes28/urchin/common"
)

type DatabaseMock struct {
	GetPostHandler     func(int) (common.Post, error)
	GetPostsHandler    func(int, int) ([]common.Post, error)
	AddImageHandler    func(string, string, string, string) error
	GetImageHandler    func(string) (common.Image, error)
	DeleteImageHandler func(string) error
}

func (db DatabaseMock) GetPosts(limit int, offset int) ([]common.Post, error) {
	return db.GetPostsHandler(limit, offset)
}

func (db DatabaseMock) GetPost(post_id int) (common.Post, error) {
	return db.GetPostHandler(post_id)
}

func (db DatabaseMock) AddPost(title string, excerpt string, content string) (int, error) {
	return 0, nil
}

func (db DatabaseMock) ChangePost(id int, title string, excerpt string, content string) error {
	return nil
}

func (db DatabaseMock) DeletePost(id int) error {
	return fmt.Errorf("not implemented")
}

func (db DatabaseMock) AddImage(id string, file_name string, alt_text string, ext string) error {
	return db.AddImageHandler(id, file_name, alt_text, ext)
}

func (db DatabaseMock) GetImage(id string) (common.Image, error) {
	return db.GetImageHandler(id)
}

func (db DatabaseMock) DeleteImage(id string) error {
	return db.DeleteImageHandler(id)
}
