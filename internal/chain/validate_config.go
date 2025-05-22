package chain

import (
	"crypto/rsa"
	"errors"

	"github.com/pot876/sauth/internal/util"
)

type ValidateServiceConfig struct {
	publicKeys map[string]*rsa.PublicKey
}

// SetKeys accepts pairs - [key_id, public_key]
func (c *ValidateServiceConfig) SetKeys(keys [][2]string) error {
	publicKeys := map[string]*rsa.PublicKey{}

	for i := range keys {
		keyID := keys[i][0]
		publicKey := keys[i][1]

		decodedPublicKey, err := util.PublicKeyDecode(publicKey)
		if err != nil {
			return err
		}

		publicKeys[keyID] = decodedPublicKey
	}

	c.publicKeys = publicKeys
	return nil
}

func (c *ValidateServiceConfig) validate() error {
	if len(c.publicKeys) == 0 {
		return errors.New("missing keys")
	}
	return nil
}
