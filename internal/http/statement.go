package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/rs/zerolog/log"
)

type StatementsRestHandler struct {
	RestHandler
	srv service.StatementService
}

func NewStatementRestHandler(srv service.StatementService) *StatementsRestHandler {
	return &StatementsRestHandler{srv: srv}
}

func (r *StatementsRestHandler) Register(app *fiber.App) {
	app.Get("/clientes/:id/extrato", r.GetStatements)
}

func (r *StatementsRestHandler) GetStatements(c *fiber.Ctx) error {

	ctx := c.Context()
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing user id")
		return c.Status(fiber.StatusUnprocessableEntity).Send(nil)
	}

	statementsParams := service.GetStatementsParams{
		UserId: userId,
	}

	statements, err := r.srv.GetStatements(ctx, statementsParams)
	if err != nil {
		return handleError(err, c)
	}

	return c.JSON(statements)
}
