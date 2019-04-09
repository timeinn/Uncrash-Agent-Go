package collector

import (
	"net"
	"os"
	"strings"
)

// 获取主机名
func GetHostName() (name string, err error) {
	return os.Hostname()
}

//获取网络接口
//排除 lo 和 没ip的网卡接口
func GetNetInterfaces() (interfaces []Interfaces, error error) {
	if inters, err := net.Interfaces(); err != nil {
		error = err
		return
	} else {
		for _, v := range inters {
			flag:=v.Flags.String()
			if !strings.Contains(flag,"up") || strings.Contains(flag,"loopback"){
				continue
			}
			var inter = Interfaces{}
			inter.Name = v.Name
			if Addrs, err := v.Addrs(); err != nil {
				error = err
				return
			} else {
				if len(Addrs) <= 0 {
					continue
				}
				for _, addr := range Addrs {
					Addr,e:=addr.(*net.IPNet)
					if !e {
						continue
					}
					inter.Addrs = append(inter.Addrs, Addr.IP.String())
				}
			}
			interfaces = append(interfaces, inter)
		}
		return
	}

}
