package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dbreliant/internal/db"
	httpapi "dbreliant/internal/http"
	"dbreliant/internal/metrics"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	conn, err := db.Open(ctx, db.LoadConfigFromEnv())
	if err != nil {
		log.Fatalf("connect to postgres: %v", err)
	}

	metrics.RegisterDBMetrics(conn)
	metrics.RegisterHTTPMetrics()

	stats := conn.Stats()
	log.Printf(
		"db pool configured: max_open=%d open=%d idle=%d in_use=%d",
		stats.MaxOpenConnections,
		stats.OpenConnections,
		stats.Idle,
		stats.InUse,
	)

	defer func() {
		if err := db.Close(conn); err != nil {
			log.Printf("close postgres connection: %v", err)
		}
	}()

	port := envOrDefault("HTTP_PORT", "8080")
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           metrics.HTTPMiddleware(httpapi.NewMux()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("DBReliant API listening on :%s", port)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown http server: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
