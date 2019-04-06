package main

import (
	"github.com/TimeInn/Uncrash-Agent-Go/collector"
	log "github.com/cihub/seelog"
	"runtime"
)

func init() {

	log.Info("OS:" + runtime.GOOS)
	log.Info("OS ARCH:" + runtime.GOARCH)
}
func main() {
	defer log.Flush()
	cpu, _ := collector.GetCPUInfo()
	log.Info(cpu)
	NetInterfaces, _ := collector.GetNetInterfaces()
	log.Info(NetInterfaces)
}
