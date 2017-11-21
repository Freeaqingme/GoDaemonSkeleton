package log

import "syscall"

func stdFdToLogFile(logFileFd int) {
	syscall.Dup3(logFileFd, 1, 0)
	syscall.Dup3(logFileFd, 2, 0)
}
