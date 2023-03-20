package models

type TransactionType string

type Transaction struct {
	ID            int64
	AccountNumber int64
	Amount        float64
	Type          TransactionType
	Date          string
	Description   string
}
