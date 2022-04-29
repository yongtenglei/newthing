package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"time"
)

type UserClaims struct {
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	jwt.StandardClaims
}

func CreateUserClaims(name, mobile string) (string, error) {
	// Create the Claims
	claims := UserClaims{
		name,
		mobile,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    settings.UserServiceConf.JWTConf.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(settings.UserServiceConf.JWTConf.SignKey)
	if err != nil {
		zap.S().Errorw("CreateUserClaims SignedString failed", "err", err.Error())
	}
	return ss, err
}

func ParseUserClaims(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.UserServiceConf.JWTConf.SignKey), nil
	})

	if err != nil {
		return nil, errors.New(e.ParseJWTFailed)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(e.InvalidTokenErr)
}

// RefreshToken 刷新Token
func RefreshToken(token string) (string, error) {

	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.UserServiceConf.JWTConf.SignKey), nil
	})

	if err != nil {
		return "", err
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UserClaims); ok && tokenClaims.Valid {
			claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(settings.UserServiceConf.JWTConf.ExpireTime)).Unix()
			return CreateUserClaims(claims.Name, claims.Mobile)
		}
	}

	return "", errors.New(e.InvalidTokenErr)

}
