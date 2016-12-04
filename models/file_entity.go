package models

import (
	"encoding/json"
	radix "github.com/mediocregopher/radix.v2/redis"
	r "gopkg.in/dancannon/gorethink.v2"
	"save.gg/sgg/util/errors"
	"time"
)

// A FileEntity is a file stored somewhere on the network. It can be a save pack,
// an image, a 3D render, whatever. Ideally, specific types of FileEntities would
// inherit this, e.g. SaveFileEntity.
type FileEntity struct {
	ID        string `json:"id"`
	Layer     string `json:"layer"`   // Layer is the section of the file network it's on. This is covered in docs/FileServices.md
	URL       string `json:"url"`     // URL on CDN
	Alternate string `json:"alt_url"` // Direct URL if CDN is failing

	Owner   *User  `json:"owner,omitempty" gorethink:"-"`
	OwnerID string `json:"-" gorethink:"owner_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	MarkedForFreeze bool `json:"marked_for_freeze"`

	presentable bool
}

// Takes in a map[string]interface to pass into the rethink filter query.
func FilesByFilter(f map[string]interface{}) ([]*FileEntity, error) {
	c, err := r.Table("file_entities").Filter(f).Run(rethink)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	var fe []*FileEntity
	err = c.All(fe)
	if err != nil {
		return nil, err
	}

	return fe, nil

}

func fileCachedById(filePtr *[]byte, id string) (ok bool, err error) {
	f := redis.Cmd("GET", "file_entity:"+id)
	if f.Err != nil {
		return false, f.Err
	}

	if f.IsType(radix.Nil) {
		return false, errors.CacheMiss
	}

	b, err := f.Bytes()
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(b, filePtr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (f *FileEntity) Presentable() *FileEntity {
	newFE := new(FileEntity)
	*newFE = *f

	newFE.presentable = true

	return newFE
}
