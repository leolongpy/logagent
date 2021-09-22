package common

import (
	"fmt"
	"net"
	"strings"
)

// CollectEntry 日志的配置项
type CollectEntry struct {
	Path  string `json:"path"`  //日志路径
	Topic string `json:"topic"` //发往kafka的topic
}

// GetOutboundIP 获取本机ip的函数
func GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	loaclAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(loaclAddr.String())
	ip = strings.Split(loaclAddr.IP.String(), ":")[0]
	return
}
