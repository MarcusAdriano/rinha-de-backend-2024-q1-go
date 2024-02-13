package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusadriano/rinha-de-backend-2024-q1/internal/repository/postgres"
)

type SqlcDatabaseConnection struct {
	dbpool *pgxpool.Pool
}

func NewSqlcDatabaseConnection(pool *pgxpool.Pool) *SqlcDatabaseConnection {
	return &SqlcDatabaseConnection{
		dbpool: pool,
	}
}

func (s *SqlcDatabaseConnection) New() *postgres.Queries {
	return postgres.New(s.dbpool)
}

func (s *SqlcDatabaseConnection) GetConn() *pgxpool.Pool {
	return s.dbpool
}

func (s *SqlcDatabaseConnection) Close() {
	s.dbpool.Close()
}
