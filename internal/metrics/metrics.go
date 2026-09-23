package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

func RegisterDBMetrics(conn *sql.DB) {
	prometheus.MustRegister(
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "db_max_open_connections",
				Help: "Maximum number of open database connections.",
			},
			func() float64 {
				return float64(conn.Stats().MaxOpenConnections)
			},
		),
	)

	prometheus.MustRegister(
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "db_open_connections",
				Help: "Current number of open database connections.",
			},
			func() float64 {
				return float64(conn.Stats().OpenConnections)
			},
		),
	)

	prometheus.MustRegister(
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "db_in_use_connections",
				Help: "Current number of database connections in use.",
			},
			func() float64 {
				return float64(conn.Stats().InUse)
			},
		),
	)

	prometheus.MustRegister(
		prometheus.NewGaugeFunc(
			prometheus.GaugeOpts{
				Name: "db_idle_connections",
				Help: "Current number of idle database connections.",
			},
			func() float64 {
				return float64(conn.Stats().Idle)
			},
		),
	)

	prometheus.MustRegister(
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name: "db_wait_count_total",
				Help: "Total number of waits for a database connection.",
			},
			func() float64 {
				return float64(conn.Stats().WaitCount)
			},
		),
	)

	prometheus.MustRegister(
		prometheus.NewCounterFunc(
			prometheus.CounterOpts{
				Name: "db_wait_duration_seconds_total",
				Help: "Total time spent waiting for database connections.",
			},
			func() float64 {
				return conn.Stats().WaitDuration.Seconds()
			},
		),
	)
}
