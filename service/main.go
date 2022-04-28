package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"

	"github.com/yongtenglei/newThing/dao/mysql"
	"github.com/yongtenglei/newThing/proto/pb"
	"github.com/yongtenglei/newThing/registration"
	"github.com/yongtenglei/newThing/service/biz"
	"github.com/yongtenglei/newThing/settings"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	register := registration.NewConsulRegister(
		settings.UserServiceConf.ConsulConf.Host,
		settings.UserServiceConf.ConsulConf.Port)

	id := uuid.New().String()
	err = register.RegisterCheckByGRPC(
		settings.UserServiceConf.UserWebServerConf.Name,
		id,
		settings.UserServiceConf.UserWebServerConf.Host,
		settings.UserServiceConf.UserWebServerConf.Port,
		settings.UserServiceConf.UserWebServerConf.Tags)
	if err != nil {
		zap.S().Errorw("Web Server register to Consul failed", "err", err.Error())
		panic(err)
	}

	fmt.Println("User Server running on ", addr)

	go func() {
		if err := server.Serve(listener); err != nil {
			zap.S().Errorw(fmt.Sprintf("Server Serve failed in %s", addr), "err", err.Error())
			panic(err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig

	// Deregister
	err = register.DeRegister(id)
	if err != nil {
		zap.S().Errorw(fmt.Sprintf("Register DeRegister Consul failed on %s", addr), "err", err.Error())
		fmt.Printf("Deregister failed, %s\n", err.Error())
	} else {
		fmt.Println("Deregister successfully")
	}

	zap.S().Info("Shutdown ...")

	server.GracefulStop()

	zap.S().Info("Server exiting")

}
