package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"time"
)

type scope string

const (
	scopeSignup        scope = "signup"
	scopePasswordReset scope = "password-reset"
	scopeEmailChange   scope = "email-change"
)

type token struct {
	plain string
	hash  []byte
}

func generateToken() (*token, error) {
	var token token
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.plain = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.plain))
	token.hash = hash[:]

	return &token, nil
}

func (app *application) insertToken(appUserId int64, ttl time.Duration, scope scope) (*token, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	newToken := models.Token{
		Hash:      token.hash,
		AppUserID: appUserId,
		Expiry:    time.Now().Add(ttl),
		Scope:     string(scope),
	}
	ctx, cancel := app.createDbContext()
	defer cancel()
	err = newToken.Insert(ctx, app.db, boil.Infer())

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (app *application) deleteAllTokensForUser(appUserId int64, scope scope) error {
	ctx, cancel := app.createDbContext()
	defer cancel()
	err := models.Tokens(models.TokenWhere.AppUserID.EQ(appUserId),
		models.TokenWhere.Scope.EQ(string(scope))).DeleteAll(ctx, app.db)
	return err
}

func (app *application) getAppUserIdFromToken(scope scope, tokenPlain string) (int64, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlain))

	ctx, cancel := app.createDbContext()
	defer cancel()
	token, err := models.Tokens(
		qm.Select(models.TokenColumns.AppUserID),
		models.TokenWhere.Hash.EQ(tokenHash[:]),
		models.TokenWhere.Scope.EQ(string(scope)),
		models.TokenWhere.Expiry.GT(time.Now())).One(ctx, app.db)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	if token != nil {
		return token.AppUserID, nil
	}
	return 0, nil

}
