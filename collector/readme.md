# collector
## API
### 基础
* GetHostName() (name string, err error) 
* GetOSArch() string 
* GetOS() string 
* GetNetInterfaces() (interfaces []net.Interface, error error) 
* GetOutboundIP() (*net.IP, error) 
* GetOutboundNetInterfaces() (*net.Interface, error) 
### 平台注册
* Register(os string, ext ExtCollector) 
### ExtCollector 平台接口
* GetCPUInfo() ([]Cpu, error) 
* GetDiskInfo() ([]Disk, error) 
* GetMemoryInfo() (Memory, error) 
* GetUptime() (int, error) 
* GetKernel() (string, error) 
* GetSession() (int, error) 
* GetProcess() ([]Process, error) 
* GetInterfacesTraffic(i net.Interface) (*InterfacesTraffic, error) 
* GetLimit() (*Limit, error) 
* GetLoadAvg() ([]float32, error) 
* GetLoad() (*Load, error) 
## 平台
- [x] Ubuntu
- [ ] Windows
# test
- [x] Ubuntu On WSL
