package tokenx

import (
	"github.com/stretchr/testify/require"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/pkg/util"
	"github.com/yongtenglei/newThing/settings"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	settings.ParseConfig("../../setting.yaml")

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	mobile := util.RandomMobile()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(mobile, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.ParseToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.NotEmpty(t, payload.Mobile)
	require.Equal(t, mobile, payload.Mobile)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, duration)
	require.WithinDuration(t, expiredAt, payload.ExpireAt, duration)
}

func TestPasetoMakerWithExpiredToken(t *testing.T) {
	settings.ParseConfig("../../setting.yaml")

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	mobile := util.RandomMobile()
	duration := -time.Minute

	token, err := maker.CreateToken(mobile, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.ParseToken(token)
	require.Error(t, err)
	require.EqualError(t, err, e.ExpiredTokenErr)
	require.Nil(t, payload)
}
