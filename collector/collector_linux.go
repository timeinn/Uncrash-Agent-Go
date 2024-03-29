package collector

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
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

//Linux获取CPU信息
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

var ignoreFilesystem = mapset.NewSet("squashfs")

// Linux获取挂在的硬盘信息
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
		if useMounts {
			lines := strings.Fields(s.Text())
			if !strings.Contains(lines[0], "/dev/") && !strings.Contains(lines[0], "overlay") {
				continue
			}

			if ignoreFilesystem.Contains(lines[2]) {
				continue
			}
			storage.Name = lines[0]
			storage.FileSystem = lines[2]
			path = lines[1]
		} else {
			mounts := strings.Split(s.Text(), "-")
			tli := strings.Fields(mounts[1])
			preli := strings.Fields(mounts[0])
			if !strings.Contains(tli[1], "/dev/") && !strings.Contains(tli[1], "overlay") {
				continue
			}
			if ignoreFilesystem.Contains(tli[0]) {
				continue
			}
			storage.Name = tli[1]
			storage.FileSystem = tli[0]
			path = preli[4]
		}
		if devSet.Add(storage.Name) {

			fs := syscall.Statfs_t{}
			err = syscall.Statfs(path, &fs)
			if err != nil {
				continue
			}
			storage.Mount = path
			storage.Total = fs.Blocks * uint64(fs.Bsize) //byte
			storage.Free = fs.Bfree * uint64(fs.Bsize)
			Storages = append(Storages, storage)
		}
	}
	defer func() {
		devSet.Clear()
	}()
	return Storages, nil
}

//Linux内存信息
func (l *linuxCo) GetMemoryInfo() (Memory, error) {
	sysInfo := new(syscall.Sysinfo_t)
	if err := syscall.Sysinfo(sysInfo); err != nil {
		return Memory{}, err
	}
	memory := Memory{}
	memory.Physical.Free = sysInfo.Freeram //byte
	memory.Physical.Total = sysInfo.Totalram
	memory.Swap.Free = sysInfo.Freeswap
	memory.Swap.Total = sysInfo.Totalswap
	return memory, nil
}

func (l *linuxCo) GetUptime() (int, error) {
	b, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	if t, err := strconv.ParseFloat(strings.Fields(string(b))[0], 64); err != nil {
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

//获取用户会话数
//读取 /var/run/utmp 文件不存在返回0
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

type process struct {
	pid       int
	comm      string
	utime     int
	stime     int
	cutime    int
	cstime    int
	starttime int
	rss       uint64
	uid       string
	defunct   bool
}

func (p *process) getInfo() error {
	if err := p.getCmd(); err != nil {
		return err
	}
	if f, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", p.pid)); err != nil {
		return err
	} else {
		stat := strings.Fields(string(f))
		if len(stat) >= 24 {
			p.defunct = stat[2] == "Z"
			p.utime, err = strconv.Atoi(stat[13])
			if err != nil {
				return err
			}
			p.stime, err = strconv.Atoi(stat[14])
			if err != nil {
				return err
			}
			p.cutime, err = strconv.Atoi(stat[15])
			if err != nil {
				return err
			}
			p.cstime, err = strconv.Atoi(stat[16])
			if err != nil {
				return err
			}
			p.starttime, err = strconv.Atoi(stat[21])
			if err != nil {
				return err
			}
		}
		return nil
	}
}
func (p *process) getCmd() error {
	if f, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", p.pid)); err != nil {
		return err
	} else {
		//p.comm = string(bytes.Trim(f, "\x00"))
		p.comm = string(f)
		return nil
	}
}
func (p *process) getUidAndRam() error {
	if f, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", p.pid)); err != nil {
		return err
	} else {
		uid := false
		ram := false
		for _, v := range strings.Split(string(f), "\n") {
			if v == "" {
				continue
			}
			if sp := strings.Fields(v); len(sp) >= 2 {
				switch strings.TrimSpace(sp[0]) {
				case "Uid:":
					p.uid = strings.TrimSpace(sp[1])
					uid = true
				case "VmRSS":
					ram = true
					p.rss, _ = strconv.ParseUint(strings.TrimSpace(sp[1]), 10, 64)
				}
			}
			if uid && ram {
				break
			}
		}
		return nil
	}
}

const hz float64 = 100

// 计算CPU使用率
func (p *process) getCPU(utime int) float64 {
	total_time := p.utime + p.stime + p.cstime + p.cutime
	seconds := float64(utime) - (float64(p.starttime) / hz)
	if seconds <= 0 || total_time <= 0 {
		return 0
	}
	return 100 * ((float64(total_time) / hz) / seconds)
}

func (l *linuxCo) GetProcess() ([]Process, error) {
	selfPid := os.Getpid()
	fs, err := ioutil.ReadDir("/proc/")
	ut, uterr := l.GetUptime()
	if err != nil {
		return nil, err
	}
	ps := make([]Process, 0)
	for _, v := range fs {
		if v.IsDir() {
			pid, err := strconv.Atoi(v.Name())
			if err != nil || pid == selfPid {
				continue
			}
			var p = &process{pid: pid}
			p.getUidAndRam()
			if p.getInfo() != nil || uterr != nil {
				continue
			}
			pa := Process{}
			pa.CPU = p.getCPU(ut)
			pa.Command = p.comm
			pa.Memory = p.rss
			if puser, err := user.LookupId(p.uid); err == nil {
				pa.User = puser.Username
			}
			pa.IsDefunct = p.defunct
			ps = append(ps, pa)
		}
	}
	return ps, nil
}
func (l *linuxCo) GetInterfacesTraffic(i net.Interface) (*InterfacesTraffic, error) {
	//interfaceName := i.Name
	if f, err := ioutil.ReadFile("/proc/net/dev"); err != nil {
		return nil, err
	} else {
		for _, info := range strings.Split(string(f), "\n")[2:] {
			dataLine := strings.Fields(info)
			if len(dataLine) < 17 {
				continue
			}
			if i.Name+":" == dataLine[0] {
				traffic := &InterfacesTraffic{}
				traffic.ReceiveBytes, _ = strconv.ParseUint(dataLine[1], 10, 64)
				traffic.ReceivePackets, _ = strconv.ParseUint(dataLine[2], 10, 64)
				traffic.TransmitBytes, _ = strconv.ParseUint(dataLine[8], 10, 64)
				traffic.TransmitPackets, _ = strconv.ParseUint(dataLine[9], 10, 64)
				return traffic, nil
			}

		}
	}
	return nil, ErrorNotFound
}
func (l *linuxCo) GetLimit() (*Limit, error) {
	data, err := ioutil.ReadFile("/proc/sys/fs/file-nr")
	if err != nil {
		return nil, err
	}
	limit := &Limit{}
	line := strings.Fields(string(data))
	if len(line) < 3 {
		return limit, nil
	}
	limit.Cur, _ = strconv.ParseUint(line[0], 10, 64)
	limit.Max, _ = strconv.ParseUint(line[2], 10, 64)
	return limit, nil
}
func (l *linuxCo) GetLoadAvg() ([]float32, error) {
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return []float32{0, 0, 0}, nil
	}
	avg := make([]float32, 0)
	for _, v := range strings.Fields(string(data))[:3] {
		l, _ := strconv.ParseFloat(v, 32)
		avg = append(avg, float32(l))
	}
	return avg, nil
}
func (l *linuxCo) GetLoad() (*Load, error) {
	data, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	info := strings.Fields(strings.Split(string(data), "\n")[0])
	load := &Load{}
	l1, _ := strconv.ParseUint(info[1], 10, 64)
	l2, _ := strconv.ParseUint(info[2], 10, 64)
	l3, _ := strconv.ParseUint(info[3], 10, 64)
	l4, _ := strconv.ParseUint(info[4], 10, 64)
	load.Cpu = l1 + l2 + l3 + l4
	load.Io = l3 + l4
	load.Idle = l3
	return load, nil
}

//注册实现
func init() {
	Register("linux", &linuxCo{})
}
