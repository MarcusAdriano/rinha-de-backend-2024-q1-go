// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Transaction struct {
	ID          pgtype.UUID
	UserID      int32
	Amount      int64
	Description string
	CreatedAt   pgtype.Timestamp
	Ttype       string
}

type User struct {
	ID           int32
	Name         string
	Balance      int64
	BalanceLimit int64
}
