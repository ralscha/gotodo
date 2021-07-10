package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexedwards/argon2id"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/internal/models"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("reading config failed")
	}

	dbstring := fmt.Sprintf("%s:%s@%s/%s?%s",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Connection, cfg.Db.Database, cfg.Db.Parameter)

	// Open handle to database like normal
	db, err := sql.Open("mysql", dbstring)
	if err != nil {
		log.Fatal().Err(err).Msg("opening database failed")
	}

	// If you don't want to pass in db to all generated methods
	// you can use boil.SetDB to set it globally, and then use
	// the G variant methods like so (--add-global-variants to enable)
	ctx := context.Background()
	users, err := models.AppUsers().All(ctx, db)
	if err != nil {
		log.Fatal().Err(err).Msg("select app_users")
	}
	fmt.Println(users)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// UNIX Time is faster and smaller than most timestamps
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Print("hello world")
	log.Warn().Msg("a warning")

	// CreateHash returns a Argon2id hash of a plain-text password using the
	// provided algorithm parameters. The returned hash follows the format used
	// by the Argon2 reference C implementation and looks like this:
	// $argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
	hash, err := argon2id.CreateHash("pa$$word", argon2id.DefaultParams)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
	fmt.Println(hash)
	// ComparePasswordAndHash performs a constant-time comparison between a
	// plain-text password and Argon2id hash, using the parameters and salt
	// contained in the hash. It returns true if they match, otherwise it returns
	// false.
	match, err := argon2id.ComparePasswordAndHash("pa$$word", hash)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	log.Printf("Match: %v", match)
}
