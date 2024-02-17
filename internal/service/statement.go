package service

import (
	"context"
	"time"

	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository/postgres"
	"github.com/rs/zerolog/log"
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

	qtx := s.dbconn.GetQueries()

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

	return &Statements{
		Balance:          balance,
		LastTransactions: transactions,
	}, nil
}
