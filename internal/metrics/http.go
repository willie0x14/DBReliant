package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}

	r.statusCode = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration in seconds.",
			Buckets: []float64{
				0.00005, // 0.05 ms
				0.0001,  // 0.1 ms
				0.00025, // 0.25 ms
				0.0005,  // 0.5 ms
				0.001,   // 1 ms
				0.0025,  // 2.5 ms
				0.005,   // 5 ms
				0.01,    // 10 ms
				0.025,   // 25 ms
				0.05,    // 50 ms
				0.1,     // 100 ms
				0.25,    // 250 ms
				0.5,     // 500 ms
				1.0,     // 1 s
				2.5,     // 2.5 s
				5.0,     // 5 s
			},
			// Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func RegisterHTTPMetrics() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
	)
}

func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(
			r.Method,
			strconv.Itoa(recorder.statusCode),
		).Inc()

		httpRequestDuration.WithLabelValues(
			r.Method,
		).Observe(duration)
	})
}
