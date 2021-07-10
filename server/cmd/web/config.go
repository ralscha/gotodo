package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Db struct {
		Url      string
		User     string
		Password string
	}
	Http struct {
		Port string
	}
}

func loadConfig() (Config, error) {
	var cfg Config

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	viper.SetEnvPrefix("GOTODO")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
