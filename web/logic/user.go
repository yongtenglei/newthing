package logic

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/pkg/util"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

var UserServiceClient pb.UserServiceClient

func InitUserServiceClient() {
	//addr := fmt.Sprintf("%s:%d",
	//	settings.UserServiceConf.UserWebServerConf.Host,
	//	settings.UserServiceConf.UserWebServerConf.Port)
	//
	//conn, err := grpc.Dial(addr, grpc.WithInsecure())
	//if err != nil {
	//	zap.S().Errorw("grpc Dial failed", "err", err.Error())
	//}
	//UserServiceClient = pb.NewUserServiceClient(conn)

	addr := fmt.Sprintf("consul://%s:%d/%s?wait=14s",
		settings.UserServiceConf.ConsulConf.Host,
		settings.UserServiceConf.ConsulConf.Port,
		settings.UserServiceConf.UserWebServerConf.Name)
	//"consul://127.0.0.1:8500/user_web_server?wait=14s"
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		zap.S().Errorw("grpc Dial failed", "err", err.Error())

		panic(err)
	}

	UserServiceClient = pb.NewUserServiceClient(conn)

}

type RegisterReq struct {
	Mobile   string `json:"mobile" binding:"required,max=20"`
	Password string `json:"password" binding:"required,min=6,max=16"`
	Name     string `json:"name" binding:"required,min=3,max=20"`
	Gender   int    `json:"gender" binding:"oneof=0 1"`
	Mail     string `json:"mail,omitempty" binding:"min=7,max=36"`
}

type RegisterRes struct {
	Mobile   string `json:"mobile,required"`
	Password string `json:"password,required"`
	Name     string `json:"name,required"`
	Gender   int    `json:"gender,required"`
	Mail     string `json:"mail,omitempty"`
}

func RegisterHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorw("RegisterHandler ShouldBindJSON failed", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	r, err := UserServiceClient.Register(context.Background(), &pb.RegisterReq{
		Mobile:   req.Mobile,
		Password: req.Password,
		Name:     req.Name,
		Gender:   int32(req.Gender),
		Mail:     req.Mail,
	})

	if err != nil {
		zap.S().Errorw("UserServiceClient Register failed", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res RegisterRes
	res.Name = r.Name
	res.Gender = int(r.Gender)
	res.Mail = r.Mail
	res.Mobile = r.Mobile

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
}

type LoginReq struct {
	Mobile   string `json:"mobile,required" binding:"required,max=20"`
	Password string `json:"password,required" binding:"required,min=6,max=16"`
}

type LoginRes struct {
	Ok           int
	Mobile       string
	Token        string
	Issuer       string
	IssueAt      int64
	ExpiredAt    int64
	TokenSession RefreshTokenInfo
}

type RefreshTokenInfo struct {
	Uuid         string
	Mobile       string
	RefreshToken string
	Issuer       string
	UserAgent    string
	ClientIP     string
	IssuedAt     int64
	ExpiredAt    int64
}

func LoginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorw("LoginHandler ShouldBindJSON failed", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	ctx := util.NewStrContext(context.Background(), "UserAgent", c.Request.UserAgent())
	ctx = util.NewStrContext(ctx, "ClientIP", c.ClientIP())

	zap.S().Info(c.Request.UserAgent())

	r, err := UserServiceClient.Login(ctx, &pb.LoginReq{
		Mobile:   req.Mobile,
		Password: req.Password,
	})

	if err != nil {
		zap.S().Errorw("UserServiceClient Login failed", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res LoginRes
	res = LoginRes{
		Ok:        int(r.Ok),
		Mobile:    r.Mobile,
		Token:     r.Token,
		Issuer:    r.Issuer,
		IssueAt:   r.IssueAt,
		ExpiredAt: r.ExpiredAt,
		TokenSession: RefreshTokenInfo{
			Uuid:         r.TokenSession.Uuid,
			Mobile:       r.TokenSession.Mobile,
			RefreshToken: r.TokenSession.RefreshToken,
			Issuer:       r.TokenSession.Issuer,
			UserAgent:    r.TokenSession.UserAgent,
			ClientIP:     r.TokenSession.ClientIP,
			IssuedAt:     r.TokenSession.IssuedAt,
			ExpiredAt:    r.TokenSession.ExpiredAt,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
}

type InfoReq struct {
	Mobile string `uri:"mobile" json:"mobile,required" binding:"required,max=20"`
}

type InfoRes struct {
	Mobile string `json:"mobile,required"`
	Name   string `json:"name,required"`
	Gender int    `json:"gender,required"`
	Mail   string `json:"mail,omitempty"`
}

func InfoHandler(c *gin.Context) {
	var req InfoReq
	if err := c.ShouldBindUri(&req); err != nil {
		zap.S().Errorw("InfoHandler ShouldBindUri failed", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	r, err := UserServiceClient.Info(context.Background(), &pb.InfoReq{
		Mobile: req.Mobile,
	})

	if err != nil {
		zap.S().Errorw("UserServiceClient Info failed", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res InfoRes
	res.Name = r.Name
	res.Gender = int(r.Gender)
	res.Mail = r.Mail
	res.Mobile = r.Mobile

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
}

type DeleteReq struct {
	Mobile string `uri:"mobile" json:"mobile,required" binding:"required,max=20"`
}

type DeleteRes struct {
	Ok int `json:"ok"`
}

func DeleteHandler(c *gin.Context) {
	var req DeleteReq
	if err := c.ShouldBindUri(&req); err != nil {
		zap.S().Errorw("DeleteHandler ShouldBindUri failed", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	r, err := UserServiceClient.Delete(context.Background(), &pb.DeleteReq{
		Mobile: req.Mobile,
	})

	if err != nil {
		zap.S().Errorw("UserServiceClient Delete failed", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res DeleteRes
	res.Ok = int(r.Ok)

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
}

type UpdateReq struct {
	Mobile string `json:"mobile,required" binding:"required,max=20"`
	Name   string `json:"name,required" binding:"required,min=3,max=20"`
	Gender int    `json:"gender,required" binding:"oneof=0 1"`
	Mail   string `json:"mail,omitempty" binding:"min=7,max=36"`
}

type UpdateRes struct {
	Ok int `json:"ok"`
}

func UpdateHandler(c *gin.Context) {
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorw("UpdateHandler ShouldBindJSON failed", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	r, err := UserServiceClient.Update(context.Background(), &pb.UpdateReq{
		Mobile: req.Mobile,
		Name:   req.Name,
		Gender: int32(req.Gender),
		Mail:   req.Mail,
	})

	if err != nil {
		zap.S().Errorw("UserServiceClient Update failed", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res UpdateRes
	res.Ok = int(r.Ok)

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
}

type RePasswordReq struct {
	Mobile   string `json:"mobile,required" binding:"required,max=20"`
	Password string `json:"password,required" binding:"required,min=6,max=16"`
}

type RePasswordRes struct {
	Ok int `json:"ok"`
}

func RePasswordHandler(c *gin.Context) {
	var req RePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorw("RePasswordHandler ShouldBindJSON failed", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	r, err := UserServiceClient.RePassword(context.Background(), &pb.RePasswordReq{
		Mobile:   req.Mobile,
		Password: req.Password,
	})

	if err != nil {
		zap.S().Errorw("UserServiceClient RePassword failed", "err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res RePasswordRes
	res.Ok = int(r.Ok)

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
}
