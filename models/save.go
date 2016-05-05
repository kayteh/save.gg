package models

import (
	r "github.com/dancannon/gorethink"
	"save.gg/sgg/meta"
)

type Save struct {
	ID           string
	Identity     string
	Title        string
	Owner        *User
	Description  Markdown
	Game         *Game
	FileResource *FileResource
	Privacy      string
}
