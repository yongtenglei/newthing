package jwtx

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"time"
)

type UserClaims struct {
	Mobile string `json:"mobile"`
	jwt.StandardClaims
}

func CreateUserClaims(mobile string) (string, error) {
	// Create the Claims
	claims := UserClaims{
		mobile,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + settings.UserServiceConf.JWTConf.ExpireTime,
			Issuer:    settings.UserServiceConf.JWTConf.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(settings.UserServiceConf.JWTConf.SignKey))
	if err != nil {
		zap.S().Errorw("CreateUserClaims SignedString failed", "err", err.Error())
		return "", err
	}
	return ss, err
}

func ParseUserClaims(tokenString string) (*UserClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			zap.S().Error("ParseWithClaims but without SigningMethodHMAC method")
			return nil, errors.New(e.TokenMethodErr)
		}
		return []byte(settings.UserServiceConf.JWTConf.SignKey), nil
	})

	if err != nil {
		//handle error in more detail
		// TODO:
		// It may be necessary to customize the Valid function which throw a custom error if an error occurs.
		// Allow for more detailed error handling.

		//if verr , ok := err.(*jwt.ValidationError); ok && errors.Is(verr.Inner, e.ExpiredTokenErr);
		
		return nil, errors.New(err.Error())
		//return nil, errors.New(e.ParseJWTFailed)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(err.Error())
	//return nil, errors.New(e.InvalidTokenErr)
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
			return CreateUserClaims(claims.Mobile)
		}
	}

	return "", errors.New(e.InvalidTokenErr)

}
