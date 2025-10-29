package security

import (
	"errors"
	"time"

	"github.com/luthfiarsyad/mms/config"
	"github.com/o1egl/paseto"
)

// Simple PASETO V2 service
type PasetoService struct {
	paseto *paseto.V2
	key    []byte
}

func NewPasetoService() *PasetoService {
	cfg := config.Get()
	if cfg == nil {
		panic("config is not loaded")
	}
	return &PasetoService{paseto: paseto.NewV2(), key: []byte(cfg.Paseto.SymmetricKey)}
}

type tokenPayload struct {
	ID  int64     `json:"id"`
	Exp time.Time `json:"exp"`
}

func (p *PasetoService) CreateToken(userID int64, exp time.Duration) (string,
	error) {
	pl := tokenPayload{ID: userID, Exp: time.Now().Add(exp)}
	token, err := p.paseto.Encrypt(p.key, pl, nil)

	if err != nil {
		return "", err
	}
	return token, nil
}
func (p *PasetoService) VerifyToken(token string) (int64, error) {
	var pl tokenPayload
	if err := p.paseto.Decrypt(token, p.key, &pl, nil); err != nil {
		return 0, err
	}
	if time.Now().After(pl.Exp) {
		return 0, errors.New("token expired")
	}
	return pl.ID, nil
}
