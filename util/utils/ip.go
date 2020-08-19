package utils

import "net"

func GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	Must(err)
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	panic("未找到")
}
