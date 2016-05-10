package models

import (
	"database/sql"
	"encoding/json"
	radix "github.com/mediocregopher/radix.v2/redis"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"save.gg/sgg/meta"
	"save.gg/sgg/util/errors"
	"save.gg/sgg/util/pq"
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

	//TODO(kkz): Consider moving this somewhere that it could be auto-generated per environment.
	passwordNonce = "4E0933FAAF9UAddQCEKwl3CSPgXoR12bIUmR5Gq6QODPHLc6jGBbFVLG5DWl7dzwbxNluVdKh4G1rEmNefh5uL52YBeEfZzp"
)

const (
	// Allowed characters in an email. This shit is cray.
	// Copied from https://github.com/go-playground/validator/blob/v8/regexes.go#L16
	emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:\\(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22)))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

	// Allowed characters in a username.
	//
	// A-Z, a-z, 0-9, _, and - are allowed, and spaces are allowed anywhere but the start and end.
	//
	// I'm aware that if someone wants a username with two letters and 28 spaces wants to do that,
	// I'm allowing it thus far.
	usernameRegexString = `\A[a-zA-Z0-9_\-](?:(?:[a-zA-Z0-9_\-\ ]+)?[a-zA-Z0-9_\-])?\z`
)

var (
	allowedCharacters = regexp.MustCompile(usernameRegexString)
	acceptableEmail   = regexp.MustCompile(emailRegexString)
)

// A user, duh!
type User struct {
	ID         string        `db:"user_id" json:"id,omitempty"`
	Slug       string        `json:"slug"`
	Username   string        `json:"username"`
	Email      string        `json:"email,omitempty"`
	Secret     string        `json:"secret,omitempty"`
	SessionKey string        `json:"session_key,omitempty" db:"session_key"`
	ACL        *pq.TextArray `json:"acl"`
	SubLevel   string        `json:"sub_level" db:"sub_level"`
	Activated  bool          `json:"activated"`

	// NoSQL data

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"-" db:"deleted_at"`

	// Metadata
	presentable bool
}

// Finds a user by their slug.
func UserBySlug(slug string) (*User, error) {

	user := &User{}

	cacheOk, err := userCachedBySlug(user, slug)
	if cacheOk {
		return user, nil
	}

	if err != nil && err != errors.CacheMiss {
		return nil, err
	}

	err = userQuery(user, "slug", slug)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

	// if everything's ok to this point,
	// let's touch the missing user in a goroutine.
	if !cacheOk {
		go user.Touch()
	}

	return user, err

}

// Find a user by their email.
func UserByEmail(email string) (*User, error) {

	user := &User{}

	err := userQuery(user, "email", email)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

	return user, err

}

// Find a user by session key.
func UserBySessionKey(key string) (*User, error) {

	user := &User{}

	err := userQuery(user, "session_key", key)

	if err == sql.ErrNoRows {
		return nil, errors.UserNotFound
	}

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

// Try to get a user from the cache. If it's missing
func userCachedBySlug(userPtr *User, slug string) (ok bool, err error) {

	r := redis.Cmd("GET", "user:"+slug)
	if r.Err != nil {
		return false, err
	}

	if r.IsType(radix.Nil) {
		return false, errors.CacheMiss
	}

	b, err := r.Bytes()
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(b, userPtr)
	if err != nil {
		return false, err
	}

	meta.App.Log.Info("cache hit")

	return true, nil
}

// Create a new user object that can be prefilled properly. For Scanning, a bare *User is fine.
func NewUser() *User {
	return &User{
		ACL: &pq.TextArray{},
	}
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

// Patch merges a a new user map into the current.
func (u *User) Patch(incoming map[string]interface{}) (err error) {
	for key, val := range incoming {

		switch key {

		case "secret":
			err = u.CreateSecret(val.(string))
			DestroySessionsByKey(u.SessionKey)

		case "username":
			err = u.SetUsername(val.(string))

		case "email":
			err = u.SetEmail(val.(string))

		}

		if err != nil {
			return err
		}

	}

	return nil
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

// Check if a user can modify this user.
// This is either itself or someone with the ACL role "admin".
func (u *User) UserCanModify(mu *User) bool {

	if u.ID == mu.ID || mu.ACL.Has(UserACLAdmin) {
		return true
	}

	return false

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

	go u.Touch()

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
	u.ACL.Append("user")

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

// Caches the user into the cache, as opposed to a box.
func (u *User) Touch() {
	if u.presentable == true {
		meta.App.Log.WithField("user id", u.ID).Error("attempted touch of presentable user")
		return
	}

	du := &User{}
	err := userQuery(du, "user_id", u.ID)
	if err != nil {
		meta.App.Log.WithError(err).WithField("user id", u.ID).Error("fetch error")
		return
	}

	b, err := json.Marshal(du)
	if err != nil {
		meta.App.Log.WithError(err).WithField("user id", u.ID).Error("cache encode error")
		return
	}

	r := redis.Cmd("SET", "user:"+du.Slug, b, "EX", "86400")
	if r.Err != nil {
		meta.App.Log.WithError(err).WithField("user id", u.ID).Error("cache save error")
	}
}

// Removes a user from the cache (usually so it can be re-cached due to slug change.)
func (u *User) Uncache() {
	redis.Cmd("DEL", "user:"+u.Slug)
}

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
	err = bcrypt.CompareHashAndPassword([]byte(u.Secret), []byte(password+passwordNonce))

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

	s, err := bcrypt.GenerateFromPassword([]byte(password+passwordNonce), 14)

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

	u.Uncache()

	u.Username = username
	u.Slug = s

	return nil
}

// Validates then sets an email.
func (u *User) SetEmail(email string) (err error) {
	if !acceptableEmail.MatchString(email) {
		return errors.UserEmailInvalid
	}

	u.Email = email

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
