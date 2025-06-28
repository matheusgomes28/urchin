package common

type CardSchema struct {
	Uuid   string   `json:"uuid"`
	Title  string   `json:"title"`
	Schema string   `json:"schema"`
	Cards  []string `json:"cards"`
}
