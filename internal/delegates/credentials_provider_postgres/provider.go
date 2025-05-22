package credentials_provider_postgres

import (
	"context"

	"github.com/pot876/sauth/internal/chain"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Provider struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) chain.ICredentialsProvider {
	return &Provider{
		conn: conn,
	}
}

func (p *Provider) ValidateAuth(ctx context.Context, realmID uuid.UUID, login, password []byte) (any, error) {
	if len(login) == 0 {
		return nil, chain.ErrNotFound
	}

	authInfo, err := AuthInfoByLogin(ctx, p.conn, realmID, string(login))
	if err != nil {
		if err == ErrNotFound {
			return nil, chain.ErrNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(authInfo.Pwdhash, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, chain.ErrBadPassword
		}
		return nil, err
	}

	return authInfo, nil
}

func (a *Provider) RegisterMetrics(r prometheus.Registerer) {

}
