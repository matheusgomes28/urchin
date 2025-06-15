package common

type Gallery struct {
	Name        string   `toml:"name"`
	Description string   `toml:"description"`
	Link        string   `toml:"link"`
	Thumbnail   string   `toml:"thumbail"`
	Images      []string `toml:"images"`
}
