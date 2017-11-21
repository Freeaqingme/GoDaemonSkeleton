// +build !arm64

package log

import (
	"os"
	"syscall"
)

// If there are any existing fd's (e.g. we're reopening logs), we rely
// on garbage collection to clean them up for us.
func LogRedirectStdOutToFile(logPath string) {
	path = logPath
	if logPath == "" {
		Log.Fatal("Log Path not set")
	}

	logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		Log.Fatal(err)
	}

	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
}