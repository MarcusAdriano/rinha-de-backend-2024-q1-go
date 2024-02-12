package http

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RestHandler interface {
	Register(app *fiber.App)
}

type RestApp struct {
	app *fiber.App
}

func (r *RestApp) RegisterHandler(handler ...RestHandler) {
	for _, h := range handler {
		h.Register(r.app)
	}
}

func (r *RestApp) Run() {

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	log.Error().Err(r.app.Listen(port)).Msg("Server startup error")
}

func NewRestApp() *RestApp {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Logger()

	app := fiber.New()
	app.Use(fiberzerolog.New(
		fiberzerolog.Config{
			Logger: &log.Logger,
		},
	))

	return &RestApp{app: app}
}

func handleError(err error, c *fiber.Ctx) error {
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
