package syslogd

import (
	"net"
	"strings"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"

	"github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/rfc3164"
	"github.com/influxdata/go-syslog/v3/rfc5424"
)

type SyslogServer struct {
	Rfc3164Parser  syslog.Machine
	Rfc5424Parser  syslog.Machine
	Rfc3164Enabled bool
	Debug          bool
}

func NewSyslogServer() *SyslogServer {
	s := &SyslogServer{}
	s.Rfc3164Parser = rfc3164.NewParser(rfc3164.WithBestEffort())
	s.Rfc5424Parser = rfc5424.NewParser(rfc3164.WithBestEffort())
	s.Debug = app.Config.Syslogd.Debug
	return s
}

// HandleRfc3164
// Handling Rfc3164 messages
func (s SyslogServer) HandleRfc3164(remoteaddr net.Addr, data []byte) (*models.Syslog, error) {
	message, err := s.Rfc3164Parser.Parse(data)
	if err != nil {
		return nil, err
	}

	slog := *message.(*rfc3164.SyslogMessage)
	logdata := &models.Syslog{
		Timestamp:       time.Now(),
		Logtype:         "rfc3164",
		MsgID:           "N/A",
		ProcID:          "N/A",
		Appname:         "N/A",
		Hostname:        *slog.Hostname,
		Priority:        int64(*slog.Priority),
		Facility:        int64(*slog.Facility),
		FacilityMessage: *slog.FacilityMessage(),
		Severity:        int64(*slog.Severity),
		SeverityMessage: *slog.SeverityMessage(),
		Version:         0,
		Message:         *slog.Message,
		Tags:            "",
	}
	if slog.Appname != nil {
		logdata.Appname = *slog.Appname
	}

	return logdata, nil
}

// HandleRfc5424
// Handling Rfc5424 messages
func (s SyslogServer) HandleRfc5424(remoteaddr net.Addr, data []byte) (*models.Syslog, error) {
	message, err := s.Rfc5424Parser.Parse(data)
	if err != nil {
		return nil, err
	}
	slog := *message.(*rfc5424.SyslogMessage)
	logdata := &models.Syslog{
		Timestamp:       time.Now(),
		Logtype:         "rfc5424",
		MsgID:           *slog.MsgID,
		ProcID:          *slog.ProcID,
		Appname:         *slog.Appname,
		Hostname:        *slog.Hostname,
		Priority:        int64(*slog.Priority),
		Facility:        int64(*slog.Facility),
		FacilityMessage: *slog.FacilityMessage(),
		Severity:        int64(*slog.Severity),
		SeverityMessage: *slog.SeverityMessage(),
		Version:         int64(slog.Version),
		Message:         *slog.Message,
		Tags:            "",
	}
	return logdata, nil
}

// HandleSyslog
// Handling Text messages
func (s SyslogServer) HandleSyslog(remoteaddr net.Addr, data []byte) {
	defer func() {
		if ret := recover(); ret != nil {
			err, ok := ret.(error)
			if ok {
				log.Error(err)
			}
		}
	}()

	logdata, err := s.HandleRfc3164(remoteaddr, data)
	if err != nil {
		logdata, err = s.HandleRfc5424(remoteaddr, data)
	}

	if err != nil {
		var message = string(data)
		logdata = &models.Syslog{
			Timestamp:       time.Now(),
			Logtype:         "text",
			MsgID:           "N/A",
			ProcID:          "N/A",
			Appname:         "N/A",
			Hostname:        remoteaddr.String(),
			Facility:        3,
			FacilityMessage: "system daemons messages",
			Message:         message,
			Tags:            "",
		}
		switch {
		case strings.Contains(message, "[DEBU]"):
			logdata.Severity = 7
			logdata.SeverityMessage = "debug-level messages"
		case strings.Contains(message, "[ERRO]"):
			logdata.Severity = 3
			logdata.SeverityMessage = "error conditions"
		case strings.Contains(message, "[WARN]"):
			logdata.Severity = 4
			logdata.SeverityMessage = "warning conditions"
		default:
			logdata.Severity = 6
			logdata.SeverityMessage = "informational messages"
		}
	}
	app.PubSyslog(logdata)

}

func (s SyslogServer) StartSyslogServer() error {
	ip := net.ParseIP(app.Config.Syslogd.Host)
	port := app.Config.Syslogd.Port
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return err
	}
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := listener.ReadFrom(data)
		if err != nil {
			log.Error(err)
		}
		var logdata = data[:n]
		if s.Debug {
			log.Info(string(logdata))
		}
		go s.HandleSyslog(remoteAddr, logdata)
	}
}
