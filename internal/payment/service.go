package payment

import (
	"context"
	"database/sql"
	"time"
)

type Service struct {
	db *sql.DB
}

type Payment struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ListRequest struct {
	AccountID int64
	Status    string
	Limit     int
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) List(ctx context.Context, req ListRequest) ([]Payment, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id, account_id, amount, status, created_at
		FROM payments
		WHERE account_id = $1
		  AND status = $2
		ORDER BY created_at DESC
		LIMIT $3
		`,
		req.AccountID,
		req.Status,
		req.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]Payment, 0, req.Limit)

	for rows.Next() {
		var payment Payment

		if err := rows.Scan(
			&payment.ID,
			&payment.AccountID,
			&payment.Amount,
			&payment.Status,
			&payment.CreatedAt,
		); err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
