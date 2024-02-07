package service

type CreateTransactionParams struct {
	UserId      float64
	Value       float64
	Type        TransactionType
	Description string
}

type TransactionCreated struct {
	Limit   float64 `json:"limite"`
	Balance float64 `json:"balance"`
}

type TransactionService interface {
	Create(params CreateTransactionParams) (TransactionCreated, error)
}
