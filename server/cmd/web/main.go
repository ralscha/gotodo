package main

import (
	"context"
	"database/sql"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/procyon-projects/chrono"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/exp/slog"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/database"
	"gotodo.rasc.ch/internal/mailer"
	"gotodo.rasc.ch/internal/version"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type application struct {
	config         *config.Config
	db             *sql.DB
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
	mailer         mailer.Mailer
	taskScheduler  chrono.TaskScheduler
}

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("reading config failed %v\n", err)
	}

	var logger *slog.Logger

	switch cfg.Environment {
	case config.Development:
		boil.DebugMode = true
		logger = slog.New(slog.NewTextHandler(os.Stdout))
	case config.Production:
		logger = slog.New(slog.NewJSONHandler(os.Stdout))
	}

	slog.SetDefault(logger)

	db, err := database.New(cfg)
	if err != nil {
		logger.Error("opening database connection failed", err)
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	logger.Info("database connection pool established")

	sm := scs.New()
	sm.Store = mysqlstore.NewWithCleanupInterval(db, 30*time.Minute)
	sm.Lifetime = cfg.Cleanup.SessionLifetime
	sm.Cookie.SameSite = http.SameSiteStrictMode
	if cfg.CookieDomain != "" {
		sm.Cookie.Domain = cfg.CookieDomain
	}
	sm.Cookie.Secure = cfg.SecureCookie
	logger.Info("secure cookie", "secure", sm.Cookie.Secure)

	err = initAuth(cfg)
	if err != nil {
		logger.Error("init auth failed", err)
		os.Exit(1)
	}

	app := &application{
		config:         &cfg,
		db:             db,
		sessionManager: sm,
		mailer:         mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
		taskScheduler:  chrono.NewDefaultTaskScheduler(),
	}

	_, err = app.taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		app.cleanup()
	}, 20*time.Minute)

	if err != nil {
		logger.Error("scheduling cleanup task failed", err)
		os.Exit(1)
	}

	logger.Info("starting server", "addr", app.config.HTTP.Port, "version", version.Get().Version)

	err = app.serve()
	if err != nil {
		logger.Error("http serve failed", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
