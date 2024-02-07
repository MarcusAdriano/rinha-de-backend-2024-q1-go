package service

type TransactionType string

const (
	Debit  TransactionType = "d"
	Credit TransactionType = "c"
)

type Statements struct {
	Balance struct {
		Amount float64 `json:"total"`
		Date   string  `json:"data_extrato"`
		Limit  float64 `json:"limite"`
	}
	LastTransactions []struct {
		Value           float64         `json:"valor"`
		TransactionType TransactionType `json:"tipo"`
		Description     string          `json:"descricao"`
		Date            string          `json:"realizado_em"`
	} `json:"ultimas_transacoes"`
}

type GetStatementsParams struct {
	UserId float64
}

type StatementService interface {
	GetStatements(params GetStatementsParams) (Statements, error)
}
