package tokenx

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"time"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < settings.UserServiceConf.TokenConf.MinSignKeySize {
		zap.S().Errorf("NewJWTMaker invalid key size: must be at least %d characters", settings.UserServiceConf.TokenConf.MinSignKeySize)
		return nil, fmt.Errorf("NewJWTMaker invalid key size: must be at least %d characters", settings.UserServiceConf.TokenConf.MinSignKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker JWTMaker) CreateToken(mobile string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(mobile, duration)
	if err != nil {
		return "", nil, err
	}

	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := JWTToken.SignedString([]byte(settings.UserServiceConf.TokenConf.SignKey))

	return token, payload, nil
}

func (maker JWTMaker) ParseToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			zap.S().Error("ParseWithClaims but without SigningMethodHMAC method")
			return nil, errors.New(e.TokenMethodErr)
		}
		return []byte(settings.UserServiceConf.TokenConf.SignKey), nil
	})

	//handle error in more detail
	if err != nil {

		verr, ok := err.(*jwt.ValidationError)
		if ok {
			switch {
			case verr.Inner.Error() == e.ExpiredTokenErr:
				return nil, errors.New(e.ExpiredTokenErr)
			case verr.Inner.Error() == e.TokenMethodErr:
				return nil, errors.New(e.TokenMethodErr)
			}

		}

		return nil, errors.New(e.InvalidTokenErr)
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, errors.New(e.InternalBusy)
	}

	return payload, nil

}
