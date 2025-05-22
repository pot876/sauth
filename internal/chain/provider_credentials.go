package chain

import (
	"context"

	"github.com/google/uuid"
)

type ICredentialsProvider interface {
	ValidateAuth(ctx context.Context, realmID uuid.UUID, login, password []byte) (any, error)
}
