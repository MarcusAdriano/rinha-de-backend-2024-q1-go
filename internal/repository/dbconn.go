package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository/postgres"
)

var (
	sharedQueries *postgres.Queries
)

type SqlcDatabaseConnection struct {
	dbpool *pgxpool.Pool
}

func NewSqlcDatabaseConnection(pool *pgxpool.Pool) *SqlcDatabaseConnection {
	sharedQueries = postgres.New(pool)
	return &SqlcDatabaseConnection{
		dbpool: pool,
	}
}

func (s *SqlcDatabaseConnection) GetConn() *pgxpool.Pool {
	return s.dbpool
}

func (s *SqlcDatabaseConnection) GetQueries() *postgres.Queries {
	return sharedQueries
}

func (s *SqlcDatabaseConnection) Close() {
	s.dbpool.Close()
}
