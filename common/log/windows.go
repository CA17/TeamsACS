// +build windows

package log

import "github.com/op/go-logging"

func SetupSyslog(level logging.Level, syslogaddr string, module string) logging.LeveledBackend {
	return nil
}
