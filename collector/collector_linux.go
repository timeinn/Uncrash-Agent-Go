package collector

import (
	"bufio"
	"encoding/binary"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	mapset "github.com/deckarep/golang-set"
)

type linuxCo struct {
}
type _cinfo struct {
	Cpu
	PhysicalId int
}

func (l *linuxCo) GetCPUInfo() ([]Cpu, error) {
	b, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	content := string(b)
	var cpuinfoRegExp = regexp.MustCompile(`([^:]*?)\s*:\s*(.*)$`)
	var e = make(map[int]Cpu)
	var _c = &_cinfo{}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if len(line) == 0 && i != len(lines)-1 {
			if _, _ok := e[_c.PhysicalId]; !_ok {
				e[_c.PhysicalId] = Cpu{
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
	c := make([]Cpu, 0)
	for _, v := range e {
		c = append(c, v)
	}
	defer func() {
		for k := range e {
			delete(e, k)
		}
	}()
	return c, nil
}

func (l *linuxCo) GetDiskInfo() ([]Disk, error) {
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
	Storages := make([]Disk, 0)
	devSet := mapset.NewSet()
	for s.Scan() {
		var storage Disk
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
			storage.Mount = path
			storage.Total = int(fs.Blocks * uint64(fs.Bsize))
			storage.Free = int(fs.Bfree * uint64(fs.Bsize))
			Storages = append(Storages, storage)
		}
	}
	defer func() {
		devSet.Clear()
	}()
	return Storages, nil
}

func (l *linuxCo) GetMemoryInfo() (Memory, error) {
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

func (l *linuxCo) GetUptime() (int, error) {
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

func (l *linuxCo) GetKernel() (string, error) {
	if b, err := ioutil.ReadFile("/proc/sys/kernel/osrelease"); err != nil {
		return "linux", err
	} else {
		return string(b), nil
	}
}

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

func (l *linuxCo) GetSession() (int, error) {

	file, err := os.Open("/var/run/utmp")
	if err != nil {
		return 0, err
	}
	defer file.Close()
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

func init() {
	Register("linux", &linuxCo{})
}
