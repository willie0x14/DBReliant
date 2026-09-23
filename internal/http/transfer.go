package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"dbreliant/internal/transfer"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type transferResponse struct {
	Status string `json:"status"`
}

func transferHandler(service *transfer.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req transferRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err := service.Transfer(r.Context(), transfer.Request{
			FromAccountID: req.FromAccountID,
			ToAccountID:   req.ToAccountID,
			Amount:        req.Amount,
		})

		switch {
		case err == nil:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			_ = json.NewEncoder(w).Encode(transferResponse{
				Status: "completed",
			})

		case errors.Is(err, transfer.ErrInvalidAmount):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, transfer.ErrSameAccount):
			http.Error(w, err.Error(), http.StatusBadRequest)

		case errors.Is(err, transfer.ErrAccountNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)

		case errors.Is(err, transfer.ErrInsufficientFunds):
			http.Error(w, err.Error(), http.StatusConflict)

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
