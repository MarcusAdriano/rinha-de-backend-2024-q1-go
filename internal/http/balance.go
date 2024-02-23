package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository"
)

type BalanceRestHandler struct {
	RestHandler
	conn *repository.SqlcDatabaseConnection
}

func NewBalanceRestHandler(conn *repository.SqlcDatabaseConnection) *BalanceRestHandler {
	return &BalanceRestHandler{conn: conn}
}

func (r *BalanceRestHandler) Register(app *fiber.App) {
	app.Get("/clientes/saldos", r.GetBalances)
}

type GetBalanceResponse struct {
	UserID          int32  `json:"user_id"`
	Name            string `json:"name"`
	Limit           int64  `json:"limit"`
	Balance         int64  `json:"balance"`
	SumTransactions int64  `json:"sum_transactions"`
	Transactions    int64  `json:"transactions"`
}

// GetBalances godoc
//	@Summary		Obtem todos os saldos e a soma de todas as transacoes.
//	@Description	Saldo e somatoria das transacoes.
//	@Tags			clientes
//	@Produce		json
//	@Success		200	{array}	GetBalanceResponse
//	@Router			/clientes/saldos [get]``
func (r *BalanceRestHandler) GetBalances(c *fiber.Ctx) error {
	ctx := c.Context()

	balances, err := r.conn.GetQueries().GetAllBalance(ctx)
	if err != nil {
		return globalErrorHandler(err, c)
	}

	balance := make([]GetBalanceResponse, 0, len(balances))
	for _, b := range balances {
		balance = append(balance, GetBalanceResponse{
			UserID:          b.UserID,
			Name:            b.UserName,
			Limit:           b.Limit,
			Balance:         b.Balance,
			SumTransactions: b.SumTransactions,
			Transactions:    b.Transactions,
		})
	}

	return c.JSON(balance)
}
