package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello Rinha de backend 2024!!")
	})

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal().Msgf("Error: %s", err)
	}
	log.Info().Msg("App running on port :3000")
}
