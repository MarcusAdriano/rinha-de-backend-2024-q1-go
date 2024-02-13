package http

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
)

var (
	Validator = validator.New()
)

type TransactionRestHandler struct {
	RestHandler
	srv service.TransactionService
}

func NewTransactionRestHandler(srv service.TransactionService) *TransactionRestHandler {
	return &TransactionRestHandler{srv: srv}
}

func (r *TransactionRestHandler) Register(app *fiber.App) {
	app.Post("/clientes/:id/transacoes", createTransactionReqValidator, r.CreateTransaction)
}

type createTransactionRequest struct {
	Amount      int64  `json:"valor" validate:"required"`
	Type        string `json:"tipo" validate:"required,oneof=c d"`
	Description string `json:"descricao" validate:"required,min=1,max=10"`
}

func createTransactionReqValidator(c *fiber.Ctx) error {
	var body = new(createTransactionRequest)
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Send(nil)
	}

	err = Validator.Struct(body)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Send(nil)
	}
	return c.Next()
}

func (r *TransactionRestHandler) CreateTransaction(c *fiber.Ctx) error {

	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Send(nil)
	}

	var req createTransactionRequest
	c.BodyParser(&req)

	result, err := r.srv.Create(c.Context(), service.CreateTransactionParams{
		UserId:      int32(userId),
		Value:       req.Amount,
		Type:        service.TransactionType(req.Type),
		Description: req.Description,
	})
	if err != nil {
		return handleError(err, c)
	}

	return c.JSON(result)
}
