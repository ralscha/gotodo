package main

import (
	"database/sql"
	"github.com/alexedwards/argon2id"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upInitialUsers, downInitialUsers)
}

func upInitialUsers(tx *sql.Tx) error {
	hash, err := argon2id.CreateHash("password", argon2id.DefaultParams)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO app_user (email, password_hash, authority, activated, expired, last_access) 
	VALUES ('admin@test.ch', ?, 'ADMIN', 1, null, null)
    `
	_, err = tx.Exec(stmt, hash)
	if err != nil {
		return err
	}
	return nil
}

func downInitialUsers(tx *sql.Tx) error {
	stmt := `DELETE FROM app_user where email = 'admin@test.ch'`
	_, err := tx.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
