package chain

import (
	"context"
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v4"
	"github.com/prometheus/client_golang/prometheus"
)

type IServiseValidate interface {
	Validate(context.Context, string) error
	Info(context.Context, string) (any, error)
}

type ValidateService struct {
	publicKeys map[string]*rsa.PublicKey
}

func NewValidateService(config *ValidateServiceConfig) (IServiseValidate, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &ValidateService{
		publicKeys: config.publicKeys,
	}, nil
}

func (s *ValidateService) Validate(_ context.Context, token string) error {
	tokenDecoded, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		keyID, ok := t.Header["kid"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		key, ok := s.publicKeys[keyID]
		if !ok {
			return nil, ErrKeyNotFound
		}
		return key, nil
	})
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			switch e.Inner {
			case ErrUnexpectedSigningMethod, ErrInvalidToken, ErrKeyNotFound:
				return e.Inner
			}
		}
		return err
	}

	_, ok := tokenDecoded.Claims.(*Claims)
	if !ok || !tokenDecoded.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (s *ValidateService) Info(_ context.Context, token string) (any, error) {
	return nil, nil
}

func (s *ValidateService) RegisterMetrics(r prometheus.Registerer, prefix string) {

}
