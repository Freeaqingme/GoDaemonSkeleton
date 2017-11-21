// +build !arm64

package log

import "syscall"

func stdFdToLogFile(logFileFd int) {
	syscall.Dup2(logFileFd, 1)
	syscall.Dup2(logFileFd, 2)
}
