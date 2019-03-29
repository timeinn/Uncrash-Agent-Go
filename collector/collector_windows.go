//+build windows

package collector

import (
	"github.com/StackExchange/wmi"
	log "github.com/cihub/seelog"
)

type cpuInfo struct {
	Name          string
	NumberOfCores uint32
	ThreadCount   uint32
}

func GetCPUInfo() (Cpu, error) {
	var cpuinfo []cpuInfo
	err := wmi.Query("Select * from Win32_Processor", &cpuinfo)
	if err != nil {
		_ = log.Error("windows Cannt get Cpuinfo " + err.Error())
		return Cpu{}, err
	}
	var info Cpu
	info.Num = len(cpuinfo)
	for _, v := range cpuinfo {
		info.Info = append(info.Info, CpuInfo{
			Thread: int(v.ThreadCount),
			Name:   v.Name,
			Core:   int(v.NumberOfCores),
		})

	}
	return info, nil
}
