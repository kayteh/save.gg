package models

type Markdown struct {
	Raw  string `json:"raw" gorethink:"raw"`
	HTML string `json:"html" gorethink:"html"`
}
