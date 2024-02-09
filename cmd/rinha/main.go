package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/http"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/rs/zerolog/log"
)

func main() {

	app := http.NewRestApp()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error().Msgf("Error connecting to database: %v", err)
		panic(err)
	}

	config.ConnConfig.Config.ConnectTimeout = time.Second * 2
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Msgf("Error connecting to database: %v", err)
		panic(err)
	}
	defer pool.Close()

	srvStatement := service.NewStatementService(pool)
	srvTransaction := service.NewTransactionService(pool)

	app.RegisterHandler(
		http.NewStatementRestHandler(srvStatement),
		http.NewTransactionRestHandler(srvTransaction),
	)

	app.Run()
}
