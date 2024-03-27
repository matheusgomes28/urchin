package common

type Post struct {
	Title   string
	Content string
	Excerpt string
	Id      int
}

type Card struct {
	Uuid          string
	ImageLocation string
	JsonData      string
	SchemaName    string
}
