package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository/postgres"
)

type CreateTransactionParams struct {
	UserId      int32
	Value       int64
	Type        TransactionType
	Description string
}

type TransactionCreated struct {
	Limit   int64 `json:"limite"`
	Balance int64 `json:"saldo"`
}

type TransactionService interface {
	Create(ctx context.Context, params CreateTransactionParams) (*TransactionCreated, error)
}

type transactionService struct {
	dbpool  *pgxpool.Pool
	queries *postgres.Queries
}

func NewTransactionService(dbpool *pgxpool.Pool) TransactionService {
	return &transactionService{
		dbpool: dbpool,
	}
}

func (s *transactionService) Create(ctx context.Context, params CreateTransactionParams) (*TransactionCreated, error) {

	tx, err := s.dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := postgres.New(s.dbpool)
	qtx := queries.WithTx(tx)

	var operationValue = params.Value

	if params.Type == Debit {
		operationValue *= -1
	}

	u, err := qtx.UpdateUserBalance(ctx, postgres.UpdateUserBalanceParams{
		Balance: operationValue,
		ID:      params.UserId,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	if params.Type == Debit && u.Balance < u.BalanceLimit*-1 {
		return nil, ErrInsufficientBalance
	}

	query := postgres.CreateTransactionParams{
		UserID:      params.UserId,
		Amount:      params.Value,
		Description: params.Description,
		Ttype:       string(params.Type),
	}

	err = qtx.CreateTransaction(ctx, query)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return &TransactionCreated{
		Limit:   u.BalanceLimit,
		Balance: u.Balance,
	}, nil
}
