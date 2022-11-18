package main

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"gotodo.rasc.ch/internal/config"
	"gotodo.rasc.ch/migrations"
	"log"
	"os"
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

	dbstring := fmt.Sprintf("%s:%s@%s/%s?%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Connection, cfg.DB.Database, cfg.DB.Parameter)

	db, err := goose.OpenDBWithDriver("mysql", dbstring)
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
