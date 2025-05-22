package main

import (
	"context"
	"errors"

	"github.com/pot876/sauth/internal/chain"
	"github.com/pot876/sauth/internal/config"
	"github.com/pot876/sauth/internal/delegates/credentials_provider_file"
	"github.com/pot876/sauth/internal/delegates/credentials_provider_postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAuth(ctx context.Context, cfg *config.Config) (chain.IServiseAuth, error) {
	authConfig := &chain.AuthServiceConfig{}
	authConfig.SetAccessKeys(cfg.Auth.AccessTokenPrivateKey, cfg.Auth.AccessTokenPublicKey, cfg.Auth.AccessTokenKeyID)
	authConfig.SetRefreshKeys(cfg.Auth.RefreshTokenPrivateKey, cfg.Auth.RefreshTokenPublicKey, cfg.Auth.RefreshTokenKeyID)
	authConfig.AccessTokenTTL = cfg.Auth.AccessTokenExpiresIn
	authConfig.RefreshTokenTTL = cfg.Auth.RefreshTokenExpiresIn

	cp, err := buildCredentialsProvider(ctx, cfg)
	if err != nil {
		return nil, err
	}

	s, err := chain.NewAuthService(ctx, cp, authConfig)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func buildCredentialsProvider(ctx context.Context, cfg *config.Config) (chain.ICredentialsProvider, error) {
	var credentialProvider chain.ICredentialsProvider
	if cfg.CredentialsProviderPostgresEnabled {
		pgxcfg, err := pgxpool.ParseConfig(cfg.CredentialsProviderPostgres.ConnectionString)
		if err != nil {
			return nil, err
		}

		pool, err := pgxpool.NewWithConfig(ctx, pgxcfg)
		if err != nil {
			return nil, err
		}

		err = pool.Ping(ctx)
		if err != nil {
			return nil, err
		}

		credentialProvider = credentials_provider_postgres.New(pool)
	} else if cfg.CredentialsProviderFileEnabled {
		if cfg.CredentialsProviderFile.Format != "json" {
			return nil, errors.New("component: credentials_provider_file, error: unsupported file format")
		}

		users, err := credentials_provider_file.ReadUsers(cfg.CredentialsProviderFile.Path)
		if err != nil {
			return nil, err
		}

		_credentialProvider, err := credentials_provider_file.New(users)
		if err != nil {
			return nil, err
		}

		credentialProvider = _credentialProvider
	}

	return credentialProvider, nil
}

func NewValidate(ctx context.Context, cfg *config.Config) (chain.IServiseValidate, error) {
	validateConfig := &chain.ValidateServiceConfig{}
	err := validateConfig.SetKeys([][2]string{
		{cfg.Auth.AccessTokenKeyID, cfg.Auth.AccessTokenPublicKey},
	})
	if err != nil {
		return nil, err
	}

	s, err := chain.NewValidateService(validateConfig)
	if err != nil {
		return nil, err
	}

	return s, nil
}
