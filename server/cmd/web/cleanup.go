package main

import (
	"context"
	"github.com/volatiletech/null/v8"
	"go.uber.org/zap"
	"gotodo.rasc.ch/internal/models"
	"time"
)

func (app *application) cleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokens, err := models.Tokens(models.TokenWhere.Expiry.LT(time.Now())).All(ctx, app.db)
	if err != nil {
		app.logger.Error("deleting expired tokens failed", zap.Error(err))
	}

	for _, token := range tokens {
		switch scope(token.Scope) {
		case scopeSignup:
			// Delete all users that created a registration but never confirmed it
			err := models.AppUsers(models.AppUserWhere.ID.EQ(token.AppUserID)).DeleteAll(ctx, app.db)
			if err != nil {
				app.logger.Error("deleting user failed", zap.Error(err))
			}
		case scopeEmailChange:
			// Reset all email change requests where the confirmation token is expired
			err := models.AppUsers(models.AppUserWhere.ID.EQ(token.AppUserID)).UpdateAll(ctx, app.db,
				models.M{models.AppUserColumns.EmailNew: null.NewString("", false)})
			if err != nil {
				app.logger.Error("updating user failed", zap.Error(err))
			}
		}

		err := token.Delete(ctx, app.db)
		if err != nil {
			app.logger.Error("deleting token failed", zap.Error(err))
		}
	}

	// Delete all users that are expired for the configured amount of time
	expired := time.Now().Add(-app.config.Cleanup.ExpiredUsersMaxAge)
	err = models.AppUsers(models.AppUserWhere.Expired.LT(null.NewTime(expired, true))).DeleteAll(ctx, app.db)
	if err != nil {
		app.logger.Error("deleting expired users failed", zap.Error(err))
	}

	// Inactivate all users where the last access was older than the configured max age
	inactive := time.Now().Add(-app.config.Cleanup.InactiveUsersMaxAge)
	err = models.AppUsers(models.AppUserWhere.LastAccess.LT(null.NewTime(inactive, true))).UpdateAll(ctx, app.db,
		models.M{models.AppUserColumns.Expired: null.NewTime(time.Now(), true)})
	if err != nil {
		app.logger.Error("updating user failed", zap.Error(err))
	}

}
