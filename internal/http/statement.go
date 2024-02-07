package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
)

type StatementsRestHandler struct {
	srv service.StatementService
}

func (r *StatementsRestHandler) GetStatements(c *fiber.Ctx) error {
	return c.SendString("Endpoint extrato")
}
