package chain

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/pot876/sauth/internal/util"
)

type AuthServiceConfig struct {
	Issuer string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	AccessKeyID      string
	AccessPrivateKey *rsa.PrivateKey
	AccessPublicKey  *rsa.PublicKey

	RefreshKeyID      string
	RefreshPrivateKey *rsa.PrivateKey
	RefreshPublicKeys map[string]*rsa.PublicKey
}

func (c *AuthServiceConfig) SetAccessKeys(privateKey, publicKey, keyID string) error {
	privateKeyDecoded, publicKeyDecoded, err := util.DecodeKeyPair(privateKey, publicKey)
	if err != nil {
		return err
	}

	c.AccessPrivateKey = privateKeyDecoded
	c.AccessPublicKey = publicKeyDecoded
	c.AccessKeyID = keyID

	return nil
}

func (c *AuthServiceConfig) SetRefreshKeys(privateKey, publicKey, keyID string) error {
	privateKeyDecoded, publicKeyDecoded, err := util.DecodeKeyPair(privateKey, publicKey)
	if err != nil {
		return err
	}

	c.RefreshPrivateKey = privateKeyDecoded
	c.RefreshPublicKeys = map[string]*rsa.PublicKey{
		keyID: publicKeyDecoded,
	}
	c.RefreshKeyID = keyID

	return nil
}

func (c *AuthServiceConfig) validate() error {
	if c.AccessPrivateKey == nil {
		return errors.New("missing AccessPrivateKey")
	}
	if c.AccessPublicKey == nil {
		return errors.New("missing AccessPrivateKey")
	}

	if c.RefreshPrivateKey == nil {
		return errors.New("missing RefreshPrivateKey")
	}
	if c.RefreshPublicKeys == nil {
		return errors.New("missing RefreshPublicKeys")
	}

	if c.AccessTokenTTL == 0 {
		return errors.New("missing AccessTokenTTL")
	}
	if c.RefreshTokenTTL == 0 {
		return errors.New("missing RefreshTokenTTL")
	}

	return nil
}
