package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gotodo.rasc.ch/internal/config"
)

func New(cfg config.Config) (*sql.DB, error) {
	dbstring := fmt.Sprintf("postgres://%s:%s@%s/%s?%s",
		url.QueryEscape(cfg.DB.User), url.QueryEscape(cfg.DB.Password), cfg.DB.Connection, cfg.DB.Database, cfg.DB.Parameter)

	db, err := sql.Open("pgx", dbstring)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	connMaxIdleTime, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(connMaxIdleTime)

	connMaxLifetime, err := time.ParseDuration(cfg.DB.MaxLifetime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
