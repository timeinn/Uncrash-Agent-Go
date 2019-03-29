// +build linux

package collector

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
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
