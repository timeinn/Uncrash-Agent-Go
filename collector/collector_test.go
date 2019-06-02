package collector

//data_post="token=${auth[0]}&data=$(base "$version") $(base "$uptime") $(base "$sessions") $(base "$processes") $(base "$processes_array") $(base "$file_handles") $(base "$file_handles_limit") $(base "$os_kernel") $(base "$os_name") $(base "$os_arch") $(base "$cpu_name") $(base "$cpu_cores") $(base "$cpu_freq") $(base "$ram_total") $(base "$ram_usage") $(base "$swap_total") $(base "$swap_usage") $(base "$disk_array") $(base "$disk_total") $(base "$disk_usage") $(base "$connections") $(base "$nic") $(base "$ipv4") $(base "$ipv6") $(base "$rx") $(base "$tx") $(base "$rx_gap") $(base "$tx_gap") $(base "$load") $(base "$load_cpu") $(base "$load_io") $(base "$ping_cn") $(base "$ping_hk") $(base "$ping_jp") $(base "$ping_sg") $(base "$ping_eu") $(base "$ping_us") $(base "$ping_as")"
import (
	"testing"
)

func TestGetCPUInfo(t *testing.T) {
	if cpuInfo, err := GetCPUInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(cpuInfo)
	}
}
func TestGetHostName(t *testing.T) {
	if name, err := GetHostName(); err != nil {
		t.Error(err)
	} else {
		t.Log(name)
	}
}
func TestGetNetInterfaces(t *testing.T) {
	if NetInterfaces, err := GetNetInterfaces(); err != nil {
		t.Error(err)
	} else {
		t.Log(NetInterfaces)
	}
}

func TestGetDiskInfo(t *testing.T) {
	if diskInfo, err := GetDiskInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(diskInfo)
	}
}
func TestGetMemoryInfo(t *testing.T) {
	if meminfo, err := GetMemoryInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(meminfo)
	}
}
func TestGetUptime(t *testing.T) {
	if uptime, err := GetUptime(); err != nil {
		t.Error(err)
	} else {
		t.Log(uptime)
	}
}
func TestGetKernel(t *testing.T) {
	if kernel, err := GetKernel(); err != nil {
		t.Error(err)
	} else {
		t.Log(kernel)
	}
}
func TestGetOSArch(t *testing.T) {
	t.Log(GetOSArch())
}
func TestGetSession(t *testing.T) {
	if sess, err := GetSession(); err != nil {
		t.Error(err)
	} else {
		t.Log(sess)
	}
}
