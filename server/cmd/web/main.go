package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gotodo.rasc.ch/internal/config"
	"net/http"
	"time"
)

var (
	appBuildTime string
	appVersion   string
)

type application struct {
	config         *config.Config
	db             *sql.DB
	sessionManager *scs.SessionManager
}

func main() {
	// TODO: make this configurable
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("reading config failed")
	}

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("opening database connection failed")
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	log.Info().Msg("database connection pool established")

	sm := scs.New()
	sm.Store = mysqlstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = 24 * time.Hour
	sm.Cookie.SameSite = http.SameSiteStrictMode

	// TODO: make this configurable
	sm.Cookie.Secure = true

	app := &application{
		config:         &cfg,
		db:             db,
		sessionManager: sm,
	}

	err = app.serve()
	if err != nil {
		log.Fatal().Err(err).Msg("http serve failed")
	}

}

func openDB(cfg config.Config) (*sql.DB, error) {
	dbstring := fmt.Sprintf("%s:%s@%s/%s?%s",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Connection, cfg.Db.Database, cfg.Db.Parameter)

	db, err := sql.Open("mysql", dbstring)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
