package collector

type Cpu struct {
	Core   int     `json:"core"`
	Name   string  `json:"name"`
	Thread int     `json:"thread"`
	Freq   float64 `json:"freq"`
}

type Disk struct {
	Mount      string `json:"mount"`
	Name       string `json:"name"`
	FileSystem string `json:"file_system"`
	Total      uint64 `json:"total"`
	Free       uint64 `json:"free"`
}
type Memory struct {
	Physical MemoryInfo `json:"physical"`
	Swap     MemoryInfo `json:"swap"`
}
type MemoryInfo struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}
type Process struct {
	User      string  `json:"user"`
	CPU       float64 `json:"cpu"`
	Memory    uint64  `json:"memory"`
	Command   string  `json:"command"`
	IsDefunct bool    `json:"is_defunct"`
}

// 网卡流量
type InterfacesTraffic struct {
	ReceiveBytes, ReceivePackets, TransmitBytes, TransmitPackets uint64
}

type Limit struct {
	Cur uint64
	Max uint64
}
type Load struct {
	Cpu, Io, Idle uint64
}
