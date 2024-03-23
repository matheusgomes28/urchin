package admin_app

import "github.com/matheusgomes28/urchin/common"

// Extracted all bindings and requests structs into a single package to
// organize the data in a simpler way. Every domain object supporting
// CRUD endpoints has their own structures to handle the http methods.

type DeletePostBinding struct {
	common.IdBinding
}

type AddPostRequest struct {
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type ChangePostRequest struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Content string `json:"content"`
}

type DeleteImageBinding struct {
	common.IdBinding
}
