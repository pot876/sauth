package credentials_provider_postgres

import (
	"context"
	"time"

	"github.com/pot876/sauth/internal/chain"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Provider struct {
	conn *pgxpool.Pool

	metricsBcryptTimings prometheus.Histogram
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

	err = p.compareHashAndPassword([]byte(authInfo.Pwdhash), password)
	if err != nil {
		return nil, err
	}

	authInfo.Pwdhash = nil
	return authInfo, nil
}

func (p *Provider) compareHashAndPassword(hashedPassword, password []byte) error {
	t0 := time.Now()
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	p.metricsBcryptTimings.Observe(float64(time.Since(t0).Seconds()))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return chain.ErrBadPassword
		}
		return err
	}

	return nil
}

func (a *Provider) RegisterMetrics(r prometheus.Registerer, prefix string) {
	a.metricsBcryptTimings = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    prefix + "bcrypt_timings",
		Help:    "",
		Buckets: []float64{.1, .2, .3, .5, 1., 2.},
	})

	r.MustRegister(a.metricsBcryptTimings)
}
