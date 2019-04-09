package collector

type Message struct {
	Cpu Cpu `json:"cpu"`
}
type Cpu struct {
	Num  int       `json:"num"`
	Info []CpuInfo `json:"info"`
}
type CpuInfo struct {
	Core   int    `json:"core"`
	Name   string `json:"name"`
	Thread int    `json:"thread"`
}
type Interfaces struct {
	Name  string   `json:"name"`
	Addrs []string `json:"addrs"`
}

type Storage struct {
	Name       string `json:"name"`
	FileSystem string `json:"file_system"`
	Total      uint64 `json:"total"`
	Free       uint64 `json:"free"`
}
