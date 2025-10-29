package user

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidCreds = errors.New("invalid credentials")
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service { return &Service{repo: r} }
func (s *Service) Register(ctx context.Context, u *User) error {
	// repository should check unique email; we keep simple here
	return s.repo.Create(ctx, u)
}
func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}
	return u, nil
}
