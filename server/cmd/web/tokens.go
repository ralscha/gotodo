package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"gotodo.rasc.ch/internal/models"
	"time"
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

func (app *application) insertToken(ctx context.Context, appUserID int64, ttl time.Duration, scope models.TokensScope) (*token, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	newToken := models.Token{
		Hash:      token.hash,
		AppUserID: appUserID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}
	err = newToken.Insert(ctx, app.db, boil.Infer())

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (app *application) deleteAllTokensForUser(ctx context.Context, appUserID int64, scope models.TokensScope) error {
	err := models.Tokens(models.TokenWhere.AppUserID.EQ(appUserID),
		models.TokenWhere.Scope.EQ(scope)).DeleteAll(ctx, app.db)
	return err
}

func (app *application) getAppUserIDFromToken(ctx context.Context, scope models.TokensScope, tokenPlain string) (int64, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlain))

	token, err := models.Tokens(
		qm.Select(models.TokenColumns.AppUserID),
		models.TokenWhere.Hash.EQ(tokenHash[:]),
		models.TokenWhere.Scope.EQ(scope),
		models.TokenWhere.Expiry.GT(time.Now())).One(ctx, app.db)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	if token != nil {
		return token.AppUserID, nil
	}
	return 0, nil

}
