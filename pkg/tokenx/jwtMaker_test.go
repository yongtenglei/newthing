package tokenx

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/pkg/util"
	"github.com/yongtenglei/newThing/settings"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	settings.ParseConfig("../../setting.yaml")

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	mobile := util.RandomMobile()
	duration := time.Duration(settings.UserServiceConf.TokenConf.ExpireTime) * time.Second

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(mobile, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.ParseToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.NotEmpty(t, payload.Mobile)
	require.Equal(t, mobile, payload.Mobile)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpireAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, payload.ExpireAt, duration)
}

func TestJWTMakerWithExpiredToken(t *testing.T) {
	settings.ParseConfig("../../setting.yaml")

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	mobile := util.RandomMobile()
	duration := -time.Minute

	token, payload, err := maker.CreateToken(mobile, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.ParseToken(token)
	require.Error(t, err)
	require.EqualError(t, err, e.ExpiredTokenErr)
	require.Nil(t, payload)

}

func TestJWTMakerWithNoneAlgorithm(t *testing.T) {
	settings.ParseConfig("../../setting.yaml")

	mobile := util.RandomMobile()
	duration := time.Minute

	payload, err := NewPayload(mobile, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	JWTToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	require.NotNil(t, JWTToken)

	token, err := JWTToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.ParseToken(token)
	require.Error(t, err)
	require.EqualError(t, err, e.TokenMethodErr)
	require.Nil(t, payload)

}
