package common

type Card struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	Schema  string `json:"schema"`
	Content string `json:"content"`
}
