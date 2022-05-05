package biz

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/model"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/pkg/tokenx"
	"github.com/yongtenglei/newThing/pkg/util"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TokenSessionServiceServer struct {
	pb.UnimplementedTokenSessionServiceServer
}

func (t TokenSessionServiceServer) CreateTokenSession(ctx context.Context, req *pb.CreateReq) (*pb.CreateRes, error) {
	uuidFromReq, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, errors.New("解析uuid失败")
	}

	var ts model.TokenSession
	ts = model.TokenSession{
		Model:        gorm.Model{},
		Uuid:         uuidFromReq,
		Mobile:       req.Mobile,
		RefreshToken: req.RefreshToken,
		Issuer:       req.Issuer,
		UserAgent:    req.UserAgent,
		ClientIP:     req.ClientIP,
		ExpiredAt:    req.ExpiredAt,
	}

	if err := mysql.DB.Create(&ts).Error; err != nil {
		zap.S().Errorw("CreateTokenSession mysql create failed", "err", err.Error())
		//return nil, errors.New(e.InternalBusy)
		return nil, err
	}

	var res pb.CreateRes
	res.Ok = 1

	return &res, nil
}

// GetTokenSession TODO: 查找过期的Token, 删除后在进行查找
func (t TokenSessionServiceServer) GetTokenSession(ctx context.Context, req *pb.GetReq) (*pb.GetRes, error) {
	var ts model.TokenSession
	r := mysql.DB.Model(&model.TokenSession{}).Where("uuid=? AND mobile=?", req.Uuid, req.Mobile).First(&ts)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.TokenSessionDoesNotFound)
	}

	var res pb.GetRes
	res = pb.GetRes{
		Uuid:         ts.Uuid.String(),
		Mobile:       ts.Mobile,
		RefreshToken: ts.RefreshToken,
		Issuer:       ts.Issuer,
		UserAgent:    ts.UserAgent,
		ClientIP:     ts.ClientIP,
		IssuedAt:     ts.IssuedAt.Unix(),
		ExpiredAt:    ts.ExpiredAt,
	}

	return &res, nil
}

func (t TokenSessionServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshReq) (*pb.RefreshRes, error) {
	jwtMaker, err := tokenx.NewJWTMaker(settings.UserServiceConf.TokenConf.SignKey)
	if err != nil {
		zap.S().Errorw("RefreshToken NewJWTMaker failed", "err", err.Error())
		return nil, errors.New(e.InternalBusy)
	}
	_, err = jwtMaker.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, err
		//return nil, errors.New(e.InvalidTokenErr)
	}

	var ts model.TokenSession
	r := mysql.DB.Model(&model.TokenSession{}).Where("uuid=? AND mobile=? AND refresh_token=?", req.Uuid, req.Mobile, req.RefreshToken).First(&ts)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.InvalidRefreshTokenErr)
	}

	JWTToken, payload, err := jwtMaker.CreateToken(req.Mobile, time.Duration(settings.UserServiceConf.TokenConf.ExpireTime)*time.Second)
	if err != nil {
		return nil, errors.New(e.InternalBusy)
	}

	var userAgent string
	var clientIP string
	var ok bool
	userAgent, ok = util.FromContextForStr(ctx, "UserAgent")
	if !ok {
		userAgent = ""
	}
	clientIP, ok = util.FromContextForStr(ctx, "ClientIP")
	if !ok {
		clientIP = ""
	}

	var res pb.RefreshRes
	res = pb.RefreshRes{
		Uuid:      payload.ID.String(),
		Mobile:    payload.Mobile,
		Token:     JWTToken,
		Issuer:    payload.Issuer,
		UserAgent: userAgent,
		ClientIP:  clientIP,
		IssuedAt:  payload.IssuedAt.Unix(),
		ExpiredAt: payload.ExpireAt.Unix(),
	}

	return &res, nil
}
