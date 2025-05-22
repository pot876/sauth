package config

import "time"

type Config struct {
	HTTPListenAddr       string `env:"HTTP_LISTEN_ADDR"`
	HTTPEndpointPrefix   string `env:"HTTP_ENDPOINT_PREFIX"`
	HTTPEndpointLogin    string `env:"HTTP_ENDPOINT_LOGIN,default=/login"`
	HTTPEndpointRefresh  string `env:"HTTP_ENDPOINT_REFRESH,default=/refresh"`
	HTTPEndpointInfo     string `env:"HTTP_ENDPOINT_LOGIN,default=/info"`
	HTTPEndpointValidate string `env:"HTTP_ENDPOINT_REFRESH,default=/validate"`

	// prefix not applied to metrics endpoint !
	HTTPEndpointMetrics string `env:"HTTP_ENDPOINT_METRICS,default=/metrics"`

	Auth Auth `env:",prefix=AUTH_"`

	CredentialsProviderPostgresEnabled bool                        `env:"CP_PG_ENABLED"`
	CredentialsProviderPostgres        CredentialsProviderPostgres `env:",prefix=CP_PG_"`

	CredentialsProviderFileEnabled bool                    `env:"CP_FILE_ENABLED"`
	CredentialsProviderFile        CredentialsProviderFile `env:",prefix=CP_FILE_"`
}

type Auth struct {
	AccessTokenKeyID      string `env:"ACCESS_TOKEN_KEY_ID"`
	AccessTokenPrivateKey string `env:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string `env:"ACCESS_TOKEN_PUBLIC_KEY"`

	RefreshTokenKeyID      string `env:"REFRESH_TOKEN_KEY_ID"`
	RefreshTokenPrivateKey string `env:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string `env:"REFRESH_TOKEN_PUBLIC_KEY"`

	AccessTokenExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRED_IN,default=6h"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRED_IN,default=72h"`
}

type CredentialsProviderPostgres struct {
	ConnectionString string `env:"CONNECTION_STRING"`
}
type CredentialsProviderFile struct {
	Format string `env:"FORMAT,default=json"`
	Path   string `env:"PATH"`
}
