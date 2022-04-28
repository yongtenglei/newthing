package biz

import (
	"context"
	"errors"
	"github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/model"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/proto/pb"
	"go.uber.org/zap"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (us UserServiceServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterRes, error) {
	var user model.User
	r := mysql.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).First(&user)
	if r.RowsAffected > 0 {
		return nil, errors.New(e.UserAlreadyExists)
	}

	user = model.User{
		Mobile:   req.Mobile,
		Password: req.Password,
		Name:     req.Name,
		Gender:   int(req.Gender),
		Mail:     req.Mail,
	}

	err := mysql.DB.Model(&model.User{}).Save(&user).Error
	if err != nil {
		zap.S().Errorw("Register save failed", "err", err.Error())
		return nil, errors.New(e.InternalBusy)
	}

	var res pb.RegisterRes
	res.Name = req.Name
	res.Gender = req.Gender
	res.Mail = req.Mail
	res.Password = req.Password
	res.Mobile = req.Mobile

	return &res, nil
}

func (us UserServiceServer) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {
	var user model.User
	r := mysql.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).First(&user)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.UserDoesNotFound)
	}

	if user.Password != req.Password {
		return nil, errors.New(e.PasswordErr)
	}

	var res pb.LoginRes
	res.Ok = 1

	return &res, nil
}

func (us UserServiceServer) Info(ctx context.Context, req *pb.InfoReq) (*pb.InfoRes, error) {
	var user model.User
	r := mysql.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).First(&user)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.UserDoesNotFound)
	}

	var res pb.InfoRes
	res.Name = user.Name
	res.Gender = int32(user.Gender)
	res.Mail = user.Mail
	res.Mobile = user.Mobile

	return &res, nil
}

func (us UserServiceServer) Delete(ctx context.Context, req *pb.DeleteReq) (*pb.DeleteRes, error) {
	var user model.User
	r := mysql.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).First(&user)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.UserDoesNotFound)
	}

	if err := mysql.DB.Delete(&user).Error; err != nil {
		zap.S().Errorw("Delete  failed", "err", err.Error())
		return nil, errors.New(e.InternalBusy)
	}

	var res pb.DeleteRes
	res.Ok = 1

	return &res, nil
}

func (us UserServiceServer) Update(ctx context.Context, req *pb.UpdateReq) (*pb.UpdateRes, error) {
	var user model.User
	r := mysql.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).First(&user)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.UserDoesNotFound)
	}

	m := make(map[string]interface{})
	if req.Mobile != "" {
		m["mobile"] = req.Mobile
	}

	if req.Gender == 0 || user.Gender == 1 {
		m["gender"] = req.Gender
	}

	if req.Mail != "" {
		m["mail"] = req.Mail
	}

	if req.Name != "" {
		m["name"] = req.Name
	}

	if err := mysql.DB.Model(&user).Updates(m).Error; err != nil {
		zap.S().Errorw("Updates  failed", "err", err.Error())
		return nil, errors.New(e.InternalBusy)
	}

	var res pb.UpdateRes
	res.Ok = 1

	return &res, nil
}

func (us UserServiceServer) RePassword(ctx context.Context, req *pb.RePasswordReq) (*pb.RePasswordRes, error) {
	var user model.User
	r := mysql.DB.Model(&model.User{}).Where("mobile=?", req.Mobile).First(&user)
	if r.RowsAffected < 1 {
		return nil, errors.New(e.UserDoesNotFound)
	}

	m := make(map[string]interface{})
	if req.Password != "" {
		m["password"] = req.Password
	}

	if err := mysql.DB.Model(&user).Updates(m).Error; err != nil {
		zap.S().Errorw("RePassword Updates  failed", "err", err.Error())
		return nil, errors.New(e.InternalBusy)
	}

	var res pb.RePasswordRes
	res.Ok = 1

	return &res, nil
}
