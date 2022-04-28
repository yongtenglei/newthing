package main

import (
	"flag"
	"fmt"
	"github.com/yongtenglei/newThing/settings"
	"github.com/yongtenglei/newThing/web/logic"
	"github.com/yongtenglei/newThing/web/routers"
)

func init() {
	configPath := flag.String("f", "setting.yaml", "config file path")
	settings.ParseConfig(*configPath)

	logic.InitClient()
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

}
