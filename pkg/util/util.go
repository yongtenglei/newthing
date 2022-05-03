package util

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"strings"
	"time"
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

func RandomString(size int) string {
	rand.Seed(time.Now().Unix())
	strBuilder := strings.Builder{}
	options := []byte("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < size; i++ {
		n := rand.Int31n(int32(len(options)))

		//error is always nil
		_ = strBuilder.WriteByte(options[n])
	}

	return strBuilder.String()

}

func RandomMobile() string {
	rand.Seed(time.Now().Unix())
	mobileBuilder := strings.Builder{}
	mobileBuilder.WriteByte('1')
	options := []byte("0123456789")
	for i := 0; i < 10; i++ {
		n := rand.Int31n(int32(len(options)))

		//error is always nil
		_ = mobileBuilder.WriteByte(options[n])
	}

	return mobileBuilder.String()
}

func main() {
	//fmt.Println(RandomPort("localhost"))
	fmt.Println(RandomString(32))
}
