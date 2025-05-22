package credentials_provider_postgres

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// PG_TEST_CREATE_USERS=1 go test -v -count=1 ./internal/delegates/credentials_provider_postgres
func TestCreateUsers(t *testing.T) {
	if os.Getenv("PG_TEST_CREATE_USERS") == "" {
		t.SkipNow()
	}

	cfg, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable")
	require.NoError(t, err)
	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)

	err = pool.Ping(ctx)
	require.NoError(t, err)

	for range 10 {
		err := exampleCreateUser(ctx, pool, uuid.New(), uuid.New(), "test", "test", "test")
		require.NoError(t, err)
	}
}
