package collector

import (
	"testing"
)

func TestGetCPUInfo(t *testing.T) {
	if cpuinfo, err := GetCPUInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(cpuinfo)
	}
}

func TestGetDiskInfo(t *testing.T) {
	if DiskInfo, err := GetDiskInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(DiskInfo)
	}
}
func TestGetMemoryInfo(t *testing.T) {
	if GetMemoryInfo, err := GetMemoryInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(GetMemoryInfo)
	}
}
func TestGetProcess(t *testing.T) {
	if GetProcess, err := GetProcess(); err != nil {
		t.Error(err)
	} else {
		t.Log(GetProcess)
	}
}

func TestGetNet(t *testing.T) {
	if GetProcess, err := GetNetInterfaces(); err != nil {
		t.Error(err)
	} else {
		t.Log(GetProcess)
	}
}

func TestGetOutboundNetInterfaces(t *testing.T) {
	if inter, err := GetOutboundNetInterfaces(); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}
func TestGetInterfacesTraffic(t *testing.T) {
	i, err := GetOutboundNetInterfaces()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(i)
	if inter, err := GetInterfacesTraffic(*i); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}

func TestGetLimit(t *testing.T) {
	if inter, err := GetLimit(); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}

func TestGetLoadAvg(t *testing.T) {
	if inter, err := GetLoadAvg(); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}
func TestGetLoad(t *testing.T) {
	if inter, err := GetLoad(); err != nil {
		t.Error(err)
	} else {
		t.Log(inter)
	}
}
