package biz

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/model"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/proto/pb"
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
