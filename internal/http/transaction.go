package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
)

type TransactionRestHandler struct {
	srv service.TransactionService
}

func (r *TransactionRestHandler) CreateTransaction(c *fiber.Ctx) error {
	return c.SendString("tbd")
}
