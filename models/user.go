package models

import (
	"fmt"
	// "github.com/jmoiron/sqlx"
	"bytes"
	"database/sql"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"save.gg/sgg/meta"
	"save.gg/sgg/util/errors"
	"strings"
	"time"
)

const (
	UserSubFree  = "free"
	UserSubPro   = "pro"
	UserSubElite = "elite"

	UserACLAdmin      = "admin"
	UserACLGlobalMod  = "glomod"
	UserACLPartner    = "partner"
	UserACLTester     = "tester"
	UserACLContinuum  = "dev:c7m"
	UserACLAtmosphere = "dev:atmos"
	UserACLMomentum   = "dev:m6m"
	UserACLSupernova  = "dev:nova"
)

var (
	// Allowed characters in a username.
	//
	// A-Z, a-z, 0-9, _, and - are allowed, and spaces are allowed anywhere but the start and end.
	//
	// I'm aware that if someone wants a username with two letters and 28 spaces wants to do that,
	// I'm allowing it thus far.
	allowedCharacters = regexp.MustCompile(`\A[a-zA-Z0-9_\-](?:(?:[a-zA-Z0-9_\-\ ]+)?[a-zA-Z0-9_\-])?\z`)
)

// A user, duh!
//
// Some notes
//
// Due to postgres being better than every other SQL DB, any special
// data structures (jsonb, arrays, etc), need to be []byte slices and converted before and after
// DB fetches/commits. An example of this workaround are the ACL functions on this type. In the future,
// helpers might be written.
type User struct {
	ID         string   `db:"user_id" json:"id,omitempty"`
	Slug       string   `json:"slug"`
	Username   string   `json:"username"`
	Email      string   `json:"email,omitempty"`
	Secret     string   `json:"secret,omitempty"`
	SessionKey string   `json:"session_key,omitempty" db:"session_key"`
	ACL        []byte   `json:"-"`   // sqlx-readable ACL
	RealACL    []string `json:"acl"` // Human-readable ACL
	SubLevel   string   `json:"sub_level" db:"sub_level"`
	Activated  bool     `json:"activated"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"-" db:"deleted_at"`

	// Joined data
	Saves []Save `json:"saves,omitempty"`

	// Metadata
	presentable bool
}

// Finds a user by their slug.
func UserBySlug(slug string) (*User, error) {

	user := &User{}

	err := userQuery(user, "slug", slug)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

	user.setACLPq()

	return user, err

}

// Find a user by their email.
func UserByEmail(email string) (*User, error) {

	user := &User{}

	err := userQuery(user, "email", email)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

	user.setACLPq()

	return user, err

}

// Find a user by session key.
func UserBySessionKey(key string) (*User, error) {

	user := &User{}

	err := userQuery(user, "session_key", key)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

	user.setACLPq()

	return user, err

}

// Abstraction over various single-constraint queries for getting one user.
func userQuery(userPtr *User, constraint, value string) error {
	return db.Get(userPtr, `
	SELECT 
		user_id, slug, username, email,
		secret, acl, sub_level,
		activated, created_at, 
		updated_at, session_key
	FROM users 
	WHERE deleted_at IS NULL 
		AND `+constraint+` = $1 
	LIMIT 1`, value)
}

// Strips secret or un-needed values out of public-facing, or presented, data.
// The idea is that if you still need one of these fields, you would fill it after calling this.
// The model pointer this returns is non-commitable.
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

// Authenicates a user based on their username or email, and their password.
func UserAuth(handle, password string) (u *User, err error) {
	if strings.Contains(handle, "@") {
		u, err = UserByEmail(handle)
	} else {
		u, err = UserBySlug(TransformUsernameToSlug(handle))
	}

	if err == errors.UserNotFound {
		return nil, errors.UserAuthBadHandle
	}

	if err != nil {
		return nil, err
	}

	pass, err := u.CheckSecret(password)

	if !pass {
		return nil, errors.UserAuthBadPassword
	}

	return u, nil
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
			secret=:secret,
			slug=:slug,
			username=:username,
			email=:email,
			acl=:acl,
			sub_level=:sub_level,
			activated=:activated,
			updated_at=:updated_at,
			session_key=:session_key
		WHERE user_id=:user_id
	`, &u)

	return err
}

// Inserts a user into the database.
func (u *User) Insert() (err error) {
	if u.ID == "" {
		u.ID = uuid.NewV1().String()
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.SessionKey = generateKey()
	u.AppendACL("user")
	u.ACL = u.getACLPq()

	_, err = db.NamedExec(`
		INSERT INTO users ( 
			 user_id,  email,  slug,  secret,  username,  acl,  activated,  created_at,  updated_at,  session_key
		) VALUES (
			:user_id, :email, :slug, :secret, :username, :acl, :activated, :created_at, :updated_at, :session_key
		)`,
		&u)

	return err
}

// Soft-deletes the user from the database.
// This should trigger cleanup jobs before being fully removed.
func (u *User) Delete() (err error) {
	u.DeletedAt = time.Now()
	_, err = db.NamedExec(`UPDATE users SET deleted_at=:deleted_at WHERE user_id = :user_id`, &u)
	return err
}

// Hard-deletes a user from the database. This is uncoverable.
func (u *User) Purge() (err error) {
	_, err = db.Exec(`DELETE FROM users WHERE user_id = $1 LIMIT 1`, u.ID)
	return err
}

// Gets the byte array from the human-readable ACL to store in postgres.
func (u *User) getACLPq() []byte {
	// if len(u.realACL) == 0 {
	// 	return []byte("{}")
	// }

	return []byte(fmt.Sprintf("{%s}", strings.Join(u.RealACL, ",")))
}

// Sets the human-readable ACL from go-sql's byte array
func (u *User) setACLPq() {
	b := u.ACL

	if len(b) <= 2 {
		return
	}

	p := bytes.Split(b[1:len(b)-1], []byte(","))

	for _, r := range p {
		if len(r) == 0 {
			continue
		} else {
			u.RealACL = append(u.RealACL, string(r))
		}
	}

	return
}

// Appends a new ACL to the list
func (u *User) AppendACL(role string) {
	u.RealACL = append(u.RealACL, role)

	return
}

// Gets the real ACL values
func (u *User) GetACL() []string {
	return u.RealACL
}

// func (u *User) HasACL(role string) bool {
// 	r := []byte(role)

// }

// Creates a new Session for the user.
func (u *User) CreateSession() (s *Session, err error) {
	s, err = NewSession(u.SessionKey)
	if err != nil {
		return nil, err
	}

	meta.App.Log.WithField("session", s).Info("session")

	if s.SessionKey != u.SessionKey {
		u.SessionKey = s.SessionKey
		err = u.Save()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Test a hashed secret/password for correctness.
func (u *User) CheckSecret(password string) (o bool, err error) {
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

// Transforms a username into a "url-safe" slug
func TransformUsernameToSlug(username string) string {
	slug := strings.ToLower(username)
	slug = strings.Replace(slug, " ", "-", -1)

	return slug
}
