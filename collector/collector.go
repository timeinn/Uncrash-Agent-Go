package collector

import (
	"net"
	"os"
	"runtime"
)

// 获取主机名
func GetHostName() (name string, err error) {
	return os.Hostname()
}

//获取网络接口
//排除 lo 和 没ip的网卡接口
func GetNetInterfaces() (interfaces []net.Interface, error error) {
	if inters, err := net.Interfaces(); err != nil {
		error = err
		return
	} else {
		for _, v := range inters {
			if v.Flags&net.FlagUp == 0 || v.Flags&net.FlagLoopback != 0 {
				continue
			}
			if Addrs, err := v.Addrs(); err == nil {
				if len(Addrs) <= 0 {
					continue
				}
				interfaces = append(interfaces, v)
			}
		}
		return
	}
}

func GetOSArch() string {
	return runtime.GOARCH
}
func GetOS() string {
	return runtime.GOOS
}
func GetOutboundIP() (*net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return &localAddr.IP, nil
}
