package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		return
	}
	command := args[0]

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("reading config failed", err)
	}

	dbstring := fmt.Sprintf("postgres://%s:%s@%s/%s?%s",
		url.QueryEscape(cfg.DB.User), url.QueryEscape(cfg.DB.Password), cfg.DB.Connection, cfg.DB.Database, cfg.DB.Parameter)

	db, err := goose.OpenDBWithDriver("pgx", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	goose.SetBaseFS(migrations.EmbeddedFiles)

	if err := goose.Run(command, db, ".", arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
