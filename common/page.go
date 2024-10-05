package common

type Page struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Link    string `json:"link"`
}
