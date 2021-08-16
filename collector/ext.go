package collector

import (
	"fmt"
	"runtime"
	"sync"
)

type ExtCollector interface {
	GetCPUInfo() ([]Cpu, error)
	GetDiskInfo() ([]Disk, error)
	GetMemoryInfo() (Memory, error)
	GetUptime() (int, error)
	GetKernel() (string, error)
	GetSession() (int, error)
}

type safeExt struct {
	sync.Mutex
	ExtCollector
}

var _extCollector *safeExt = &safeExt{ExtCollector: nil}

func checkExt() {
	if _extCollector == nil || _extCollector.ExtCollector == nil {
		panic("no register " + runtime.GOOS + " collector")
	}
}
func Register(os string, ext ExtCollector) {
	fmt.Println("Register ExtCollector", os)
	_extCollector.Lock()
	_extCollector.ExtCollector = ext
	_extCollector.Unlock()
}
func GetCPUInfo() ([]Cpu, error) {
	checkExt()
	return _extCollector.GetCPUInfo()
}
func GetDiskInfo() ([]Disk, error) {
	checkExt()
	return _extCollector.GetDiskInfo()
}
func GetMemoryInfo() (Memory, error) {
	checkExt()
	return _extCollector.GetMemoryInfo()
}
func GetUptime() (int, error) {
	checkExt()
	return _extCollector.GetUptime()
}
func GetKernel() (string, error) {
	checkExt()
	return _extCollector.GetKernel()
}
func GetSession() (int, error) {
	checkExt()
	return _extCollector.GetSession()
}
