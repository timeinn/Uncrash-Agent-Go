package Test

import (
	"testing"

	"github.com/TimeInn/Uncrash-Agent-Go/collector"
)

func TestGetCPUInfo(t *testing.T) {
	if cpuinfo, err := collector.GetCPUInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(cpuinfo)
	}
}

func TestGetDiskInfo(t *testing.T) {
	if DiskInfo, err := collector.GetDiskInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(DiskInfo)
	}
}
func TestGetMemoryInfo(t *testing.T) {
	if GetMemoryInfo, err := collector.GetMemoryInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(GetMemoryInfo)
	}
}
func TestGetProcess(t *testing.T) {
	if GetProcess, err := collector.GetProcess(); err != nil {
		t.Error(err)
	} else {
		t.Log(GetProcess)
	}
}

func TestGetNet(t *testing.T) {
	if GetProcess, err := collector.GetNetInterfaces(); err != nil {
		t.Error(err)
	} else {
		t.Log(GetProcess)
	}
}

func TestGetOutboundNetInterfaces(t *testing.T) {
	if inter, err := collector.GetOutboundNetInterfaces(); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}
func TestGetInterfacesTraffic(t *testing.T) {
	i, err := collector.GetOutboundNetInterfaces()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(i)
	if inter, err := collector.GetInterfacesTraffic(*i); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}
