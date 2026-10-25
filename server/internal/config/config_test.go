package config

import (
	"strings"
	"testing"
)

func TestLoadConfigFromEnvironment(t *testing.T) {
	t.Setenv("GOTODO_ENVIRONMENT", "development")
	t.Setenv("GOTODO_BASEURL", "https://example.test/app")
	t.Setenv("GOTODO_DB_USER", "user")
	t.Setenv("GOTODO_DB_PASSWORD", "password")
	t.Setenv("GOTODO_DB_CONNECTION", "db.example.test:5432")
	t.Setenv("GOTODO_DB_DATABASE", "gotodo")
	t.Setenv("GOTODO_DB_PARAMETER", "sslmode=require")
	t.Setenv("GOTODO_HTTP_PORT", ":8080")
	t.Setenv("GOTODO_HTTP_TRUSTEDPROXYCIDRS", "10.1.0.4/32,::1/128")
	t.Setenv("GOTODO_SMTP_HOST", "smtp.example.test")
	t.Setenv("GOTODO_SMTP_PORT", "2525")
	t.Setenv("GOTODO_SMTP_SENDER", "no-reply@example.test")

	cfg, err := loadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}
	if cfg.BaseURL != "https://example.test/app/" {
		t.Fatalf("BaseURL = %q", cfg.BaseURL)
	}
	if cfg.DB.Connection != "db.example.test:5432" {
		t.Fatalf("DB.Connection = %q", cfg.DB.Connection)
	}
	if len(cfg.HTTP.TrustedProxyCIDRs) != 2 || cfg.HTTP.TrustedProxyCIDRs[0] != "10.1.0.4/32" || cfg.HTTP.TrustedProxyCIDRs[1] != "::1/128" {
		t.Fatalf("TrustedProxyCIDRs = %v", cfg.HTTP.TrustedProxyCIDRs)
	}
	if cfg.Cleanup.SessionLifetime == 0 {
		t.Fatal("default session lifetime was not applied")
	}
}

func TestLoadConfigRejectsMissingRequiredValues(t *testing.T) {
	_, err := loadConfig(t.TempDir())
	if err == nil {
		t.Fatal("loadConfig() error = nil")
	}
	if !strings.Contains(err.Error(), "baseURL") {
		t.Fatalf("loadConfig() error = %v, want baseURL error", err)
	}
}

func TestLoadConfigRejectsInvalidEnvironment(t *testing.T) {
	t.Setenv("GOTODO_ENVIRONMENT", "staging")
	t.Setenv("GOTODO_BASEURL", "https://example.test")

	_, err := loadConfig(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "environment") {
		t.Fatalf("loadConfig() error = %v, want environment error", err)
	}
}

func TestConfigRejectsInvalidTrustedProxy(t *testing.T) {
	var cfg Config
	cfg.HTTP.TrustedProxyCIDRs = []string{"not-a-cidr"}
	err := cfg.validate()
	if err == nil || !strings.Contains(err.Error(), "trustedProxyCIDRs") {
		t.Fatalf("validate() error = %v, want invalid proxy CIDR", err)
	}
}
