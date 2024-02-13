package main

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/http"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/service"
	"github.com/rs/zerolog/log"
)

func createApp() *http.RestApp {
	app := http.NewRestApp()

	connectionString := os.Getenv("DATABASE_URL")
	dbconn := createDatabaseConnection(connectionString)

	srvStatement := service.NewStatementService(dbconn)
	srvTransaction := service.NewTransactionService(dbconn)

	app.RegisterHandler(
		http.NewStatementRestHandler(srvStatement),
		http.NewTransactionRestHandler(srvTransaction),
	)

	return app
}

func createDatabaseConnection(connectionString string) *repository.SqlcDatabaseConnection {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Error().Msgf("Error connecting to database: %v", err)
		panic(err)
	}

	config.ConnConfig.Config.ConnectTimeout = time.Second * 1
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Msgf("Error connecting to database: %v", err)
		panic(err)
	}
	return repository.NewSqlcDatabaseConnection(pool)
}

func main() {

	app := createApp()
	app.Run()
}
