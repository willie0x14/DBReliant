package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	conn, err := sql.Open("pgx", cfg.ConnString())
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	conn.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := Ping(ctx, conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

func Ping(ctx context.Context, conn *sql.DB) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return conn.PingContext(pingCtx)
}

func Close(conn *sql.DB) error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}
