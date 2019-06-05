//+build windows

package windows

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"github.com/TimeInn/Uncrash-Agent-Go/collector"
	log "github.com/cihub/seelog"
	"golang.org/x/sys/windows"
	"syscall"
	"time"
	"unsafe"
)

var modpsapi = windows.NewLazySystemDLL("psapi.dll")
var kernel = windows.NewLazyDLL("Kernel32.dll")
var procGetProcessMemoryInfo = modpsapi.NewProc("GetProcessMemoryInfo")

type windowsC struct {
	collector.BaseCollector
}
type cpuInfo struct {
	Name                      string
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	MaxClockSpeed             uint32
}

func (w windowsC) GetCPUInfo() (collector.Cpu, error) {
	var cpuinfo []cpuInfo
	err := wmi.Query("Select * from Win32_Processor", &cpuinfo)
	if err != nil {
		_ = log.Error("windows Cannt get Cpuinfo " + err.Error())
		return collector.Cpu{}, err
	}
	var info collector.Cpu
	info.Num = len(cpuinfo)
	for _, v := range cpuinfo {
		info.Info = append(info.Info, collector.CpuInfo{
			Thread: int(v.NumberOfLogicalProcessors),
			Name:   v.Name,
			Core:   int(v.NumberOfCores),
			Freq:   float64(v.MaxClockSpeed),
		})

	}
	return info, nil
}
func (w windowsC) GetDiskInfo() ([]collector.Storage, error) {
	type storageInfo struct {
		Name       string
		Size       uint64
		FreeSpace  uint64
		FileSystem string
	}

	var storageinfo []storageInfo
	var loaclStorages []collector.Storage
	err := wmi.Query("Select * from Win32_LogicalDisk", &storageinfo)
	if err != nil {
		return nil, err
	}

	for _, storage := range storageinfo {
		if storage.Size <= 0 {
			continue
		}
		info := collector.Storage{
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

func (w windowsC) GetMemoryInfo() (collector.Memory, error) {
	type query struct {
		FreePhysicalMemory      uint64
		FreeSpaceInPagingFiles  uint64
		TotalVisibleMemorySize  uint64
		SizeStoredInPagingFiles uint64
	}
	var s []query
	err := wmi.Query("Select FreePhysicalMemory,FreeSpaceInPagingFiles,TotalVisibleMemorySize,SizeStoredInPagingFiles from Win32_OperatingSystem", &s)
	if err != nil {
		return collector.Memory{}, err
	}
	if len(s) != 1 {
		return collector.Memory{}, fmt.Errorf("query fail")
	}
	memory := collector.Memory{}
	memory.Physical.Free = int(s[0].FreePhysicalMemory)
	memory.Physical.Total = int(s[0].TotalVisibleMemorySize)
	memory.Swap.Free = int(s[0].FreeSpaceInPagingFiles)
	memory.Swap.Total = int(s[0].SizeStoredInPagingFiles)
	return memory, nil
}
func (w windowsC) GetUptime() (int, error) {
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
func (w windowsC) GetKernel() (string, error) {
	return "Windows NT", nil
}

func (w windowsC) GetSession() (int, error) {
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
	     | Where-Object { $_.name -inotmatch '_total|idle' } `
	    | ForEach-Object {
	         "Process={0,-25} CPU_Usage={1,-12} Memory_Usage_(MB)={2,-16}" -f `
	            $_.Name,$_.PercentProcessorTime,([math]::Round($_.WorkingSetPrivate/1Mb,2))
	     }
	*/
	fmt.Println("wmi")
	t := time.Now()
	type query struct {
		IDProcess            uint32
		PercentProcessorTime uint64
		Name                 string
		WorkingSetPrivate    uint64
	}
	var q []query
	_ = wmi.Query("Select IDProcess,PercentProcessorTime,Name,WorkingSetPrivate from Win32_PerfFormattedData_PerfProc_Process", &q)
	c := 0
	for _, v := range q {
		if getuser(v.IDProcess, false) {
			fmt.Println(v.Name, v.PercentProcessorTime)
			c++
		}
	}
	fmt.Println(time.Now().Sub(t))
	fmt.Println(c)
	fmt.Println("dll")
	t = time.Now()
	h, _ := windows.CreateToolhelp32Snapshot(0x00000002, 0)
	if h < 0 {
		fmt.Println(syscall.GetLastError())
		return
	}
	var ret windows.ProcessEntry32
	ret.Size = uint32(unsafe.Sizeof(ret))
	c = 0
	for windows.Process32Next(h, &ret) == nil {
		if getuser(ret.ProcessID, true) {
			fmt.Println(windows.UTF16ToString(ret.ExeFile[:]))
			c++
		}
	}
	fmt.Println(time.Now().Sub(t))
	fmt.Println(c)
}
func getuser(pid uint32, b bool) bool {
	var token windows.Token
	defer token.Close()
	h2, e := windows.OpenProcess(0x1000, false, pid)
	if e != nil {
		return false
	}
	if b {
		var c, e, k, u windows.Filetime
		if windows.GetProcessTimes(h2, &c, &e, &k, &u) == nil {
			user := float64(u.HighDateTime)*429.4967296 + float64(u.LowDateTime)*1e-7
			kernel := float64(k.HighDateTime)*429.4967296 + float64(k.LowDateTime)*1e-7
			fmt.Println(user, kernel)
		}
	}
	_ = windows.OpenProcessToken(h2, syscall.TOKEN_QUERY, &token)
	//if e == windows.ERROR_INVALID_HANDLE  || e==windows.ERROR_ACCESS_DENIED {
	//	continue
	//}
	tu, e := token.GetTokenUser()
	if e != nil {
		return false
	}
	_, _, _, _ = tu.User.Sid.LookupAccount("")
	return true
}
func init() {
	collector.Register(func() collector.Collector {
		return collector.Collector(windowsC{})
	})
}
