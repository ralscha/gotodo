package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Db struct {
		User       string
		Password   string
		Connection string
		Database   string
	}
	Http struct {
		Port                  string
		ReadTimeoutInSeconds  time.Duration
		WriteTimeoutInSeconds time.Duration
		IdleTimeoutInSeconds  time.Duration
	}
}

func applyDefaults() {
	viper.SetDefault("http.readTimeoutInSeconds", 30)
	viper.SetDefault("http.writeTimeoutInSeconds", 30)
	viper.SetDefault("http.idleTimeoutInSeconds", 120)
}

func LoadConfig() (Config, error) {
	var cfg Config

	applyDefaults()
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
