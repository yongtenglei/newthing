package settings

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var NacosConf NacosConfig
var UserServiceConf UserServiceConfig
var fileName string
var ext string

func ParseConfig(filepath string) {
	ext = path.Ext(filepath)
	fileName = strings.TrimSuffix(filepath, ext)
	InitNacos(fileName, ext[1:])
	InitFromNacos()

}

func InitNacos(filename, ext string) {
	viper.SetConfigName(filename) // name of config file (without extension)
	viper.SetConfigType(ext)      // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../..")  // optionally look for config in the working directory
	viper.AddConfigPath("..")     // optionally look for config in the working directory
	viper.AddConfigPath(".")      // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			zap.S().Errorw("Init config failed", "info", "config file not found", "err", err)
			panic(err)
		} else {
			// Config file was found but another error was produced
			zap.S().Errorw("Init config failed", "info", "config file found", "err", err)
			panic(err)
		}
	}

	// Config file found and successfully parsed
	if err := viper.Unmarshal(&NacosConf); err != nil {
		zap.S().Errorw("Unmarshal config failed", "err", err)
		panic(err)
	}

	fmt.Println(NacosConf)

	//once := sync.Once{}
	//once.Do(func() {
	//	// watch viper config changes
	//	viper.OnConfigChange(func(e fsnotify.Event) {
	//		ParseConfig(filename)
	//		zap.S().Infow("viper config file changed", "info", filename)
	//		fmt.Println("Config file changed:", e.Name)
	//	})
	//	viper.WatchConfig()
	//
	//})
}

func InitFromNacos() {
	sc := []constant.ServerConfig{
		{
			IpAddr: NacosConf.IpAddr,
			Port:   uint64(NacosConf.Port),
		},
	}
	//or a more graceful way to create ServerConfig
	//_ = []constant.ServerConfig{
	//*constant.NewServerConfig(
	//"console.nacos.io",
	//80,
	//constant.WithScheme("http"),
	//constant.WithContextPath("/nacos")),
	//}

	cc := constant.ClientConfig{
		NamespaceId:         NacosConf.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              NacosConf.LogDir,
		CacheDir:            NacosConf.CacheDir,
		LogLevel:            NacosConf.LogLevel,
	}
	//or a more graceful way to create ClientConfig
	//_ = *constant.NewClientConfig(
	//constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"),
	//constant.WithTimeoutMs(5000),
	//constant.WithNotLoadCacheAtStart(true),
	//constant.WithLogDir("/tmp/nacos/log"),
	//constant.WithCacheDir("/tmp/nacos/cache"),
	//constant.WithLogLevel("debug"),
	//)

	// a more graceful way to create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		zap.S().Errorw("New nacos client failed", "err", err)
		panic(err)
	}

	//get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: NacosConf.DataId,
		Group:  NacosConf.Group,
	})
	if err != nil {
		zap.S().Errorw("Get nacos config failed", "err", err)
		panic(err)

	}

	fmt.Println("GetConfig,config :" + content)

	err = json.Unmarshal([]byte(content), &UserServiceConf)
	if err != nil {
		zap.S().Errorw("Parse nacos config to user service config failed", "err", err)
		panic(err)
	}

	fmt.Printf("parsed UserWebServerConf: %#v\n", UserServiceConf.UserWebServerConf)
	fmt.Printf("parsed UserWebClientConf: %#v\n", UserServiceConf.UserWebClientConf)
	fmt.Printf("parsed ConsulConf: %#v\n", UserServiceConf.ConsulConf)

	once := sync.Once{}
	once.Do(func() {
		//Listen config change,key=dataId+group+namespaceId.
		err = client.ListenConfig(vo.ConfigParam{
			DataId: NacosConf.DataId,
			Group:  NacosConf.Group,
			OnChange: func(namespace, group, dataId, data string) {
				InitFromNacos()
				zap.S().Infow("nacos config has changed ")
				fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			},
		})

		if err != nil {
			zap.S().Errorw("Start listen nacos failed", "err", err)
			panic(err)
		}

	})

}
