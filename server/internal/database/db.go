package database

import (
	"context"
	"database/sql"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gotodo.rasc.ch/internal/config"
)

func New(cfg config.Config) (*sql.DB, error) {
	connMaxIdleTime, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	connMaxLifetime, err := time.ParseDuration(cfg.DB.MaxLifetime)
	if err != nil {
		return nil, err
	}

	dbURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DB.User, cfg.DB.Password),
		Host:     cfg.DB.Connection,
		Path:     "/" + cfg.DB.Database,
		RawQuery: cfg.DB.Parameter,
	}
	db, err := sql.Open("pgx", dbURL.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetConnMaxLifetime(connMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
