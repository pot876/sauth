package credentials_provider_file

import (
	"context"

	"github.com/google/uuid"
	"github.com/pot876/sauth/internal/chain"
	"golang.org/x/crypto/bcrypt"
)

type Provider struct {
	users map[string]FileUserInfo
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

	err := bcrypt.CompareHashAndPassword([]byte(authInfo.Pwdhash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, chain.ErrBadPassword
		}
		return nil, err
	}

	return authInfo, nil
}
