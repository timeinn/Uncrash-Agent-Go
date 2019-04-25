package collector

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
