package common

type Image struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	AltText string `json:"alt_text"`
	Ext     string `json:"extension"`
}
