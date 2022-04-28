package main

import (
	"flag"
	"fmt"
	"github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/service/biz"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

func init() {
	configPath := flag.String("f", "setting.yaml", "config file path")
	settings.ParseConfig(*configPath)
}
func main() {

	mysql.Init()

	addr := fmt.Sprintf("%s:%d",
		settings.UserServiceConf.UserWebServerConf.Host,
		settings.UserServiceConf.UserWebServerConf.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Errorw("Create listener failed", "err", err.Error())
		panic(err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &biz.UserServiceServer{})

	fmt.Println("User Server running on ", addr)

	if err := server.Serve(listener); err != nil {
		zap.S().Errorw(fmt.Sprintf("Server Serve failed in %s", addr), "err", err.Error())
		panic(err)
	}
}
