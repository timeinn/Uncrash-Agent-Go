package collector

import (
	"fmt"
	"runtime"
	"sync"
)

// 根据各自平台获取信息的接口
type ExtCollector interface {
	GetCPUInfo() ([]Cpu, error)
	GetDiskInfo() ([]Disk, error)
	GetMemoryInfo() (Memory, error)
	GetUptime() (int, error)
	GetKernel() (string, error)
	GetSession() (int, error)
	GetProcess() ([]Process, error)
}

type safeExt struct {
	sync.Mutex
	ExtCollector
}

//保存平台的实现的接口
var _extCollector *safeExt = &safeExt{ExtCollector: nil}

// 检查接口实现
func checkExt() {
	if _extCollector == nil || _extCollector.ExtCollector == nil {
		panic("no register " + runtime.GOOS + " collector")
	}
}

// 注册平台实现的接口 以最后一次调用为准
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
func GetProcess() ([]Process, error) {
	checkExt()
	return _extCollector.GetProcess()
}
