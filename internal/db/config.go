package db

import (
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL string
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	SSLMode     string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
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

		MaxOpenConns:    envInt("DB_MAX_OPEN_CONNS", 10),
		MaxIdleConns:    envInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		ConnMaxIdleTime: envDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute),
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

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}
