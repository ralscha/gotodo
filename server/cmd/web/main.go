package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gotodo.rasc.ch/internal/config"
)

var (
	appBuildTime string
	appVersion   string
)

type application struct {
	config config.Config
}

func main() {
	//TODO configurable
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("reading config failed")
	}

	app := &application{
		config: cfg,
	}

	err = app.serve()
	if err != nil {
		log.Fatal().Err(err).Msg("http serve failed")
	}

}
