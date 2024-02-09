package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type TransactionRestHandler struct {
	RestHandler
	srv service.TransactionService
}

func NewTransactionRestHandler(srv service.TransactionService) *TransactionRestHandler {
	return &TransactionRestHandler{srv: srv}
}

func (r *TransactionRestHandler) Register(app *fiber.App) {
	app.Post("/clientes/:id/transacoes", r.CreateTransaction)
}

type createTransactionRequest struct {
	Amount      int32  `json:"valor"`
	Type        string `json:"tipo"`
	Description string `json:"descricao"`
}

func (r *TransactionRestHandler) CreateTransaction(c *fiber.Ctx) error {

	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var req createTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	result, err := r.srv.Create(c.Context(), service.CreateTransactionParams{
		UserId:      int32(userId),
		Value:       req.Amount,
		Type:        service.TransactionType(req.Type),
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, service.ErrInsufficientLimit) {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}
		if errors.Is(err, service.ErrCustomerNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		log.Error().Err(err).Msg("Error creating transaction")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(result)
}
