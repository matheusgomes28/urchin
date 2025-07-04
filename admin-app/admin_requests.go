package admin_app

import (
	"mime/multipart"

	"github.com/matheusgomes28/urchin/common"
)

// swagger:ignore Extracted all bindings and requests structs into a single package to
// swagger:ignore organize the data in a simpler way. Every domain object supporting
// swagger:ignore CRUD endpoints has their own structures to handle the http methods.

// swagger:parameters add_page_request AddPageRequest
type AddPageRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}

// swagger:parameters delete_post_ DeletePostBinding
type DeletePostBinding struct {
	common.IntIdBinding
}

// swagger:parameters addPost
type AddPostRequest struct {
	// Title of the post
	// in: body
	// required: true
	Title string `json:"title"`
	// Excerpt of the post
	// in: body
	Excerpt string `json:"excerpt"`
	// Content of the post
	// in: body
	Content string `json:"content"`
}

// swagger:parameters changePost
type ChangePostRequest struct {
	// ID of the post
	// in: body
	// required: true
	Id int `json:"id"`
	// Title of the post
	// in: body
	Title string `json:"title"`
	// Excerpt of the post
	// in: body
	Excerpt string `json:"excerpt"`
	// Content of the post
	// in: body
	Content string `json:"content"`
}

// swagger:parameters postImage
type AddImageRequest struct {
	// The image file
	// in: formData
	// required: true
	Filedata *multipart.FileHeader `form:"filedata"`
	// Excerpt for the image
	// in: formData
	Excerpt string `form:"excerpt"`
}

type DeleteImageBinding struct {
	Name string `uri:"name" binding:"required"`
}
