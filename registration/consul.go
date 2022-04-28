package registration

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

//type IRegistrar interface {
//	Register(name, id string, port int, tags []string) error
//	DeRegister(id string) error
//}

type ConsulRegister struct {
	Host string
	Port int
}

// NewConsulRegister 构造向host:port的consul注册器
func NewConsulRegister(host string, port int) ConsulRegister {
	cr := ConsulRegister{
		Host: host,
		Port: port,
	}
	return cr
}

// RegisterByCheckHTTP 将服务信息为参数的服务注册进cr, check method is http.
func (cr ConsulRegister) RegisterCheckByHTTP(name string, id string, host string, port int, tags []string) error {
	// consul的配置, 指向consul地址
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", cr.Host, cr.Port)

	// 创建consul地址的的客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		zap.S().Errorw("Register New Consul client failed", "err", err.Error())
		return err
	}

	// 注册信息准备
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Address: host, // 服务自己的host
		Port:    port, // 服务自己的port
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port), // 服务所提供的api接口
			Timeout:                        "5s",
			Interval:                       "3s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	// 客户端使用api进行注册
	return client.Agent().ServiceRegister(registration)
}

// RegisterCheckByGRPC 将服务信息为参数的服务注册进cr, check method is grpc.
func (cr ConsulRegister) RegisterCheckByGRPC(name string, id string, host string, port int, tags []string) error {
	// consul的配置, 指向consul地址
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", cr.Host, cr.Port)

	// 创建consul地址的的客户端
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		zap.S().Errorw("Register New Consul client failed", "err", err.Error())
		return err
	}

	// 注册信息准备
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Address: host, // 服务自己的host
		Port:    port, // 服务自己的port
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", host, port), // 服务所提供的api接口
			Timeout:                        "5s",
			Interval:                       "3s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	// 客户端使用api进行注册
	return client.Agent().ServiceRegister(registration)
}

func (cr ConsulRegister) DeRegister(id string) error {
	// 创建consul地址的的客户端
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", cr.Host, cr.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		zap.S().Errorw("DeRegister New Consul client failed", "err", err.Error())
		return err
	}

	// 客户端使用api进行反注册
	return client.Agent().ServiceDeregister(id)

}
