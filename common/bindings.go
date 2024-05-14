package common

type IntIdBinding struct {
	Id int `uri:"id" binding:"required"`
}

type StringIdBinding struct {
	Id string `uri:"id" binding:"required"`
}

type PostIdBinding struct {
	IntIdBinding
}

type ImageIdBinding struct {
	// This is the uuid of an image to be retrieved
	Filename string `uri:"name" binding:"required"`
}
