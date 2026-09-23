package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"dbreliant/internal/payment"
)

const (
	defaultPaymentLimit = 100
	maxPaymentLimit     = 1000
)

func paymentHandler(service *payment.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		accountID, err := strconv.ParseInt(
			r.URL.Query().Get("account_id"),
			10,
			64,
		)
		if err != nil || accountID <= 0 {
			http.Error(w, "invalid account_id", http.StatusBadRequest)
			return
		}

		status := r.URL.Query().Get("status")
		if !validPaymentStatus(status) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}

		limit := defaultPaymentLimit

		if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
			parsedLimit, err := strconv.Atoi(rawLimit)
			if err != nil || parsedLimit <= 0 || parsedLimit > maxPaymentLimit {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}

			limit = parsedLimit
		}

		payments, err := service.List(
			r.Context(),
			payment.ListRequest{
				AccountID: accountID,
				Status:    status,
				Limit:     limit,
			},
		)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"payments": payments,
			"count":    len(payments),
		})
	}
}

func validPaymentStatus(status string) bool {
	switch status {
	case "processing", "completed", "failed":
		return true
	default:
		return false
	}
}
