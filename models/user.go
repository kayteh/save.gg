package models

import (
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"save.gg/sgg/meta"
	"time"
)

type User struct {
	ID         string   `json:"id",omitempty`
	Slug       string   `json:"slug"`
	Username   string   `json:"username"`
	Email      string   `json:"email",omitempty`
	Secret     string   `json:"secret",omitempty`
	SessionKey string   `json:"session_key",omitempty`
	ACL        []string `json:"acl"`
	SubLevel   string   `json:"sub_level"`
	Activated  bool     `json:"activated"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at",omitempty`

	OldSecrets UserOldSecrets `json:"old_secrets",omitempty`

	KnownIPs UserKnownIPs `json:"known_ips",omitempty`

	// Joined data
	Saves []Save `json:"saves",omitempty`
}

type UserOldSecrets []struct {
	Secret        string
	InvalidatedAt time.Time
}

type UserKnownIPs []struct {
	IP       string
	LastSeen time.Time
}

func UserBySlug(slug string) (*User, error) {

	var user User
	err := meta.App.Cache.Get(fmt.Sprintf("user:%s", slug), &user)

	if err != nil && err.Error() != "cache miss" {
		return nil, err
	}

	if err.Error() == "cache miss" {
		ur, err := UserBySlugMiss(slug)
		if err != nil {
			return nil, err
		}

		ur.Touch()

		return ur, nil
	}

	return &user, nil

}

func UserBySlugMiss(slug string) (u *User, err error) {

	var user User

	c, err := r.Table("users").Filter(map[string]interface{}{"slug": slug}).Run(meta.App.Rethink, nil)
	if err != nil {
		return nil, err
	}

	err = c.One(&user)
	if err != nil {
		return nil, err
	}

	err = c.Close()
	if err != nil {
		return user, err
	}

	return &user, nil

}

func (u *User) Presentable() *User {
	newUser := new(User)
	*newUser = *u

	newUser.ID = ""
	newUser.Email = ""
	newUser.Secret = ""
	newUser.Secret = ""
	newUser.SessionKey = ""
	newUser.OldSecrets = []UserOldSecrets{}
	newUser.KnownIPs = []UserKnownIPs{}

	return newUser
}

func (u *User) FetchSaves(f *FetchConfig) (s *[]Save, err error) {
	s, err = SavesByUser(u, f)

	return s, err
}

// Saves a user to the truth store, then touches.
func (u *User) Save() (err error) {

	if u.ID == "" {
		return nil
	}

	_, err = r.Table("users").Get(u.ID).Update(u).WriteRun(meta.App.Rethink, nil)
	if err != nil {
		return err
	}

	return u.Touch()
}

// Touch re-caches data from the source of truth.
func (u *User) Touch() (err error) {
	if u.ID == "" {
		return nil
	}

	return meta.App.Cache.Set(fmt.Sprintf("user:%s", u.Slug), u, meta.App.Conf.Cache.TTL.User)
}

// Test a hashed secret/password for correctness.
func (u *User) TestSecret(password string) (o bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Secret), []byte(password))
	o = err != nil
	return o, err
}

// Creates a new secret from the password.
func (u *User) CreateSecret(password string) (err error) {
	s, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return err
	}

	u.Secret = string(s)

	return nil
}
