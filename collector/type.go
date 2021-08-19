package collector

import "net"

type Cpu struct {
	Core   int     `json:"core"`
	Name   string  `json:"name"`
	Thread int     `json:"thread"`
	Freq   float64 `json:"freq"`
}
type Interfaces struct {
	Name  string       `json:"name"`
	Addrs []net.IPAddr `json:"addrs"`
}

type Disk struct {
	Mount      string `json:"mount"`
	Name       string `json:"name"`
	FileSystem string `json:"file_system"`
	Total      int    `json:"total"`
	Free       int    `json:"free"`
}
type Memory struct {
	Physical MemoryInfo `json:"physical"`
	Swap     MemoryInfo `json:"swap"`
}
type MemoryInfo struct {
	Total int `json:"total"`
	Free  int `json:"free"`
}
type Process struct {
	User      string  `json:"user"`
	CPU       float64 `json:"cpu"`
	Memory    uint    `json:"memory"`
	Command   string  `json:"command"`
	IsDefunct bool    `json:"is_defunct"`
}

// 网卡流量
type InterfacesTraffic struct {
	ReceiveBytes, ReceivePackets, TransmitBytes, TransmitPackets uint64
}
