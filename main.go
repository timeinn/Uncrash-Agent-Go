package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/StackExchange/wmi"
	log "github.com/cihub/seelog"
)

func init() {

	log.Info("OS:" + runtime.GOOS)
	log.Info("OS ARCH:" + runtime.GOARCH)
}
func main() {
	defer log.Flush()
	test3()
}

func test3() {
	for {
		type query struct {
			IDProcess            uint32
			PercentProcessorTime uint64
			Name                 string
			WorkingSetPrivate    uint64
			ThreadCount          uint32
		}
		var q []query
		_ = wmi.Query("Select IDProcess,PercentProcessorTime,Name,WorkingSetPrivate,ThreadCount from Win32_PerfFormattedData_PerfProc_Process where Name='Uncrash-Agent-Go'", &q)
		fmt.Println(float64(q[0].PercentProcessorTime) / float64(q[0].ThreadCount))
		time.Sleep(time.Millisecond * 500)
	}
}
