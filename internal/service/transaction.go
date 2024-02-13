package service

import (
	"context"

	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository"
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
	dbconn *repository.SqlcDatabaseConnection
}

func NewTransactionService(conn *repository.SqlcDatabaseConnection) TransactionService {
	return &transactionService{
		dbconn: conn,
	}
}

func (s *transactionService) Create(ctx context.Context, params CreateTransactionParams) (*TransactionCreated, error) {

	tx, err := s.dbconn.GetConn().Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := s.dbconn.New()
	qtx := queries.WithTx(tx)

	var operationValue = params.Value

	if params.Type == Debit {
		operationValue *= -1
	}

	u, err := qtx.UpdateUserBalance(ctx, postgres.UpdateUserBalanceParams{
		Balance: operationValue,
		ID:      params.UserId,
	})
	if err != nil {
		if repository.IsErrNoRows(err) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	if params.Type == Debit && u.Balance+u.BalanceLimit < 0 {
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
