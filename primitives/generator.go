package primitives

import (
	"syscall"
	"fmt"
)

func NewGenerator(prefix string) func() string {
	si := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(si)
	if err != nil {
		panic("Commander, we have a problem. syscall.Sysinfo:" + err.Error())
	}

	counter, uptime := 0, si.Uptime
	return func() string {
		counter++
		return fmt.Sprintf("%s_%d_%d", prefix, uptime, counter)
	}
}

