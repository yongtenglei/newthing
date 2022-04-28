package settings

type NacosConfig struct {
	IpAddr      string `mapstructure:"ipAddr"`
	Port        int    `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespaceId"`
	DataId      string `mapstructure:"dataId"`
	Group       string `mapstructure:"group"`
	LogDir      string `mapstructure:"logDir"`
	CacheDir    string `mapstructure:"cacheDir"`
	LogLevel    string `mapstructure:"logLevel"`
}

type UserWebServerConfig struct {
	Host string   `json:"host,required"`
	Port int      `json:"port,required"`
	Name string   `json:"name,required"`
	Tags []string `json:"tags,required"`
}

type UserWebClientConfig struct {
	Host string   `json:"host,required"`
	Port int      `json:"port,required"`
	Name string   `json:"name,required"`
	Tags []string `json:"tags,required"`
}

type ConsulConfig struct {
	Host string `json:"host,required"`
	Port int    `json:"port,required"`
}

type UserServiceConfig struct {
	UserWebServerConf *UserWebServerConfig `json:"server"`
	UserWebClientConf *UserWebClientConfig `json:"client"`
	ConsulConf        *ConsulConfig        `json:"consul"`
	MySQLConf         *MySQLConfig         `json:"mysql"`
}

type MySQLConfig struct {
	Host     string `json:"host,required"`
	Port     int    `json:"port,required"`
	User     string `json:"user,required"`
	Password string `json:"password,required"`
	DbName   string `json:"dbName,required"`
}