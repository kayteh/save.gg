package models

import (
	"fmt"
	// "github.com/jmoiron/sqlx"
	"bytes"
	"database/sql"
	"github.com/satori/go.uuid"
	"github.com/tv42/slug"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"save.gg/sgg/meta"
	"save.gg/sgg/util/errors"
	"strings"
	"time"
)

var (
	allowedCharacters = regexp.MustCompile(`\A[a-zA-Z0-9_\-](?:[a-zA-Z0-9_\-\ ]+[a-zA-Z0-9_\-])?\z`)
)

type User struct {
	ID         string   `db:"user_id",json:"id,omitempty"`
	Slug       string   `json:"slug"`
	Username   string   `json:"username"`
	Email      string   `json:"email,omitempty"`
	Secret     string   `json:"secret,omitempty"`
	SessionKey string   `json:"session_key,omitempty" db:"session_key"`
	ACL        []byte   `json:"-"`
	RealACL    []string `json:"acl"`
	SubLevel   string   `json:"sub_level" db:"sub_level"`
	Activated  bool     `json:"activated"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	OldSecrets []byte `json:"old_secrets,omitempty" db:"old_secrets"`

	KnownIPs []byte `json:"known_ips,omitempty" db:"known_ips"`

	// Joined data
	Saves []Save `json:"saves,omitempty"`

	// Metadata
	presentable bool
}

type UserOldSecrets []struct {
	sql.NullString
	Secret        string
	InvalidatedAt time.Time
}

type UserKnownIPs []struct {
	sql.NullString
	IP       string
	LastSeen time.Time
}

func UserBySlug(slug string) (*User, error) {

	user := &User{}

	err := db.Get(user, `
	SELECT 
		user_id, slug, username,
		secret, acl, sub_level,
		activated, created_at, updated_at,
		old_secrets, known_ips
	FROM users 
	WHERE deleted_at IS NULL 
		AND slug = $1 
	LIMIT 1`, slug)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

	user.setACLPq()

	return user, err

}

func (u *User) Presentable() *User {
	newUser := new(User)
	*newUser = *u

	newUser.ID = ""
	newUser.Email = ""
	newUser.Secret = ""
	newUser.SessionKey = ""
	// newUser.OldSecrets = UserOldSecrets{}
	// newUser.KnownIPs = UserKnownIPs{}
	newUser.presentable = true

	return newUser
}

func (u *User) FetchSaves(f *FetchConfig) (s *[]Save, err error) {
	//s, err = SavesByUser(u, f)

	return s, err
}

// Saves a user to the truth store, then touches.
func (u *User) Save() (err error) {

	if u.presentable == true {
		return errors.UserPresentableSave
	}

	if u.ID == "" {
		return errors.UserNoIDSave
	}

	u.UpdatedAt = time.Now()
	u.ACL = u.getACLPq()

	_, err = db.NamedExec(`
		UPDATE users SET
			secret=:secret:,
			slug=:slug:,
			username=:username:,
			email=:email:,
			acl=:acl:,
			sub_level=:sub_level:,
			activated=:activated:,
			updated_at=:updated_at:,
			old_secrets=:old_secrets:,
			known_ips=:known_ips:,
			session_key=:session_key:
		WHERE user_id=:user_id:
	`, &u)

	return err
}

// Inserts a user into the database.
func (u *User) Insert() (err error) {
	if u.ID == "" {
		u.ID = uuid.NewV4().String()
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.ACL = u.getACLPq()

	_, err = db.NamedExec(`
		INSERT INTO users ( 
			user_id,  email, slug, 
			secret, username, acl, 
			sub_level, activated, 
			created_at, updated_at
		) VALUES (
			:user_id:, :email:, :slug:, 
			:secret:, :username:, :acl:, 
			:sub_level:, :activated:, 
			:created_at:, :updated_at: 
		)`,
		&u)

	return err
}

func (u *User) getACLPq() []byte {
	// if len(u.realACL) == 0 {
	// 	return []byte("{}")
	// }

	return []byte(fmt.Sprintf("{%s}", strings.Join(u.RealACL, ",")))
}

func (u *User) setACLPq() {
	b := u.ACL

	if len(b) == 0 {
		return
	}

	p := bytes.Split(b[1:len(b)-1], []byte(","))

	po := make([]string, len(p))

	for _, r := range p {
		po = append(po, string(r))
	}

	u.RealACL = po

	return
}

func (u *User) AppendACL(role string) {
	u.RealACL = append(u.RealACL, role)

	return
}

func (u *User) GetACL() []string {
	return u.RealACL
}

// func (u *User) HasACL(role string) bool {
// 	r := []byte(role)

// }

// Test a hashed secret/password for correctness.
func (u *User) TestSecret(password string) (o bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Secret), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Creates a new secret from the password.
func (u *User) CreateSecret(password string) (err error) {
	err = validatePassword(password)
	if err != nil {
		return err
	}

	s, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return err
	}

	u.Secret = string(s)

	return nil
}

// Gracefully changes a username for slug, security, and cache reasons.
func (u *User) SetUsername(username string) (err error) {

	if !allowedCharacters.MatchString(username) {
		return errors.UsernameInvalid
	}

	if len(username) > meta.App.Conf.Validation.UsernameLength {
		return errors.UsernameTooLong
	}

	s := TransformUsernameToSlug(username)

	// check for a bad slug
	for _, badSlug := range meta.App.Conf.Validation.DisallowedSlugs {
		if s == badSlug {
			return errors.UsernameDisallowed
		}
	}

	// check if user exists
	cu, err := UserBySlug(s)
	if err != nil && err != errors.UserNotFound {
		return err
	}

	if err == nil && cu != nil {
		return errors.UsernameTaken
	}

	u.Username = username
	u.Slug = s

	return nil
}

// Validates a password is of a certain constraint.
func validatePassword(password string) (err error) {

	if len(password) < meta.App.Conf.Validation.PasswordLength {
		return errors.UserPasswordTooShort
	}

	return nil

}

func TransformUsernameToSlug(username string) string {
	return slug.Slug(username)
}
