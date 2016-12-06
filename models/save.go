package models

import (
	"time"
)

type Save struct {
	ID           string `db:"save_id" json:"id,omitempty"`
	CanonicalURL string `db:"canonical_url" json:"canonical_url"`
	CustomURL    string `db:"custom_url" json:"custom_url,omitempty"`

	OwnerID string `db:"owner_id" json:"-"`
	Owner   *User  `json:"owner,omitempty"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Title       string
	Description *Markdown

	MetadataID string                 `db:"metadata_id" json:"-"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`

	FileEntityID string          `db:"file_entity_id" json:"-"`
	//FileEntity   *SaveFileEntity `json:"file_entity,omitempty"`

	presentable bool
}
