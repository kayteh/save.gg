package models

import (
	_ "save.gg/sgg/meta"
)

type Save struct {
	ID          string
	Identity    string
	Title       string
	Owner       *User
	Description Markdown
	Privacy     string
}
