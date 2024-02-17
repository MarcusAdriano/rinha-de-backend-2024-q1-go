package service

import "errors"

type TransactionType string

const (
	Debit      TransactionType = "d"
	Credit     TransactionType = "c"
	DateFormat                 = "2006-01-02T15:04:05-0700"
)

var (
	ErrCustomerNotFound    error = errors.New("customer not found")
	ErrInsufficientBalance error = errors.New("insufficient balance")
)