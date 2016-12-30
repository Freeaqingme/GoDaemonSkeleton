package log

import (
	logging "github.com/op/go-logging"
	"log"
	"log/syslog"
	"os"
	"syscall"
)

type Logger struct {
	*logging.Logger
}

var (
	Log *Logger
	path string
)

func Open(name, logLevelStr string) *Logger {
	logLevel, err := logging.LogLevel(logLevelStr)
	if err != nil {
		log.Fatal("Invalid log level specified")
	}

	Log = &Logger { logging.MustGetLogger(name) }

	var formatStdout = logging.MustStringFormatter(
		"%{color}%{time:2006-01-02T15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{message}",
	)
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(stdout, formatStdout)
	stdoutLeveled := logging.AddModuleLevel(formatter)
	stdoutLeveled.SetLevel(logLevel, "")
	syslogBackend, err := logging.NewSyslogBackendPriority("cluegetter", syslog.LOG_MAIL)
	if err != nil {
		Log.Fatal(err)
	}

	logging.SetBackend(syslogBackend, stdoutLeveled)

	return Log
}

func Reopen() {
	if path == "" {
		Log.Notice("Asked to reopen logs but running in foreground. Ignoring.")
		return
	}
	Log.Notice("Reopening log file per IPC request...")
	LogRedirectStdOutToFile(path)
	Log.Notice("Reopened log file per IPC request")
}

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
