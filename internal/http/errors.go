package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/rs/zerolog/log"
)

// globalErrorHandler is a helper function to handle business errors.
func globalErrorHandler(err error, c *fiber.Ctx) error {
	switch err {
	case service.ErrInsufficientBalance:
		return c.Status(fiber.StatusUnprocessableEntity).Send(nil)
	case service.ErrCustomerNotFound:
		return c.Status(fiber.StatusNotFound).Send(nil)
	default:
		log.Error().Err(err).Msg("Error creating transaction")
		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}
}
