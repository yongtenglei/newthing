package logic

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yongtenglei/newThing/pkg/e"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

var RefreshTokenServiceClient pb.TokenSessionServiceClient

func InitRefreshTokenServiceClient() {
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

	RefreshTokenServiceClient = pb.NewTokenSessionServiceClient(conn)

}

type RefreshTokenReq struct {
	Uuid         string `json:"uuid" binding:"required"`
	Mobile       string `json:"mobile" binding:"required,max=20"`
	RefreshToken string `json:"token" binding:"required"`
}

type RefreshTokenRes struct {
	Mobile    string
	Token     string
	Issuer    string
	UserAgent string
	ClientIP  string
	IssueAt   int64
	ExpiredAt int64
}

func RefreshTokenHandler(c *gin.Context) {
	var req RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	ctx := context.WithValue(context.Background(), "UserAgent", c.Request.UserAgent())
	ctx = context.WithValue(ctx, "ClientIP", c.ClientIP())
	r, err := RefreshTokenServiceClient.RefreshToken(ctx, &pb.RefreshReq{
		Uuid:         req.Uuid,
		Mobile:       req.Mobile,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.Failed,
			"data": err.Error(),
		})
		return
	}

	var res RefreshTokenRes
	res = RefreshTokenRes{
		Mobile:    r.Mobile,
		Token:     r.Token,
		Issuer:    r.Issuer,
		UserAgent: r.UserAgent,
		ClientIP:  r.ClientIP,
		IssueAt:   r.IssuedAt,
		ExpiredAt: r.ExpiredAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"data": res,
	})
	return

}
