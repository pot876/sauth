package credentials_provider_postgres

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthInfo struct {
	RealmID uuid.UUID `json:"realm_id"`
	UserID  uuid.UUID `json:"user_id"`
	Role    string    `json:"role,omitempty"`
	Pwdhash []byte    `json:"pwdhash,omitempty"`
}

func AuthInfoByLogin(ctx context.Context, db DB, realmID uuid.UUID, login string) (*UserAuthInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	const sqlstr = `select ` +
		`user_id, pwdhash, role ` +
		`from public.credentials ` +
		`where realm_id = $1 and login = $2`

	dbRequestCtx, dbRequestCtxCancel := context.WithTimeout(context.Background(), DatabaseRequestTimeout)
	defer dbRequestCtxCancel()

	var authInfo UserAuthInfo
	err := db.QueryRow(dbRequestCtx, sqlstr, realmID, login).Scan(&authInfo.UserID, &authInfo.Pwdhash, &authInfo.Role)
	if err != nil {
		if err == errNoRows {
			return nil, ErrNotFound
		}

		logerror(ctx, err)
		return nil, err
	}

	authInfo.RealmID = realmID
	return &authInfo, nil
}

func exampleCreateUser(ctx context.Context, db DB, realmID uuid.UUID, userID uuid.UUID, login, password, role string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	const sqlstr = `insert into public.credentials (realm_id, user_id, login, pwdhash, role) values ($1, $2, $3, $4, $5)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	dbRequestCtx, dbRequestCtxCancel := context.WithTimeout(context.Background(), DatabaseRequestTimeout)
	defer dbRequestCtxCancel()

	_, err = db.Exec(dbRequestCtx, sqlstr, realmID, userID, login, hashedPassword, role)
	if err != nil {
		logerror(ctx, err)
		return err
	}

	return nil
}
