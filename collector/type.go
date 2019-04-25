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
	Total      int    `json:"total"`
	Free       int    `json:"free"`
}
type Memory struct {
	Physical struct {
		Total int `json:"total"`
		Free  int `json:"free"`
	} `json:"physical"`
	Swap struct {
		Total int `json:"total"`
		Free  int `json:"free"`
	} `json:"swap"`
}
