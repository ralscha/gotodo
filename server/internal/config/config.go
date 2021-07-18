package config

import (
	"github.com/spf13/viper"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	Environment  Environment
	SecureCookie bool
	BaseUrl      string
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
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	Argon2 struct {
		Memory      uint32
		Iterations  uint32
		Parallelism uint8
		SaltLength  uint32
		KeyLength   uint32
	}
}

func applyDefaults() {
	viper.SetDefault("environment", Production)
	viper.SetDefault("http.readTimeoutInSeconds", 30)
	viper.SetDefault("http.writeTimeoutInSeconds", 30)
	viper.SetDefault("http.idleTimeoutInSeconds", 120)
	viper.SetDefault("db.maxOpenConns", 25)
	viper.SetDefault("db.maxIdleConns", 25)
	viper.SetDefault("db.maxIdleTime", "15m")
	viper.SetDefault("secureCookie", true)
	viper.SetDefault("argon2.memory", 1<<17)
	viper.SetDefault("argon2.iterations", 20)
	viper.SetDefault("argon2.parallelism", 8)
	viper.SetDefault("argon2.saltLength", 16)
	viper.SetDefault("argon2.keyLength", 32)
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
