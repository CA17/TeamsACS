// +build !windows

package log

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/op/go-logging"
)

func SetupSyslog(level logging.Level, syslogaddr string, module string) logging.LeveledBackend {
	var format = logging.MustStringFormatter(
		`%{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
	)
	backend, err := NewSyslogBackend("", syslogaddr, syslog.LOG_INFO)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backend2Formatter := logging.NewBackendFormatter(backend, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(level, module)
	return backend1Leveled
}
