package service

import (
	"context"
	"errors"
	"time"

	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository/postgres"
	"github.com/rs/zerolog/log"
)

type TransactionType string

const (
	Debit      TransactionType = "d"
	Credit     TransactionType = "c"
	DateFormat                 = "2006-01-02T15:04:05-0700"
)

var (
	ErrCustomerNotFound    error = errors.New("customer not found")
	ErrInsufficientBalance error = errors.New("insufficient balance")
)

type Statements struct {
	Balance          Balance       `json:"saldo"`
	LastTransactions []Transaction `json:"ultimas_transacoes"`
}

type Balance struct {
	Amount int64  `json:"total"`
	Date   string `json:"data_extrato"`
	Limit  int64  `json:"limite"`
}

type Transaction struct {
	Value           int64           `json:"valor"`
	TransactionType TransactionType `json:"tipo"`
	Description     string          `json:"descricao"`
	CreateAt        string          `json:"realizado_em"`
}

type GetStatementsParams struct {
	UserId int
}

type StatementService interface {
	GetStatements(ctx context.Context, params GetStatementsParams) (*Statements, error)
}

type statementService struct {
	dbconn *repository.SqlcDatabaseConnection
}

func NewStatementService(conn *repository.SqlcDatabaseConnection) StatementService {
	return &statementService{
		dbconn: conn,
	}
}

func (s *statementService) GetStatements(ctx context.Context, params GetStatementsParams) (*Statements, error) {

	tx, err := s.dbconn.GetConn().Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Error starting transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	queries := s.dbconn.GetQueries()
	qtx := queries.WithTx(tx)

	user, err := qtx.GetUser(ctx, int32(params.UserId))
	if err != nil {
		if repository.IsErrNoRows(err) {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}

	query := postgres.GetTransactionsByUserParams{
		UserID: int32(params.UserId),
		Limit:  10,
	}

	rows, err := qtx.GetTransactionsByUser(ctx, query)
	if err != nil {
		log.Err(err).Msg("Error getting transactions")
		return nil, err
	}

	transactions := make([]Transaction, 0)
	var balance Balance
	balance.Date = time.Now().Format(DateFormat)

	for _, row := range rows {
		var transaction Transaction

		transaction.Value = row.Amount
		transaction.Description = row.Description

		transaction.CreateAt = row.CreatedAt.Time.Format(DateFormat)
		transaction.TransactionType = TransactionType(row.Ttype)

		transactions = append(transactions, transaction)
	}

	balance.Amount = user.Balance
	balance.Limit = user.BalanceLimit
	balance.Date = time.Now().Format(DateFormat)

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("Error commiting transaction")
		return nil, err
	}

	return &Statements{
		Balance:          balance,
		LastTransactions: transactions,
	}, nil
}
