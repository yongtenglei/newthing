package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/yongtenglei/newThing/registration"
	"github.com/yongtenglei/newThing/settings"
	"github.com/yongtenglei/newThing/web/logic"
	"github.com/yongtenglei/newThing/web/routers"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	register registration.ConsulRegister
	id       string
)

func init() {
	configPath := flag.String("f", "setting.yaml", "config file path")
	settings.ParseConfig(*configPath)

	register = registration.NewConsulRegister(
		settings.UserServiceConf.ConsulConf.Host,
		settings.UserServiceConf.ConsulConf.Port)

	id = uuid.New().String()

	err := register.RegisterCheckByHTTP(
		settings.UserServiceConf.UserWebClientConf.Name,
		id,
		settings.UserServiceConf.UserWebClientConf.Host,
		settings.UserServiceConf.UserWebClientConf.Port,
		settings.UserServiceConf.UserWebClientConf.Tags,
	)

	if err != nil {
		panic(err)
	}

	logic.InitUserServiceClient()
	logic.InitRefreshTokenServiceClient()
}

func main() {

	router := routers.SetUp()

	addr := fmt.Sprintf("%s:%d",
		settings.UserServiceConf.UserWebClientConf.Host,
		settings.UserServiceConf.UserWebClientConf.Port)

	fmt.Println("User Web Client running on ", addr)
	if err := router.Run(addr); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	<-sig
	// Deregister
	//internal.ConsulDeRegister(setting.ProductServiceConf.ProductWebClientConfig.ID)
	err := register.DeRegister(id)
	if err != nil {
		zap.S().Errorw(fmt.Sprintf("Register DeRegister Consul failed on %s", addr), "err", err.Error())
		fmt.Printf("Deregister failed, %s\n", err.Error())
	} else {
		fmt.Println("Deregister ok")
	}

	zap.S().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server Shutdown", zap.Error(err))
	}

	zap.S().Info("Server exiting")

}
