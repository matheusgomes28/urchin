package common

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
	Id      int    `json:"id"`
}
