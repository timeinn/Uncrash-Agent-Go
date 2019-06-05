package collector

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

//data_post="
// token=${auth[0]}&data=
// $(base "$version") √
// $(base "$uptime") √
// $(base "$sessions") √
// $(base "$processes")
// $(base "$processes_array")
// $(base "$file_handles")
// $(base "$file_handles_limit")
// $(base "$os_kernel") √
// $(base "$os_name") √
// $(base "$os_arch") √
// $(base "$cpu_name") √
// $(base "$cpu_cores") √
// $(base "$cpu_freq") √
// $(base "$ram_total") √
// $(base "$ram_usage") √
// $(base "$swap_total") √
// $(base "$swap_usage") √
// $(base "$disk_array")
// $(base "$disk_total")
// $(base "$disk_usage")
// $(base "$connections")
// $(base "$nic")
// $(base "$ipv4")
// $(base "$ipv6")
// $(base "$rx")
// $(base "$tx")
// $(base "$rx_gap")
// $(base "$tx_gap")
// $(base "$load")
// $(base "$load_cpu")
// $(base "$load_io")
// $(base "$ping_cn")
// $(base "$ping_hk")
// $(base "$ping_jp")
// $(base "$ping_sg")
// $(base "$ping_eu")
// $(base "$ping_us")
// $(base "$ping_as")"

type Collector interface {
	Base
	GetCPUInfo() (Cpu, error)
	GetDiskInfo() ([]Storage, error)
	GetMemoryInfo() (Memory, error)
	GetUptime() (int, error)
	GetKernel() (string, error)
	GetSession() (int, error)
}
type Base interface {
	GetHostName() (name string, err error)
	GetNetInterfaces() (interfaces []Interfaces, error error)
	GetOSArch() string
	Ext() (interface{}, error)
}
type safeRegFunc struct {
	sync.Mutex
	regFunc func() Collector
}

var _regFunc *safeRegFunc

func init() {
	_regFunc = &safeRegFunc{}
	fmt.Println(_regFunc)
}
func Register(regFunc func() Collector) {
	_regFunc.Lock()
	_regFunc.regFunc = regFunc
	_regFunc.Unlock()
}
func Default() Collector {
	fmt.Println("register collector")
	return _regFunc.regFunc()
}

type BaseCollector struct {
}

// 获取主机名
func (base BaseCollector) GetHostName() (name string, err error) {
	return os.Hostname()
}

//获取网络接口
//排除 lo 和 没ip的网卡接口
func (base BaseCollector) GetNetInterfaces() (interfaces []Interfaces, error error) {
	if inters, err := net.Interfaces(); err != nil {
		error = err
		return
	} else {
		for _, v := range inters {
			flag := v.Flags.String()
			if !strings.Contains(flag, "up") || strings.Contains(flag, "loopback") {
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
					Addr, e := addr.(*net.IPNet)
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

func (base BaseCollector) GetOSArch() string {
	return runtime.GOARCH
}

// 获取主机名
func (base BaseCollector) Ext() (interface{}, error) {
	return nil, fmt.Errorf("default ext")
}
