package chain

import "time"

type Claims struct {
	ExpiresAt int64  `json:"exp,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	KeyID     string `json:"kid,omitempty"`

	Issuer    string `json:"iss,omitempty"`
	SessionID string `json:"sid,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Seq       int    `json:"seq,omitempty"`
	Siat      int64  `json:"siat,omitempty"`

	Info any `json:"info,omitempty"`
}

func (c *Claims) Valid() error {
	now := time.Now().Unix()

	if !c.verifyExpiresAt(now) {
		return ErrExpired
	}
	return nil
}

func (c *Claims) verifyExpiresAt(cmp int64) bool {
	now := time.Unix(cmp, 0)
	exp := time.Unix(c.ExpiresAt, 0)
	return now.Before(exp)
}

func (c *Claims) verifyIssuedAt(cmp int64) bool {
	now := time.Unix(cmp, 0)
	iat := time.Unix(c.IssuedAt, 0)
	return now.After(iat) || now.Equal(iat)
}
