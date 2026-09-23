package transfer

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrInvalidAmount     = errors.New("transfer amount must be greater than zero")
	ErrSameAccount       = errors.New("source and destination accounts must be different")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type Service struct {
	db *sql.DB
}

type Request struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Transfer(ctx context.Context, req Request) error {
	if req.Amount <= 0 {
		return ErrInvalidAmount
	}

	if req.FromAccountID == req.ToAccountID {
		return ErrSameAccount
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	firstID := req.FromAccountID
	secondID := req.ToAccountID

	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}

	var firstBalance int64
	err = tx.QueryRowContext(
		ctx,
		`
		SELECT balance
		FROM accounts
		WHERE id = $1
		FOR UPDATE
		`,
		firstID,
	).Scan(&firstBalance)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	var secondBalance int64
	err = tx.QueryRowContext(
		ctx,
		`
		SELECT balance
		FROM accounts
		WHERE id = $1
		FOR UPDATE
		`,
		secondID,
	).Scan(&secondBalance)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrAccountNotFound
	}
	if err != nil {
		return err
	}

	var fromBalance int64

	if req.FromAccountID == firstID {
		fromBalance = firstBalance
	} else {
		fromBalance = secondBalance
	}

	if fromBalance < req.Amount {
		return ErrInsufficientFunds
	}

	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		`,
		req.Amount,
		req.FromAccountID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		`,
		req.Amount,
		req.ToAccountID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO transfers (
			from_account_id,
			to_account_id,
			amount,
			status
		)
		VALUES ($1, $2, $3, 'completed')
		`,
		req.FromAccountID,
		req.ToAccountID,
		req.Amount,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
