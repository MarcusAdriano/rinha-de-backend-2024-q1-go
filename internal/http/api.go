package http

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
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

	// swagger
	r.app.Get("/swagger/*", swagger.HandlerDefault)
}

func (r *RestApp) GetApp() *fiber.App {
	return r.app
}

//	@title			Rinha Backend API - Concorrencia
//	@version		0.0.2
//	@description	Servidor Web "Rinha de Backend 2 - Concorrencia".
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	https://github.com/marcusadriano
//	@contact.email	marcusadriano.pereira@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host	localhost:8080
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
