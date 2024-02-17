package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	constraintViolationErrId = "23503"
)

func IsErrNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func IsConstraintViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == constraintViolationErrId
	}
	return false
}
