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
func GetDiskInfo() ([]Storage,error) {
	type storageInfo struct {
		Name       string
		Size       uint64
		FreeSpace  uint64
		FileSystem string
	}

	var storageinfo []storageInfo
	var loaclStorages []Storage
	err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
	if err != nil {
		return nil,err
	}

	for _, storage := range storageinfo {
		if storage.Size <=0 {
			continue
		}
		info := Storage{
			Name:       storage.Name,
			FileSystem: storage.FileSystem,
			Total:      storage.Size,
			Free:       storage.FreeSpace,
		}
		loaclStorages = append(loaclStorages, info)
	}
	return loaclStorages, nil

}
