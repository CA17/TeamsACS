// +build !windows

package log

import (
	"log/syslog"

	"github.com/op/go-logging"
)

// SyslogBackend is a simple logger to syslog backend. It automatically maps
// the internal log levels to appropriate syslog log levels.
type SyslogBackend struct {
	Writer *syslog.Writer
}

// NewSyslogBackend connects to the syslog daemon using UNIX sockets with the
// given prefix. If prefix is not given, the prefix will be derived from the
// launched command.
func NewSyslogBackend(prefix string, syslogaddr string, priority syslog.Priority) (b *SyslogBackend, err error) {
	var w *syslog.Writer
	if syslogaddr != "" {
		w, err = syslog.Dial("udp", syslogaddr, priority, prefix)
	} else {
		w, err = syslog.Dial("", "", priority, prefix)
	}
	return &SyslogBackend{w}, err
}

// NewSyslogBackendPriority is the same as NewSyslogBackend, but with custom
// syslog priority, like syslog.LOG_LOCAL3|syslog.LOG_DEBUG etc.
func NewSyslogBackendPriority(prefix string, priority syslog.Priority) (b *SyslogBackend, err error) {
	var w *syslog.Writer
	w, err = syslog.New(priority, prefix)
	return &SyslogBackend{w}, err
}

// Log implements the Backend interface.
func (b *SyslogBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	line := rec.Formatted(calldepth + 1)
	switch level {
	case logging.CRITICAL:
		return b.Writer.Crit(line)
	case logging.ERROR:
		return b.Writer.Err(line)
	case logging.WARNING:
		return b.Writer.Warning(line)
	case logging.NOTICE:
		return b.Writer.Notice(line)
	case logging.INFO:
		return b.Writer.Info(line)
	case logging.DEBUG:
		return b.Writer.Debug(line)
	default:
	}
	panic("unhandled log level")
}
