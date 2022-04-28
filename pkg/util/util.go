package util

import (
	"fmt"
	"go.uber.org/zap"
	"net"
)

func RandomPort(host string) int {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, 0))
	if err != nil {
		zap.S().Info("RandomPort ResolveTCPAddr failed")
		return -1
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		zap.S().Info("RandomPort ListenTCP failed")
		return -1
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func main() {
	fmt.Println(RandomPort("localhost"))
}
