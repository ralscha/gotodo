package main

import (
	"flag"
	"fmt"
	"gotodo.rasc.ch/internal/config"
	"log"
	"os"

	"github.com/pressly/goose"

	_ "github.com/go-sql-driver/mysql"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", "./migrations", "directory with migration files")
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
		cfg.Db.User, cfg.Db.Password, cfg.Db.Connection, cfg.Db.Database, cfg.Db.Parameter)

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

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
