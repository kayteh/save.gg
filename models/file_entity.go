package models

import (
	r "github.com/dancannon/gorethink"
	"time"
)

// A FileEntity is a file stored somewhere on the network. It can be a save pack,
// an image, a 3D render, whatever. Ideally, specific types of FileEntities would
// inherit this, e.g. SaveFileEntity.
type FileEntity struct {
	ID        string `json:"id"`
	Layer     string `json:"layer"`
	URL       string `json:"url"`
	Alternate string `json:"alt_url"`

	Owner   *User  `json:"owner,omitempty"`
	OwnerID string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	MarkedForFreeze bool `json:"marked_for_freeze"`
}

func Fetch(key, value string) []*FileEntity {

}
