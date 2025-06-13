package common

type Image struct {
	Uuid     string   `json:"uuid"`
	Name     string   `json:"name"`
	Filename string   `json:"filename"`
	Filepath string   `json:"filepath"`
	Ext      string   `json:"extension"`
	Excerpt  string   `json:"excerpt"`
	Location Location `json:"location"`
	Date     string   `json:"date"`
}
