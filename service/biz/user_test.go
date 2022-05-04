package biz

import (
	"context"
	"fmt"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"github.com/yongtenglei/newThing/dao/mysql"
	_ "github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/settings"
	"testing"
)

var userServiceServer UserServiceServer

func init() {
	settings.ParseConfig("setting.yaml")
	mysql.Init()
}

func TestUserServiceServer_Register(t *testing.T) {

	for i := 0; i < 10; i++ {
		res, err := userServiceServer.Register(context.Background(), &pb.RegisterReq{
			Mobile:   fmt.Sprintf("1580000000%d", i),
			Password: fmt.Sprintf("1580000000%d", i),
			Name:     fmt.Sprintf("%d", i),
			Gender:   int32(i % 2),
			Mail:     fmt.Sprintf("1580000000%d", i),
		})

		assert.Empty(t, err)
		assert.Equal(t, fmt.Sprintf("1580000000%d", i), res.Mobile)
		assert.Equal(t, fmt.Sprintf("%d", i), res.Name)
		assert.Equal(t, int32(i%2), res.Gender)
		assert.Equal(t, fmt.Sprintf("1580000000%d", i), res.Mail)

		fmt.Println(res)
	}

}

func TestUserServiceServer_Login(t *testing.T) {
	for i := 0; i < 10; i++ {
		res, err := userServiceServer.Login(context.Background(), &pb.LoginReq{
			Mobile:   fmt.Sprintf("1580000000%d", i),
			Password: fmt.Sprintf("1580000000%d", i),
		})

		assert.Empty(t, err)
		assert.Equal(t, int32(1), res.Ok)
	}

}

func TestUserServiceServer_Info(t *testing.T) {
	for i := 0; i < 3; i++ {
		res, err := userServiceServer.Info(context.Background(), &pb.InfoReq{
			Mobile: fmt.Sprintf("1580000000%d", i),
		})

		assert.Empty(t, err)
		assert.Equal(t, fmt.Sprintf("1580000000%d", i), res.Mobile)
		assert.Equal(t, fmt.Sprintf("%d", i), res.Name)
		assert.Equal(t, int32(i%2), res.Gender)
		assert.Equal(t, fmt.Sprintf("1580000000%d", i), res.Mail)

	}

}

func TestUserServiceServer_Delete(t *testing.T) {
	for i := 0; i < 10; i++ {
		res, err := userServiceServer.Delete(context.Background(), &pb.DeleteReq{
			Mobile: fmt.Sprintf("1580000000%d", i),
		})

		assert.Empty(t, err)
		assert.Equal(t, int32(1), res.Ok)
	}

}

func TestUserServiceServer_Update(t *testing.T) {
	for i := 0; i < 10; i++ {
		res, err := userServiceServer.Update(context.Background(), &pb.UpdateReq{
			Mobile: fmt.Sprintf("1580000000%d", i),
			Name:   fmt.Sprintf("%d-%s", i, "updated"),
			Gender: int32(i % 2),
			Mail:   fmt.Sprintf("1580000000%d", i),
		})

		assert.Empty(t, err)
		assert.Equal(t, int32(1), res.Ok)

	}

}

func TestUserServiceServer_RePassword(t *testing.T) {
	for i := 0; i < 10; i++ {
		res, err := userServiceServer.RePassword(context.Background(), &pb.RePasswordReq{
			Mobile:   fmt.Sprintf("1580000000%d", i),
			Password: fmt.Sprintf("1580000000%d", i+1),
		})

		assert.Empty(t, err)
		assert.Equal(t, int32(1), res.Ok)

	}
}
