package models

import (
	"database/sql"
	"github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"net/http"
	"save.gg/sgg/meta"
	"save.gg/sgg/util/errors"
	"strings"
	"time"
)

type Session struct {
	SessionID  string    `db:"session_id"`
	SessionKey string    `db:"session_key"`
	CreatedAt  time.Time `db:"created_at"`
	User       *User

	userAttached bool
	token        string
}

// Retrieve a session based on the JWT token stored in either
// the Authorization: Bearer header, or the session cookie.
func SessionFromRequest(req *http.Request) (s *Session, err error) {
	ts := getTokenFromHeader(req.Header.Get("Authorization"))

	if ts == "" {
		c, err := req.Cookie(meta.App.Conf.Self.SessionCookie)
		if err != nil && err != http.ErrNoCookie {
			meta.App.Log.WithError(err).Warn("cookie fetch error")
			return nil, errors.SessionNotFound
		}

		if err != http.ErrNoCookie {
			ts = c.Value
		}
	}

	if ts == "" {
		//meta.App.Log.Warn("session not found")
		return nil, errors.SessionNotFound
	}

	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		return []byte(meta.App.Conf.Self.SigningKey), nil
	})

	if err != nil {
		meta.App.Log.WithError(err).Error("jwt errored")
		return nil, errors.SessionTokenInvalid
	}

	if !token.Valid {
		meta.App.Log.Warn("token invalid")
		return nil, errors.SessionTokenInvalid
	}

	sid := token.Claims["session_id"].(string)

	s, err = SessionByID(sid)

	return s, err

}

// Finds a session by it's ID.
func SessionByID(id string) (*Session, error) {

	session := &Session{}

	err := db.Get(session, `
		SELECT
			session_key,
			session_id,
			created_at
		FROM sessions
		WHERE session_id = $1
	`, id)

	if err == sql.ErrNoRows {
		return nil, errors.SessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return session, nil

}

// Checks if a session ID exists. This is a lot less expensive than getting the entire session.
func SessionExists(id string) (bool, error) {
	var val int
	err := db.QueryRow("SELECT count(1) as count FROM sessions WHERE session_id = $1", id).Scan(&val)
	return val > 0, err
}

// Create a new session. This commits the session to the database.
func NewSession(key string) (s *Session, err error) {

	s = &Session{}

	if key == "" {
		s.SessionKey = generateKey()
	} else {
		s.SessionKey = key
	}
	s.SessionID = generateKey()
	s.CreatedAt = time.Now()

	_, err = db.NamedExec(`
		INSERT INTO sessions (
			session_key,
			session_id,
			created_at
		) VALUES (
			:session_key,
			:session_id,
			:created_at
		)
	`, &s)

	return s, err
}

// Computes a session JWT token. This will only compute once, so is safe to call multiple times.
//
// Note that the token is not stateful. It only contains claims to be verified elsewhere.
func (s *Session) Token() (t string, err error) {
	if s.token != "" {
		return t, nil
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["session_id"] = s.SessionID

	t, err = token.SignedString([]byte(meta.App.Conf.Self.SigningKey))
	s.token = t
	return t, err
}

// Sets the token to HTTP cookie.
func (s *Session) SetCookie(w http.ResponseWriter) {
	t, _ := s.Token()

	cookie := &http.Cookie{
		Name:     meta.App.Conf.Self.SessionCookie,
		Value:    t,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

// Attaches the user to the session. This can be called multiple times with no side effects.
func (s *Session) AttachUser() (err error) {
	if s.userAttached {
		return nil
	}

	s.User, err = UserBySessionKey(s.SessionKey)
	s.userAttached = true

	return err
}

// Wrapper to DestroySessionsByKey using this session's key.
func (s *Session) DestroyKey() error {
	return DestroySessionsByKey(s.SessionKey)
}

// Destroy this specific session.
func (s *Session) Destroy() error {
	_, err := db.NamedExec(`DELETE FROM sessions WHERE session_id = :session_id LIMIT 1`, s)
	return err
}

// Destroy all sessions by a key. This is useful for password changes or.. anything else, really.
func DestroySessionsByKey(key string) (err error) {
	_, err = db.Exec(`DELETE FROM sessions WHERE session_key = $1`, key)
	return err
}

// Parse the input string as a Bearer schema. That just means stripping the "Bearer " part, but.. hey.
func getTokenFromHeader(h string) string {
	if h == "" {
		return ""
	}

	if strings.HasPrefix(h, "Bearer ") {
		return strings.Replace(h, "Bearer ", "", 1)
	}

	return ""
}

// Generate a UUIDv4
func generateKey() string {
	return uuid.NewV4().String()
}
