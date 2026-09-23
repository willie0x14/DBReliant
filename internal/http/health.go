package httpapi

import (
	"encoding/json"
	"net/http"

	"dbreliant/internal/payment"
	"dbreliant/internal/transfer"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMux(
	transferService *transfer.Service,
	paymentService *payment.Service,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/transfers", transferHandler(transferService))
	mux.HandleFunc("/payments", paymentHandler(paymentService))

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
