package config

import (
	"github.com/spf13/viper"
)

type LogLevel string

const (
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
)

type Config struct {
	LogLevel     LogLevel
	SecureCookie bool
	Db           struct {
		User         string
		Password     string
		Connection   string
		Database     string
		Parameter    string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Http struct {
		Port                  string
		ReadTimeoutInSeconds  int64
		WriteTimeoutInSeconds int64
		IdleTimeoutInSeconds  int64
	}
}

func applyDefaults() {
	viper.SetDefault("http.readTimeoutInSeconds", 30)
	viper.SetDefault("http.writeTimeoutInSeconds", 30)
	viper.SetDefault("http.idleTimeoutInSeconds", 120)
	viper.SetDefault("db.maxOpenConns", 25)
	viper.SetDefault("db.maxIdleConns", 25)
	viper.SetDefault("db.maxIdleTime", "15m")
	viper.SetDefault("secureCookie", true)
	viper.SetDefault("LogLevel", Info)
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
