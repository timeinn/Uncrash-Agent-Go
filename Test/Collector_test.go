package Test

import (
	"github.com/TimeInn/Uncrash-Agent-Go/collector"
	"testing"
)

func TestGetCPUInfo(t *testing.T) {
	if cpuInfo, err := collector.Default().GetCPUInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(cpuInfo)
	}
}
func TestGetHostName(t *testing.T) {
	if name, err := collector.Default().GetHostName(); err != nil {
		t.Error(err)
	} else {
		t.Log(name)
	}
}
func TestGetNetInterfaces(t *testing.T) {
	if NetInterfaces, err := collector.Default().GetNetInterfaces(); err != nil {
		t.Error(err)
	} else {
		t.Log(NetInterfaces)
	}
}

func TestGetDiskInfo(t *testing.T) {
	if diskInfo, err := collector.Default().GetDiskInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(diskInfo)
	}
}
func TestGetMemoryInfo(t *testing.T) {
	if meminfo, err := collector.Default().GetMemoryInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(meminfo)
	}
}
func TestGetUptime(t *testing.T) {
	if uptime, err := collector.Default().GetUptime(); err != nil {
		t.Error(err)
	} else {
		t.Log(uptime)
	}
}
func TestGetKernel(t *testing.T) {
	if kernel, err := collector.Default().GetKernel(); err != nil {
		t.Error(err)
	} else {
		t.Log(kernel)
	}
}
func TestGetOSArch(t *testing.T) {
	t.Log(collector.Default().GetOSArch())
}
func TestGetSession(t *testing.T) {
	if sess, err := collector.Default().GetSession(); err != nil {
		t.Error(err)
	} else {
		t.Log(sess)
	}
}
