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
