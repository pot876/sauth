package util

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func DecodeKeyPair(privateKey, publicKey string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	decodedPrivateKey, err := PrivateKeyDecode(privateKey)
	if err != nil {
		return nil, nil, err
	}
	decodedPublicKey, err := PublicKeyDecode(publicKey)
	if err != nil {
		return nil, nil, err
	}

	return decodedPrivateKey, decodedPublicKey, nil
}

func PublicKeyDecode(publicKey string) (*rsa.PublicKey, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("base64 decode err: %v", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("parseRSAPublicKeyFromPEM err: %v", err)
	}

	return key, nil
}

func PrivateKeyDecode(privateKey string) (*rsa.PrivateKey, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("base64 decode err: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parseRSAPrivateKeyFromPEM err: %v", err)
	}

	return key, nil
}
