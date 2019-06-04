//+build windows

package collector

import (
	"fmt"
	"github.com/StackExchange/wmi"
	log "github.com/cihub/seelog"
	"golang.org/x/sys/windows"
	"syscall"
	"time"
	"unsafe"
)

var modpsapi = windows.NewLazySystemDLL("psapi.dll")
var kernel = windows.NewLazyDLL("Kernel32.dll")
var procGetProcessMemoryInfo = modpsapi.NewProc("GetProcessMemoryInfo")

type cpuInfo struct {
	Name                      string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	MaxClockSpeed             uint32
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
			Thread: int(v.NumberOfLogicalProcessors),
			Name:   v.Name,
			Core:   int(v.NumberOfCores),
			Freq:   float64(v.MaxClockSpeed),
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
func GetUptime() (int, error) {
	type query struct {
		LastBootUpTime time.Time
	}
	var s []query
	err := wmi.Query("Select LastBootUpTime from Win32_OperatingSystem", &s)
	if err != nil {
		return 0, err
	}
	return int(time.Now().Sub(s[0].LastBootUpTime).Seconds()), nil
}
func GetKernel() (string, error) {
	return "Windows NT", nil
}

func GetSession() (int, error) {
	type query struct {
		LogonId uint32
	}
	var s []query
	err := wmi.Query("Select LogonId from Win32_LogonSession", &s)
	if err != nil {
		return 0, err
	}
	return len(s), nil
}

func GetProcess() {
	/*
		 Get-WmiObject Win32_PerfFormattedData_PerfProc_Process `
		>>     | Where-Object { $_.name -inotmatch '_total|idle' } `
		>>     | ForEach-Object {
		>>         "Process={0,-25} CPU_Usage={1,-12} Memory_Usage_(MB)={2,-16}" -f `
		>>             $_.Name,$_.PercentProcessorTime,([math]::Round($_.WorkingSetPrivate/1Mb,2))
		>>     }
	*/
	h,_:=windows.CreateToolhelp32Snapshot(0x00000002,0)
	if h < 0 {
		fmt.Println(syscall.GetLastError())
		return
	}
	var ret windows.ProcessEntry32
	ret.Size = uint32(unsafe.Sizeof(ret))
	var token windows.Token
	defer token.Close()
	for windows.Process32Next(h,&ret)==nil{
		h2,e:=windows.OpenProcess(0x1000,false,ret.ProcessID)
		if e!=nil{
			continue
		}

		_=windows.OpenProcessToken(h2,syscall.TOKEN_QUERY,&token)
		//if e == windows.ERROR_INVALID_HANDLE  || e==windows.ERROR_ACCESS_DENIED {
		//	continue
		//}
		tu,e:=token.GetTokenUser()
		if e!=nil{
			continue
		}
		user, _, _, err := tu.User.Sid.LookupAccount("")
		fmt.Printf("%+v\n",windows.UTF16ToString(ret.ExeFile[:]))
		fmt.Println()
		fmt.Println(user,  err)
	}


}
