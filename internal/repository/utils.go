package repository

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func IsErrNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
