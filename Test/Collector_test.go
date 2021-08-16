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
