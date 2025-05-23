package credentials_provider_file

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pot876/sauth/internal/chain"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/bcrypt"
)

type Provider struct {
	users map[string]FileUserInfo

	metricsBcryptTimings prometheus.Histogram
}

func New(users []*FileUserInfo) (*Provider, error) {
	usersMap := make(map[string]FileUserInfo, len(users))
	for _, u := range users {
		usersMap[u.RealmID+":"+u.Login] = FileUserInfo{
			RealmID: u.RealmID,
			UserID:  u.UserID,
			Role:    u.Role,
			Login:   u.Login,
			Pwdhash: u.Pwdhash,
		}
	}

	return &Provider{
		users: usersMap,
	}, nil
}

func (p *Provider) ValidateAuth(ctx context.Context, realmID uuid.UUID, login, password []byte) (any, error) {
	if len(login) == 0 {
		return nil, chain.ErrNotFound
	}

	authInfo, ok := p.users[realmID.String()+":"+string(login)]
	if !ok {
		return nil, chain.ErrNotFound
	}

	err := p.compareHashAndPassword([]byte(authInfo.Pwdhash), password)
	if err != nil {
		return nil, err
	}

	authInfo.Login = ""
	authInfo.Pwdhash = ""
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
