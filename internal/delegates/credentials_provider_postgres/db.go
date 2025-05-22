package credentials_provider_postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

var DatabaseRequestTimeout = 10 * time.Second

var (
	ErrNotFound = errors.New("not found")
)

var errNoRows = pgx.ErrNoRows

type DB interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func logerror(ctx context.Context, err error) {
	logger := zerolog.Ctx(ctx)
	logger.Error().Err(err).Msg("")
}
