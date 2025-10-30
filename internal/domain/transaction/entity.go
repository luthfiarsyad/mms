package transaction

import (
	"time"
)

type Transaction struct {
	ID          int64     `db:"id" json:"id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	Amount      float64   `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	Type        string    `db:"type" json:"type"` // "income" or "expense"
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)