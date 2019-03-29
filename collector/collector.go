package collector

import "os"

func GetHostName()(name string,err error)  {
	return os.Hostname()
}

