package collector

import (
	"errors"
	"net"
	"os"
	"runtime"
)

var ErrorNotFound = errors.New("data not found")

// 获取主机名
func GetHostName() (name string, err error) {
	return os.Hostname()
}
func GetOSArch() string {
	return runtime.GOARCH
}
func GetOS() string {
	return runtime.GOOS
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

//GetOutboundIP 获取访问外网的网卡本机IP
// 建立一个UDP然后提取本机IP
func GetOutboundIP() (*net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return &localAddr.IP, nil
}

// 获取访问外网的网卡
func GetOutboundNetInterfaces() (*net.Interface, error) {
	ip, err := GetOutboundIP()
	if err != nil {
		return nil, err
	}
	nets, err := GetNetInterfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range nets {
		if ips, err := v.Addrs(); err != nil {
			continue
		} else {
			for _, i := range ips {
				if ri, ok := i.(*net.IPNet); ok && ri.IP.Equal(*ip) {
					return &v, nil
				}
			}
		}

	}
	return nil, ErrorNotFound

}
