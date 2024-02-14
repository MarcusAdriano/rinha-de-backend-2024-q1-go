package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
)

type StatementsRestHandler struct {
	RestHandler
	srv service.StatementService
}

func NewStatementRestHandler(srv service.StatementService) *StatementsRestHandler {
	return &StatementsRestHandler{srv: srv}
}

func (r *StatementsRestHandler) Register(app *fiber.App) {
	app.Get("/clientes/:id/extrato", validateGetStatements, r.GetStatements)
}

func validateGetStatements(c *fiber.Ctx) error {
	userId := c.Params("id")
	if _, err := strconv.Atoi(userId); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).Send(nil)
	}
	return c.Next()
}

func (r *StatementsRestHandler) GetStatements(c *fiber.Ctx) error {

	ctx := c.Context()
	userId, _ := strconv.Atoi(c.Params("id"))

	statementsParams := service.GetStatementsParams{
		UserId: userId,
	}

	statements, err := r.srv.GetStatements(ctx, statementsParams)
	if err != nil {
		return globalErrorHandler(err, c)
	}

	return c.JSON(statements)
}
