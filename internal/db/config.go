package db

import (
	"net/url"
	"os"
)

type Config struct {
	DatabaseURL string
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string
}

func LoadConfigFromEnv() Config {
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Host:        envOrDefault("POSTGRES_HOST", "localhost"),
		Port:        envOrDefault("POSTGRES_PORT", "5432"),
		User:        envOrDefault("POSTGRES_USER", "dbreliant"),
		Password:    envOrDefault("POSTGRES_PASSWORD", "dbreliant"),
		Name:        envOrDefault("POSTGRES_DB", "dbreliant"),
		SSLMode:     envOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

func (c Config) ConnString() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + c.Port,
		Path:   c.Name,
	}

	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
