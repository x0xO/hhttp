// +build !windows

package hsyscall

import "syscall"

func RlimitStack(maxWorkers int) int {
	var rLimit syscall.Rlimit
	syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if uint64(maxWorkers) > rLimit.Cur {
		maxWorkers = int(float64(rLimit.Cur) * 0.7)
	}

	return maxWorkers
}
