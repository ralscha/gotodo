package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/mailer"
	"log"
	"net/http"
	"reflect"
	"sync"
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
	validator      *validator.Validate
	decoder        *schema.Decoder
	wg             sync.WaitGroup
	logger         *zap.SugaredLogger
	mailer         mailer.Mailer
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("reading config failed %v\n", err)
	}

	var logger *zap.Logger

	switch cfg.Environment {
	case config.Development:
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalf("can't initialize development zap logger: %v\n", err)
		}
	case config.Production:
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("can't initialize production zap logger: %v\n", err)
		}
	}

	sugar := logger.Sugar()

	db, err := openDB(cfg)
	if err != nil {
		sugar.Fatalw("opening database connection failed", zap.Error(err))
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sugar.Info("database connection pool established")

	sm := scs.New()
	sm.Store = mysqlstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = 24 * time.Hour
	sm.Cookie.SameSite = http.SameSiteStrictMode

	sm.Cookie.Secure = cfg.SecureCookie
	sugar.Infof("secure cookie: %t\n", sm.Cookie.Secure)

	vld := validator.New()
	vld.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("name")
	})

	err = initAuth(cfg)
	if err != nil {
		sugar.Fatalw("init auth failed", zap.Error(err))
	}

	app := &application{
		config:         &cfg,
		db:             db,
		sessionManager: sm,
		validator:      vld,
		decoder:        schema.NewDecoder(),
		logger:         sugar,
		mailer:         mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}

	app.schedule(func() {
		err := app.deleteExpiredTokens()
		if err != nil {
			app.logger.Error("delete expired tokens failed", zap.Error(err))
		}
	}, time.Hour)

	err = app.serve()
	if err != nil {
		sugar.Fatalw("http serve failed", zap.Error(err))
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
