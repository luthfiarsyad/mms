package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/luthfiarsyad/mms/internal/domain/user"
	"github.com/luthfiarsyad/mms/internal/infrastructure/logger"
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
	logger.L.Debug().Msg("AuthUsecase: initialized")
	return &AuthUsecase{userService: us, paseto: p}
}

func (a *AuthUsecase) Register(ctx context.Context, u *user.User, hashedPassword string) error {
	logger.L.Info().
		Str("email", u.Email).
		Msg("AuthUsecase.Register: start user registration")

	u.Password = hashedPassword
	u.CreatedAt = time.Now()

	err := a.userService.Register(ctx, u)
	if err != nil {
		logger.L.Error().
			Err(err).
			Str("email", u.Email).
			Msg("AuthUsecase.Register: failed to register user")
		return err
	}

	logger.L.Info().
		Str("email", u.Email).
		Int64("user_id", u.ID).
		Msg("AuthUsecase.Register: registration successful")

	return nil
}

func (a *AuthUsecase) Login(ctx context.Context, email, password string,
	passwordCheck func(hashed, plain string) error) (string, error) {

	logger.L.Info().
		Str("email", email).
		Msg("AuthUsecase.Login: login attempt")

	u, err := a.userService.Authenticate(ctx, email, password)
	if err != nil {
		logger.L.Warn().
			Err(err).
			Str("email", email).
			Msg("AuthUsecase.Login: authentication failed")
		return "", err
	}

	if err := passwordCheck(u.Password, password); err != nil {
		logger.L.Warn().
			Str("email", email).
			Msg("AuthUsecase.Login: invalid password")
		return "", errors.New("invalid credentials")
	}

	token, err := a.paseto.CreateToken(u.ID, 24*time.Hour)
	if err != nil {
		logger.L.Error().
			Err(err).
			Int64("user_id", u.ID).
			Str("email", u.Email).
			Msg("AuthUsecase.Login: failed to create token")
		return "", err
	}

	logger.L.Info().
		Int64("user_id", u.ID).
		Str("email", u.Email).
		Msg("AuthUsecase.Login: login success")

	return token, nil
}
