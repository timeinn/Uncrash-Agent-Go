package main

import (
	"fmt"
	"runtime"

	"github.com/TimeInn/Uncrash-Agent-Go/collector"
	log "github.com/cihub/seelog"
)

func init() {
	a, err := collector.GetNetInterfaces()
	fmt.Println(err)
	fmt.Println(a)
	log.Info("OS:" + runtime.GOOS)
	log.Info("OS ARCH:" + runtime.GOARCH)
	fmt.Println(collector.GetOutboundIP())
}
func main() {
	defer log.Flush()

}
