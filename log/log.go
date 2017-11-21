package log

import (
	logging "github.com/op/go-logging"
	"log"
	"log/syslog"
	"os"
)

type Logger struct {
	*logging.Logger
}

var (
	Log *Logger
	path string
)

func Open(name, logLevelStr string, prio syslog.Priority) *Logger {
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
	syslogBackend, err := logging.NewSyslogBackendPriority(name, prio)
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


