package admin_app

import "github.com/matheusgomes28/urchin/common"

// Extracted all bindings and requests structs into a single package to
// organize the data in a simpler way. Every domain object supporting
// CRUD endpoints has their own structures to handle the http methods.

type AddPageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

type DeletePostBinding struct {
	common.IntIdBinding
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

type AddImageRequest struct {
	Alt string `json:"alt"`
}

// TODO : Does this still need to be here?
// TODO : Are we handling images apart from file adds?
type DeleteImageBinding struct {
	Name string `uri:"name" binding:"required"`
}

type AddCardRequest struct {
	Title   string `json:"title"`
	Image   string `json:"image"`
	Schema  string `json:"schema"`
	Content string `json:"content"`
}

type AddCardSchemaRequest struct {
	JsonId     string `json:"$id"`
	JsonSchema string `json:"$schema"`
	JsonTitle  string `json:"title"`
	Schema     string `json:"schema"`
}
