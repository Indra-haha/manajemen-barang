package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port         string
	JWTSecret    string
	CookieSecure bool
	Database     DatabaseConfig
}

type DatabaseConfig struct {
	URI     string
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

func (d DatabaseConfig) DSN() string {
	if d.URI != "" {
		return d.URI
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Pass, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

func Load() Config {
	return Config{
		Port:         getenv("PORT", "8080"),
		JWTSecret:    getenv("JWT_SECRET", ""),
		CookieSecure: getenvBool("COOKIE_SECURE", false),
		Database:     LoadDatabase(),
	}
}

func getenvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

// LoadDatabase membaca konfigurasi Postgres dari environment beserta nilai
// default yang aman untuk pengembangan lokal.
func LoadDatabase() DatabaseConfig {
	return DatabaseConfig{
		URI:     getenv("POSTGRES_URI", ""),
		Host:    getenv("DB_HOST", ""),
		Port:    getenv("DB_PORT", ""),
		User:    getenv("DB_USER", ""),
		Pass:    getenv("DB_PASS", ""),
		Name:    getenv("DB_NAME", ""),
		SSLMode: getenv("DB_SSLMODE", ""),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
