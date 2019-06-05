package main

import (
	"fmt"
	"github.com/StackExchange/wmi"
	log "github.com/cihub/seelog"
	"golang.org/x/sys/windows"
	"runtime"
	"syscall"
	"time"
)

func init() {

	log.Info("OS:" + runtime.GOOS)
	log.Info("OS ARCH:" + runtime.GOARCH)
}
func main() {
	defer log.Flush()
	test3()
}
func test2()  {
	for {
		getuser(8876,true)
		time.Sleep(time.Millisecond * 500)
	}
}
func getuser(pid uint32,b bool) bool {
	var token windows.Token
	defer token.Close()
	h2,e:=windows.OpenProcess(0x1000,false,pid)
	if e!=nil{
		return false
	}
	if b {
		var	c,e,k,u windows.Filetime
		if windows.GetProcessTimes(h2,&c,&e,&k,&u)==nil{
			user := float64(u.HighDateTime)*429.4967296 + float64(u.LowDateTime)*1e-7
			kernel := float64(k.HighDateTime)*429.4967296 + float64(k.LowDateTime)*1e-7
			fmt.Println( user,kernel)
			fmt.Println( user+kernel)
		}
	}
	_=windows.OpenProcessToken(h2,syscall.TOKEN_QUERY,&token)
	//if e == windows.ERROR_INVALID_HANDLE  || e==windows.ERROR_ACCESS_DENIED {
	//	continue
	//}
	tu,e:=token.GetTokenUser()
	if e!=nil{
		return false
	}
	_, _, _, _ = tu.User.Sid.LookupAccount("")
	return true
}
func test3()  {
	for {
		type query struct {
			IDProcess            uint32
			PercentProcessorTime uint64
			Name                 string
			WorkingSetPrivate    uint64
			ThreadCount uint32
		}
		var q []query
		_ = wmi.Query("Select IDProcess,PercentProcessorTime,Name,WorkingSetPrivate,ThreadCount from Win32_PerfFormattedData_PerfProc_Process where Name='Uncrash-Agent-Go'", &q)
		fmt.Println(float64(q[0].PercentProcessorTime )/float64(q[0].ThreadCount))
		time.Sleep(time.Millisecond * 500)
	}
}
