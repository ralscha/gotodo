package main

import (
	"context"
	"database/sql"

	"github.com/alexedwards/argon2id"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInitialUsers, downInitialUsers)
}

func upInitialUsers(ctx context.Context, tx *sql.Tx) error {
	params := &argon2id.Params{
		Memory:      1 << 17,
		Iterations:  20,
		Parallelism: 8,
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := argon2id.CreateHash("password", params)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO app_user (email, password_hash, authority, activated, expired, last_access)
	VALUES ('admin@test.ch', $1, 'admin', true, null, null)
    `
	_, err = tx.ExecContext(ctx, stmt, hash)
	if err != nil {
		return err
	}
	return nil
}

func downInitialUsers(ctx context.Context, tx *sql.Tx) error {
	stmt := `DELETE FROM app_user where email = 'admin@test.ch'`
	_, err := tx.ExecContext(ctx, stmt)
	if err != nil {
		return err
	}
	return nil
}
