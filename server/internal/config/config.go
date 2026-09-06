package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

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
	CookieDomain string
	BaseURL      string
	Cleanup      struct {
		InactiveUsersMaxAge      time.Duration
		ExpiredUsersMaxAge       time.Duration
		EmailChangeTokenMaxAge   time.Duration
		SignupTokenMaxAge        time.Duration
		PasswordResetTokenMaxAge time.Duration
		SessionLifetime          time.Duration
	}
	DB struct {
		User         string
		Password     string
		Connection   string
		Database     string
		Parameter    string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
		MaxLifetime  string
	}
	HTTP struct {
		Port                  string
		ReadTimeoutInSeconds  int64
		WriteTimeoutInSeconds int64
		IdleTimeoutInSeconds  int64
	}
	SMTP struct {
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

func applyDefaults(v *viper.Viper) {
	v.SetDefault("environment", Production)
	v.SetDefault("http.readTimeoutInSeconds", 30)
	v.SetDefault("http.writeTimeoutInSeconds", 30)
	v.SetDefault("http.idleTimeoutInSeconds", 120)
	v.SetDefault("db.maxOpenConns", 4)
	v.SetDefault("db.maxIdleConns", 2)
	v.SetDefault("db.maxIdleTime", "15m")
	v.SetDefault("db.maxLifetime", "2h")
	v.SetDefault("secureCookie", true)
	v.SetDefault("argon2.memory", 1<<17)
	v.SetDefault("argon2.iterations", 20)
	v.SetDefault("argon2.parallelism", 8)
	v.SetDefault("argon2.saltLength", 16)
	v.SetDefault("argon2.keyLength", 32)
	v.SetDefault("cleanup.inactiveUsersMaxAge", "8760h")
	v.SetDefault("cleanup.expiredUsersMaxAge", "8760h")
	v.SetDefault("cleanup.emailChangeTokenMaxAge", "48h")
	v.SetDefault("cleanup.signupTokenMaxAge", "48h")
	v.SetDefault("cleanup.passwordResetTokenMaxAge", "2h")
	v.SetDefault("cleanup.sessionLifetime", "720h")
}

func LoadConfig() (Config, error) {
	return loadConfig(".")
}

func loadConfig(configPath string) (Config, error) {
	var cfg Config

	v := viper.New()
	applyDefaults(v)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AddConfigPath(configPath)
	v.SetEnvPrefix("GOTODO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	for _, key := range configKeys {
		if err := v.BindEnv(key); err != nil {
			return cfg, fmt.Errorf("bind environment variable for %s: %w", key, err)
		}
	}

	err := v.ReadInConfig()
	var notFound viper.ConfigFileNotFoundError
	if err != nil && !errors.As(err, &notFound) {
		return cfg, err
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	if err := cfg.validate(); err != nil {
		return cfg, err
	}
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL += "/"
	}

	return cfg, nil
}

var configKeys = []string{
	"environment",
	"secureCookie",
	"cookieDomain",
	"baseURL",
	"db.user",
	"db.password",
	"db.connection",
	"db.database",
	"db.parameter",
	"db.maxOpenConns",
	"db.maxIdleConns",
	"db.maxIdleTime",
	"db.maxLifetime",
	"http.port",
	"http.readTimeoutInSeconds",
	"http.writeTimeoutInSeconds",
	"http.idleTimeoutInSeconds",
	"smtp.host",
	"smtp.port",
	"smtp.username",
	"smtp.password",
	"smtp.sender",
	"argon2.memory",
	"argon2.iterations",
	"argon2.parallelism",
	"argon2.saltLength",
	"argon2.keyLength",
	"cleanup.inactiveUsersMaxAge",
	"cleanup.expiredUsersMaxAge",
	"cleanup.emailChangeTokenMaxAge",
	"cleanup.signupTokenMaxAge",
	"cleanup.passwordResetTokenMaxAge",
	"cleanup.sessionLifetime",
}

func (cfg Config) validate() error {
	if cfg.Environment != Development && cfg.Environment != Production {
		return fmt.Errorf("environment must be %q or %q", Development, Production)
	}
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil || (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		return errors.New("baseURL must be an absolute HTTP or HTTPS URL")
	}

	required := map[string]string{
		"db.user":       cfg.DB.User,
		"db.connection": cfg.DB.Connection,
		"db.database":   cfg.DB.Database,
		"http.port":     cfg.HTTP.Port,
		"smtp.host":     cfg.SMTP.Host,
		"smtp.sender":   cfg.SMTP.Sender,
	}
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", key)
		}
	}
	if cfg.SMTP.Port <= 0 || cfg.SMTP.Port > 65535 {
		return errors.New("smtp.port must be between 1 and 65535")
	}
	return nil
}
