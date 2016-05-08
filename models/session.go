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
}

func SessionFromRequest(req *http.Request) (s *Session, err error) {
	ts := getTokenFromHeader(req.Header.Get("Authorization"))

	if ts == "" {
		meta.App.Log.Warn("session not found")
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

func (s *Session) Token() (t string, err error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["session_id"] = s.SessionID

	t, err = token.SignedString([]byte(meta.App.Conf.Self.SigningKey))

	return t, err
}

func (s *Session) AttachUser() (err error) {
	if s.userAttached {
		return nil
	}

	s.User, err = UserBySessionKey(s.SessionKey)
	s.userAttached = true

	return err
}

func (s *Session) DestroyKey() error {
	return DestroySessionsByKey(s.SessionKey)
}

func (s *Session) Destroy() error {
	_, err := db.NamedExec(`DELETE FROM sessions WHERE session_id = :session_id LIMIT 1`, s)
	return err
}

func DestroySessionsByKey(key string) (err error) {
	_, err = db.Exec(`DELETE FROM sessions WHERE session_key = $1`, key)
	return err
}

func getTokenFromHeader(h string) string {
	if h == "" {
		return ""
	}

	if strings.HasPrefix(h, "Bearer ") {
		return strings.Replace(h, "Bearer ", "", 1)
	}

	return ""
}

func generateKey() string {
	return uuid.NewV4().String()
}
