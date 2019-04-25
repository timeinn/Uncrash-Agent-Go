//+build windows

package collector

import (
	"fmt"
	"github.com/StackExchange/wmi"
	log "github.com/cihub/seelog"
	"golang.org/x/sys/windows"
)

var modpsapi = windows.NewLazySystemDLL("psapi.dll")
var kernel = windows.NewLazyDLL("Kernel32.dll")

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
func GetDiskInfo() ([]Storage, error) {
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
		return nil, err
	}

	for _, storage := range storageinfo {
		if storage.Size <= 0 {
			continue
		}
		info := Storage{
			Name:       storage.Name,
			FileSystem: storage.FileSystem,
			Total:      int(storage.Size),
			Free:       int(storage.FreeSpace),
		}
		fmt.Println(storage)
		loaclStorages = append(loaclStorages, info)
	}
	return loaclStorages, nil

}

func GetMemoryInfo() (Memory, error) {
	type query struct {
		FreePhysicalMemory      uint64
		FreeSpaceInPagingFiles  uint64
		TotalVisibleMemorySize  uint64
		SizeStoredInPagingFiles uint64
	}
	var s []query
	err := wmi.Query("Select FreePhysicalMemory,FreeSpaceInPagingFiles,TotalVisibleMemorySize,SizeStoredInPagingFiles from Win32_OperatingSystem", &s)
	if err != nil {
		return Memory{}, err
	}
	if len(s) != 1 {
		return Memory{}, fmt.Errorf("query fail")
	}
	memory := Memory{}
	memory.Physical.Free = int(s[0].FreePhysicalMemory)
	memory.Physical.Total = int(s[0].TotalVisibleMemorySize)
	memory.Swap.Free = int(s[0].FreeSpaceInPagingFiles)
	memory.Swap.Total = int(s[0].SizeStoredInPagingFiles)
	return memory, nil
}
