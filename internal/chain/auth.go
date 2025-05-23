package chain

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pot876/sauth/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type Responce struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type IServiseAuth interface {
	Login(ctx context.Context, realmID uuid.UUID, login, password []byte) (*Responce, error)
	Refresh(ctx context.Context, refreshToken []byte) (*Responce, error)
}

type AuthService struct {
	credentialProvider ICredentialsProvider

	config atomic.Value
	logger *zerolog.Logger
}

func (s *AuthService) UpdateConfig(ctx context.Context, config *AuthServiceConfig) error {
	if err := config.validate(); err != nil {
		return err
	}

	s.setConfig(config)

	return nil
}

func NewAuthService(ctx context.Context, credentialProvider ICredentialsProvider, config *AuthServiceConfig) (IServiseAuth, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config err: %v", err)
	}

	if credentialProvider == nil {
		return nil, fmt.Errorf("credentialProvider was not set")
	}

	result := &AuthService{
		credentialProvider: credentialProvider,
		logger:             zerolog.Ctx(ctx),
	}
	result.setConfig(config)

	return result, nil
}

func (s *AuthService) Login(ctx context.Context, realmID uuid.UUID, login, password []byte) (*Responce, error) {
	authInfo, err := s.credentialProvider.ValidateAuth(ctx, realmID, login, password)
	if err != nil {
		return nil, err
	}

	access_token, refresh_token, claims, err := s.issueTokensOnLogin(authInfo)
	if err != nil {
		return nil, err
	}

	s.log("tl %s/%d", claims.SessionID, claims.Seq)
	return &Responce{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken []byte) (*Responce, error) {
	claims, err := s.refreshValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	access_token, refresh_token, err := s.issueTokensOnRefresh(claims)
	if err != nil {
		return nil, err
	}

	s.log("tr %s/%d", claims.SessionID, claims.Seq)
	return &Responce{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

func (s *AuthService) issueTokensOnLogin(authInfo any) (string, string, *Claims, error) {
	now := time.Now().UTC()
	sessionID := "sid" + uuid.NewString()

	cfgPointer := s.getConfig()

	claims := &Claims{
		ExpiresAt: now.Add(cfgPointer.AccessTokenTTL).Unix(),
		IssuedAt:  now.Unix(),

		Issuer:    cfgPointer.Issuer,
		SessionID: sessionID,
		Nonce:     uuid.NewString(),
		Seq:       0,
		Siat:      now.Unix(),

		Info: authInfo,
	}

	accessToken, err := s.issueAccessToken(claims)
	if err != nil {
		return "", "", nil, err
	}

	claims.ExpiresAt = now.Add(cfgPointer.RefreshTokenTTL).Unix()
	claims.Nonce = uuid.NewString()
	refreshToken, err := s.issueRefreshToken(claims)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, claims, nil
}

func (s *AuthService) refreshValidateToken(_ context.Context, refreshToken []byte) (*Claims, error) {
	cfgPointer := s.getConfig()

	tokenDecoded, err := jwt.ParseWithClaims(string(refreshToken), &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		keyID, ok := t.Header["kid"].(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		key, ok := cfgPointer.RefreshPublicKeys[keyID]
		if !ok {
			return nil, ErrKeyNotFound
		}

		return key, nil
	})
	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			switch e.Inner {
			case ErrUnexpectedSigningMethod, ErrInvalidToken, ErrKeyNotFound:
				return nil, e.Inner
			}
		}
		return nil, err
	}

	claims, ok := tokenDecoded.Claims.(*Claims)
	if !ok || !tokenDecoded.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Typ != TypRefresh {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *AuthService) issueTokensOnRefresh(claims *Claims) (string, string, error) {
	cfgPointer := s.getConfig()

	now := time.Now().UTC()
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(cfgPointer.AccessTokenTTL).Unix()
	claims.Nonce = uuid.NewString()
	claims.Seq += 1
	accessToken, err := s.issueAccessToken(claims)
	if err != nil {
		return "", "", err
	}

	claims.ExpiresAt = now.Add(cfgPointer.RefreshTokenTTL).Unix()
	claims.Nonce = uuid.NewString()
	refreshToken, err := s.issueRefreshToken(claims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) issueAccessToken(claims *Claims) (string, error) {
	cfgPointer := s.getConfig()
	claims.Typ = TypAccess
	accessToken, err := newWithClaims(jwt.SigningMethodRS256, claims, cfgPointer.AccessKeyID).SignedString(cfgPointer.AccessPrivateKey)
	if err != nil {
		return "", fmt.Errorf("sign access token err: %w", err)
	}

	return accessToken, nil
}

func (s *AuthService) issueRefreshToken(claims *Claims) (string, error) {
	cfgPointer := s.getConfig()
	claims.Typ = TypRefresh
	refreshToken, err := newWithClaims(jwt.SigningMethodRS256, claims, cfgPointer.RefreshKeyID).SignedString(cfgPointer.RefreshPrivateKey)
	if err != nil {
		return "", fmt.Errorf("sign refresh token err: %w", err)
	}

	return refreshToken, nil
}

func (s *AuthService) getConfig() *AuthServiceConfig {
	// return (*AuthServiceConfig)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&s.config))))
	return s.config.Load().(*AuthServiceConfig)
}
func (s *AuthService) setConfig(config *AuthServiceConfig) {
	// atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&s.config)), unsafe.Pointer(config))
	s.config.Store(config)
}

func (s *AuthService) log(format string, v ...any) {
	if s.logger != nil {
		s.logger.Log().Msgf(format, v...)
	}
}

func (s *AuthService) RegisterMetrics(r prometheus.Registerer, prefix string) {
	_ = util.RegisterMetrics(s.credentialProvider, r, prefix)
}

func newWithClaims(method jwt.SigningMethod, claims jwt.Claims, keyID string) *jwt.Token {
	return &jwt.Token{
		Header: map[string]any{
			"typ": "JWT",
			"alg": method.Alg(),
			"kid": keyID,
		},
		Claims: claims,
		Method: method,
	}
}
