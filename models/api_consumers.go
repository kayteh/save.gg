package models

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"save.gg/sgg/util/errors"
	"time"
)

type Consumer struct {
	Key       string    `db:"api_key" json:"api_key"`
	Public    string    `db:"public_key" json:"public_key"`
	Internal  bool      `db:"is_internal"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	OwnerID   string    `db:"owner_user_id"`
	Active    bool      `db:"active"`
}

// Fetch a consumer by it's API key. This will return an error
// if the consumer is inactive, but will still give you the consumer instance.
// If the error is relevant, use it, if not, don't. This is to allow the frontend
// and consumer api to list/show without issue.
func ConsumerByAPIKey(key string) (*Consumer, error) {
	var c *Consumer

	err := db.Get(c, `
	SELECT
		api_key,
		public_key,
		is_internal,
		created_at,
		owner_user_id,
		active
	FROM consumers
	WHERE api_key = $1
	`, key)

	if err != nil {
		return nil, err
	}

	if !c.Active {
		return c, errors.ConsumerAPIKeyInactive
	}

	return c, nil
}

func NewConsumer(owner *User) (*Consumer, error) {
	return &Consumer{
		Key:       generateKey(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    true,
		OwnerID:   owner.ID,
		Internal:  false,
	}, nil
}

func (c *Consumer) Insert() error {
	_, err := db.NamedExec(`
	INSERT INTO consumers (
		api_key, public_key,
		created_at, updated_at,
		active, owner_user_id,
		is_internal
	) VALUES (
		:api_key, :public_key,
		:created_at, :updated_at,
		:active, :owner_user_id,
		:is_internal
	)
	`, c)

	return err
}

func (c *Consumer) Save() error {
	c.UpdatedAt = time.Now()
	_, err := db.NamedExec(`
	UPDATE consumers SET
		public_key=:public_key,
		updated_at=:updated_at,
		active=:active,
		owner_user_id=:owner_user_id,
		is_internal=:is_internal
	WHERE api_key = :api_key
	`, c)

	return err
}

// Generates an ECDSA P-384 DER-encoded key pair for signing requests.
//
// We keep the public key for decoding, the consumer keeps the private key for theirself.
// This function returns a base64-encoded version of the private key, and it should be decoded
// for actual use. The consumer web UI should decode this and offer it as a blob.
//
// Since we do not keep the private key, we cannot give them their private key after it's generated.
// The consumer must take this into account.
func (c *Consumer) GenerateKeys() (privateKey string, err error) {
	k, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return "", err
	}

	pk, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return "", err
	}

	pubkDer, err := x509.MarshalPKIXPublicKey(k.Public())
	if err != nil {
		return "", err
	}

	c.Public = base64.StdEncoding.EncodeToString(pubkDer)

	return base64.StdEncoding.EncodeToString(pk), nil

}

// Decodes the public key stored in the database.
//
// This should be preferred over using *Consumer.Public since it's obviously encoded.
func (c *Consumer) DecodedPublic() (b []byte, err error) {
	b, err = base64.StdEncoding.DecodeString(c.Public)
	return b, err
}

// Lookup function for JWT consumer signatures. Returns the public key PEM in byte form.
func LookupPublicKey(key string) (k *ecdsa.PublicKey, err error) {
	c, err := ConsumerByAPIKey(key)
	if err != nil {
		return nil, err
	}

	b, err := c.DecodedPublic()
	if err != nil {
		return nil, err
	}

	kp, err := x509.ParsePKIXPublicKey(b)
	k = kp.(*ecdsa.PublicKey)
	return k, nil
}

// Adds a tick to the rate limiting of this consumer key.
func (c *Consumer) RLTick(t string) {
	consumerApiKeyLimitTick(c.Key, t)
}

//TODO(kkz): implement this
func consumerApiKeyLimitCheck(key, t string) bool {
	return true
}

//TODO(kkz): implement this
func consumerApiKeyLimitTick(key, t string) {
	return
}
