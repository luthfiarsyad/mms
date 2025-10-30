package transaction

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, t *Transaction) error
	FindByID(ctx context.Context, id int64) (*Transaction, error)
	FindByUserID(ctx context.Context, userID int64) ([]*Transaction, error)
	Update(ctx context.Context, t *Transaction) error
	Delete(ctx context.Context, id int64) error
}

type TxRepository interface {
	Repository
}