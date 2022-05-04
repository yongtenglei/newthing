package biz

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/pkg/tokenx"
	"github.com/yongtenglei/newThing/pkg/util"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/settings"
	"testing"
	"time"
)

var tokenSessionServiceServer TokenSessionServiceServer

func init() {
	settings.ParseConfig("setting.yaml")
	mysql.Init()
}

func TestTokenSessionServiceServer_CreateTokenSession_AND_GetTokenSession(t *testing.T) {

	// for CreateTokenSession
	jwtMaker, err := tokenx.NewJWTMaker(settings.UserServiceConf.TokenConf.SignKey)
	require.NoError(t, err)
	require.NotNil(t, jwtMaker)

	mobile := util.RandomMobile()
	require.NotEmpty(t, mobile)

	token, payload, err := jwtMaker.CreateToken(mobile, time.Duration(settings.UserServiceConf.TokenConf.ExpireTime)*time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	session, err := tokenSessionServiceServer.CreateTokenSession(context.Background(), &pb.CreateReq{
		Uuid:         payload.ID.String(),
		Mobile:       payload.Mobile,
		RefreshToken: token,
		Issuer:       payload.Issuer,
		UserAgent:    "",
		ClientIP:     "",
		IssuedAt:     payload.IssuedAt.Unix(),
		ExpiredAt:    payload.ExpireAt.Unix(),
	})
	require.NoError(t, err)
	require.Equal(t, int32(1), session.Ok)

	// for GetTokenSession
	getTokenSession, err := tokenSessionServiceServer.GetTokenSession(context.Background(), &pb.GetReq{
		Uuid:   payload.ID.String(),
		Mobile: mobile,
	})
	require.NoError(t, err)
	require.NotNil(t, getTokenSession)
	require.Equal(t, payload.ID.String(), getTokenSession.Uuid)
	require.Equal(t, payload.Mobile, getTokenSession.Mobile)
	require.Equal(t, token, getTokenSession.RefreshToken)
	require.Equal(t, payload.Issuer, getTokenSession.Issuer)
	require.Equal(t, payload.IssuedAt.Unix(), getTokenSession.IssuedAt)
	require.Equal(t, payload.ExpireAt.Unix(), getTokenSession.ExpiredAt)

}

func TestTokenSessionServiceServer_RefreshToken(t *testing.T) {
	// for CreateTokenSession
	jwtMaker, err := tokenx.NewJWTMaker(settings.UserServiceConf.TokenConf.SignKey)
	require.NoError(t, err)
	require.NotNil(t, jwtMaker)

	mobile := util.RandomMobile()
	require.NotEmpty(t, mobile)

	token, payload, err := jwtMaker.CreateToken(mobile, time.Duration(settings.UserServiceConf.TokenConf.ExpireTime)*time.Second)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	session, err := tokenSessionServiceServer.CreateTokenSession(context.Background(), &pb.CreateReq{
		Uuid:         payload.ID.String(),
		Mobile:       payload.Mobile,
		RefreshToken: token,
		Issuer:       payload.Issuer,
		UserAgent:    "",
		ClientIP:     "",
		IssuedAt:     payload.IssuedAt.Unix(),
		ExpiredAt:    payload.ExpireAt.Unix(),
	})
	require.NoError(t, err)
	require.Equal(t, int32(1), session.Ok)

	// for RefreshToken
	refreshedToken, err := tokenSessionServiceServer.RefreshToken(context.Background(), &pb.RefreshReq{
		Uuid:         payload.ID.String(),
		Mobile:       payload.Mobile,
		RefreshToken: token,
	})
	require.NoError(t, err)
	require.NotNil(t, refreshedToken)
	require.NotEqual(t, payload.ID.String(), refreshedToken.Uuid)
	require.Equal(t, payload.Mobile, refreshedToken.Mobile)
	require.NotEqual(t, token, refreshedToken.Token)
	require.Equal(t, payload.Issuer, refreshedToken.Issuer)
	require.WithinDuration(t, time.Unix(refreshedToken.IssuedAt, 0), time.Unix(refreshedToken.ExpiredAt, 0), time.Duration(settings.UserServiceConf.TokenConf.ExpireTime)*time.Second)
}
