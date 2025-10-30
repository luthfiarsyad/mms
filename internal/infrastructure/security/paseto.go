package security

import (
	"encoding/base64"
	"errors"
	"fmt"
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
	
	// Decode the base64 key
	key, err := base64.StdEncoding.DecodeString(cfg.Paseto.SymmetricKey)
	if err != nil {
		panic(fmt.Sprintf("failed to decode PASETO key: %v", err))
	}
	
	// Verify key length (PASETO V2 requires 32 bytes)
	if len(key) != 32 {
		panic(fmt.Sprintf("invalid PASETO key length: got %d bytes, expected 32 bytes", len(key)))
	}
	
	return &PasetoService{paseto: paseto.NewV2(), key: key}
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
