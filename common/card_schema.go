package common

type CardSchema struct {
	Title  string   `json:"title"`
	Schema string   `json:"schema"`
	Cards  []string `json:"cards"`
}
