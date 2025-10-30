package transaction

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
	ErrInvalidType         = errors.New("transaction type must be 'income' or 'expense'")
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) Create(ctx context.Context, t *Transaction) error {
	// Validate transaction
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	
	if t.Type != string(TransactionTypeIncome) && t.Type != string(TransactionTypeExpense) {
		return ErrInvalidType
	}
	
	if t.UserID <= 0 {
		return errors.New("user ID is required")
	}
	
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	
	return s.repo.Create(ctx, t)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetByUserID(ctx context.Context, userID int64) ([]*Transaction, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) Update(ctx context.Context, t *Transaction) error {
	// Validate transaction
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}
	
	if t.Type != string(TransactionTypeIncome) && t.Type != string(TransactionTypeExpense) {
		return ErrInvalidType
	}
	
	t.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, t)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}