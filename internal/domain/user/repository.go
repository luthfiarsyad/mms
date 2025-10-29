package user

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

// TxRepository is optional if you use transactions
type TxRepository interface {
	Repository
}
