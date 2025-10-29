package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/luthfiarsyad/mms/internal/domain/user"
)

type AuthUsecase struct {
	userService *user.Service
	paseto      PasetoService
}

// PasetoService minimal interface for token creation/validation
type PasetoService interface {
	CreateToken(userID int64, exp time.Duration) (string, error)
	VerifyToken(token string) (int64, error)
}

func NewAuthUsecase(us *user.Service, p PasetoService) *AuthUsecase {
	return &AuthUsecase{userService: us, paseto: p}
}
func (a *AuthUsecase) Register(ctx context.Context, u *user.User,
	hashedPassword string) error {
	u.Password = hashedPassword
	u.CreatedAt = time.Now()
	return a.userService.Register(ctx, u)
}
func (a *AuthUsecase) Login(ctx context.Context, email, password string,
	passwordCheck func(hashed, plain string) error) (string, error) {
	u, err := a.userService.Authenticate(ctx, email, password)
	if err != nil {
		return "", err
	}
	// verify password
	if err := passwordCheck(u.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}
	// create token, e.g., 24 hours
	token, err := a.paseto.CreateToken(u.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return token, nil
}
