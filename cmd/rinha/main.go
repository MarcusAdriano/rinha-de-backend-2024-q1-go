package main

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

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

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, rinha de backend 2024 q1")
	})

	log.Fatal().Err(app.Listen(":3000"))
}
