package models

type Markdown struct {
	Markdown string `db:"markdown"`
	HTML     string `db:"html"`
}
