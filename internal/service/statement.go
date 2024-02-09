package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	ErrCustomerNotFound   error = errors.New("customer not found")
	ErrInsufficientLimit  error = errors.New("insufficient limit")
	ErrInvalidTransaction error = errors.New("invalid transaction")
)

type Statements struct {
	Balance          Balance            `json:"saldo"`
	LastTransactions []LastTransactions `json:"ultimas_transacoes"`
}

type Balance struct {
	Amount int32  `json:"total"`
	Date   string `json:"data_extrato"`
	Limit  int32  `json:"limite"`
}

type LastTransactions struct {
	Value           int32           `json:"valor"`
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
	dbpool *pgxpool.Pool
}

func (s *statementService) GetStatements(ctx context.Context, params GetStatementsParams) (*Statements, error) {

	tx, err := s.dbpool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		log.Err(err).Msg("Error starting transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := postgres.New(s.dbpool)
	qtx := q.WithTx(tx)

	query := postgres.GetTransactionsByUserParams{
		UserID: int32(params.UserId),
		Limit:  10,
	}

	rows, err := qtx.GetTransactionsByUser(ctx, query)
	if err != nil {
		log.Err(err).Msg("Error getting transactions")
		return nil, err
	}

	transactions := make([]LastTransactions, 0)
	var balance Balance
	balance.Date = time.Now().Format(DateFormat)

	for _, row := range rows {
		var transaction LastTransactions

		transaction.Value = row.Amount
		transaction.Description = row.Description

		transaction.CreateAt = row.CreatedAt.Time.Format(DateFormat)
		transaction.TransactionType = TransactionType(row.Ttype)

		transactions = append(transactions, transaction)
	}

	if len(rows) == 0 {

		u, err := qtx.GetUser(ctx, int32(params.UserId))
		if err == pgx.ErrNoRows {
			return nil, ErrCustomerNotFound
		}
		if err != nil {
			log.Err(err).Msg("Error getting user")
			return nil, err
		}

		balance.Amount = u.Balance
		balance.Limit = u.BalanceLimit
	} else {
		balance.Amount = rows[0].Balance
		balance.Limit = rows[0].BalanceLimit
	}

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

func NewStatementService(dbpool *pgxpool.Pool) StatementService {
	return &statementService{
		dbpool: dbpool,
	}
}
