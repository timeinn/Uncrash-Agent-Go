// +build linux

package collector

import (
	"bufio"
	"encoding/binary"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func GetCPUInfo() (Cpu, error) {
	b, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return Cpu{}, err
	}
	content := string(b)
	type _cinfo struct {
		CpuInfo
		PhysicalId int
	}
	var cpuinfoRegExp = regexp.MustCompile("([^:]*?)\\s*:\\s*(.*)$")
	var e = make(map[int]CpuInfo, 0)
	var _c = &_cinfo{}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if len(line) == 0 && i != len(lines)-1 {
			if _, _ok := e[_c.PhysicalId]; !_ok {
				e[_c.PhysicalId] = CpuInfo{
					Core:   _c.Core,
					Thread: _c.Thread,
					Name:   _c.Name,
				}
			}
			_c = &_cinfo{}
			continue
		} else if i == len(lines)-1 {
			continue
		}
		submatches := cpuinfoRegExp.FindStringSubmatch(line)
		key := submatches[1]
		value := submatches[2]
		switch key {
		case "physical id":
			_c.PhysicalId, _ = strconv.Atoi(value)
		case "siblings":
			_c.Thread, _ = strconv.Atoi(value)
		case "model name":
			_c.Name = value
		case "cpu cores":
			_c.Core, _ = strconv.Atoi(value)
		}
	}
	var c = Cpu{}
	for _, v := range e {
		c.Info = append(c.Info, v)
	}
	c.Num = len(e)
	return c, nil
}

func GetDiskInfo() ([]Storage, error) {
	useMounts := false
	_f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		if err != err.(*os.PathError) {
			return nil, err
		}
		useMounts = true
		_f, err = os.Open("/proc/self/mounts")
		if err != nil {
			return nil, err
		}
	}
	s := bufio.NewScanner(_f)
	defer func() {
		_ = _f.Close()
	}()
	Storages := make([]Storage, 0)
	devSet := mapset.NewSet()
	for s.Scan() {
		var storage Storage
		var path string
		lines := strings.Fields(s.Text())
		if useMounts {
			if !strings.Contains(lines[0], "/dev/") {
				continue
			}
			storage.Name = lines[0]
			storage.FileSystem = lines[2]
			path = lines[1]
		} else {
			if !strings.Contains(lines[8], "/dev/") {
				continue
			}
			storage.Name = lines[8]
			storage.FileSystem = lines[7]
			path = lines[4]
		}
		if devSet.Add(storage.Name) {
			fs := syscall.Statfs_t{}
			err = syscall.Statfs(path, &fs)
			if err != nil {
				continue
			}
			storage.Total = int(fs.Blocks * uint64(fs.Bsize))
			storage.Free = int(fs.Bfree * uint64(fs.Bsize))
			Storages = append(Storages, storage)
		}
	}
	return Storages, nil
}
func GetMemoryInfo() (Memory, error) {
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err != nil {
		return Memory{}, err
	}
	memory := Memory{}
	memory.Physical.Free = int(sysInfo.Freeram)
	memory.Physical.Total = int(sysInfo.Totalram)
	memory.Swap.Free = int(sysInfo.Freeswap)
	memory.Swap.Total = int(sysInfo.Totalswap)
	return memory, nil
}

func GetUptime() (int, error) {
	b, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	if t, err := strconv.ParseFloat(strings.Split(string(b), " ")[0], 64); err != nil {
		return 0, err
	} else {
		return int(t), nil
	}
}

func GetKernel() (string, error) {
	if b, err := ioutil.ReadFile("/proc/sys/kernel/osrelease"); err != nil {
		{
			return "linux", err
		}
	} else {
		return string(b), nil
	}
}

func GetSession() (int, error) {
	type ExitStatus struct {
		X__e_termination int16
		X__e_exit        int16
	}
	type TimeVal struct {
		Sec  int32
		Usec int32
	}
	type Utmp struct {
		Type      int16      //2
		Pad_cgo_0 [2]byte    //2
		Pid       int32      //4
		Line      [32]byte   //32
		Id        [4]byte    //4
		User      [32]byte   //32
		Host      [256]byte  //256
		Exit      ExitStatus //4
		Session   int32      //4
		Tv        TimeVal    //8
		AddrV6    [4]int32   //16
		Unused    [20]byte   //20
	}
	file, err := os.Open("/var/run/utmp")
	defer file.Close()
	if err != nil {
		return 0, err
	}
	var v []Utmp
	for {
		var u Utmp
		if err := binary.Read(file, binary.LittleEndian, &u); err != nil {
			break
		}
		if u.Type == 7 {
			v = append(v, u)
		}
	}
	return len(v), nil
}
