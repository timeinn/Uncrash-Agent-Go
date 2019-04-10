// +build linux

package collector

import (
	"bufio"
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
			storage.Total = fs.Blocks * uint64(fs.Bsize)
			storage.Free = fs.Bfree * uint64(fs.Bsize)
			Storages = append(Storages, storage)
		}
	}
	return Storages, nil
}
