package collector

import "testing"

func TestGetCPUInfo(t *testing.T) {
	if cpuInfo, err := GetCPUInfo(); err != nil {
		t.Error(err)
	} else {
		t.Log(cpuInfo)
	}
}
