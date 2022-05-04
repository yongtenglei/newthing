package tokenx

import (
	"errors"
	"fmt"
	"github.com/o1egl/paseto"
	"github.com/yongtenglei/newThing/pkg/e"
	"go.uber.org/zap"
	"golang.org/x/crypto/chacha20"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

const symmetricKeySize = chacha20.KeySize

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != symmetricKeySize {
		zap.S().Errorw("NewPasetoMaker invalid key size: must be %d characters", symmetricKeySize)
		return nil, fmt.Errorf("NewPasetoMaker invalid key size: must be %d characters", symmetricKeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil

}

func (maker *PasetoMaker) CreateToken(mobile string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(mobile, duration)
	if err != nil {
		return "", nil, err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)

	return token, payload, err
}

func (maker *PasetoMaker) ParseToken(tokenString string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(tokenString, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, errors.New(e.InvalidTokenErr)
	}

	err = payload.Valid()
	if err != nil {
		return nil, errors.New(e.ExpiredTokenErr)
	}

	return payload, nil
}
