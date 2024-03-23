package common

type IdBinding struct {
	Id string `uri:"id" binding:"required"`
}

// TODO: Think this trough. Do we really need special bindings for every entity type?
type PostIdBinding struct {
	IdBinding
}

type ImageIdBinding struct {
	IdBinding
}
