/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package syslogd

import (
	"net"
	"strings"
	"time"

	"github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/rfc3164"
	"github.com/influxdata/go-syslog/v3/rfc5424"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

type SyslogServer struct {
	Manager        *models.ModelManager
	Rfc3164Parser  syslog.Machine
	Rfc5424Parser  syslog.Machine
	Rfc3164Enabled bool
	Debug          bool
}

func NewSyslogServer(manager *models.ModelManager) *SyslogServer {
	s := &SyslogServer{Manager: manager}
	s.Rfc3164Parser = rfc3164.NewParser(rfc3164.WithBestEffort())
	s.Rfc5424Parser = rfc5424.NewParser(rfc3164.WithBestEffort())
	s.Debug = s.Manager.Config.Syslogd.Debug
	return s
}

// HandleRfc3164
// Handling Rfc3164 messages
func (s SyslogServer) HandleRfc3164(remoteaddr net.Addr, data []byte) {
	defer func() {
		if ret := recover(); ret != nil {
			err, ok := ret.(error)
			if ok {
				log.Error(err)
			}
		}
	}()
	message, err := s.Rfc3164Parser.Parse(data)
	if err != nil {
		s.HandleText(remoteaddr, data)
		return
	}

	slog := *message.(*rfc3164.SyslogMessage)
	attrs := map[string]interface{}{
		"Message":         *slog.Message,
		"Facility":        *slog.Facility,
		"FacilityMessage": *slog.FacilityMessage(),
		"Severity":        *slog.Severity,
		"SeverityMessage": *slog.SeverityMessage(),
		"Priority":        *slog.Priority,
		"Timestamp":       *slog.Timestamp,
		"Hostname":        *slog.Hostname,
		"Appname":         "N/A",
	}
	if slog.Appname != nil {
		attrs["Appname"]=  *slog.Appname
	}

	s.Manager.GetOpsManager().AddSyslog(&models.Syslog{
		Logtype:   "rfc3164",
		Attrs:attrs,
		Timestamp: time.Now(),
	})
}

// HandleRfc5424
// Handling Rfc5424 messages
func (s SyslogServer) HandleRfc5424(remoteaddr net.Addr, data []byte) {
	defer func() {
		if ret := recover(); ret != nil {
			err, ok := ret.(error)
			if ok {
				log.Error(err)
			}
		}
	}()
	message, err := s.Rfc5424Parser.Parse(data)
	if err != nil {
		s.HandleText(remoteaddr, data)
		return
	}
	slog := *message.(*rfc5424.SyslogMessage)
	s.Manager.GetOpsManager().AddSyslog(&models.Syslog{
		Logtype: "rfc5424",
		Attrs: map[string]interface{}{
			"Message":         *slog.Message,
			"Facility":        *slog.Facility,
			"FacilityMessage": *slog.FacilityMessage(),
			"Severity":        *slog.Severity,
			"SeverityMessage": *slog.SeverityMessage(),
			"Priority":        *slog.Priority,
			"Timestamp":       *slog.Timestamp,
			"Hostname":        *slog.Hostname,
			"Appname":         *slog.Appname,
			"ProcID":          *slog.ProcID,
			"MsgID":           *slog.MsgID,
			"Version":         slog.Version,
		},
		Timestamp: time.Now(),
	})
}

// HandleText
// Handling Text messages
func (s SyslogServer) HandleText(remoteaddr net.Addr, data []byte) {
	defer func() {
		if ret := recover(); ret != nil {
			err, ok := ret.(error)
			if ok {
				log.Error(err)
			}
		}
	}()
	var message = string(data)
	logtext := &models.Syslog{
		Logtype: "text",
		Attrs: map[string]interface{}{
			"Hostname":        remoteaddr.String(),
			"Message":         message,
			"Facility":        3,
			"FacilityMessage": "system daemons messages",
			"Timestamp":       time.Now(),
			"Appname":         "N/A",
		},
		Timestamp: time.Now(),
	}
	switch {
	case strings.Contains(message, "[DEBU]"):
		logtext.Attrs["Severity"] = 7
		logtext.Attrs["SeverityMessage"] = "debug-level messages"
	case strings.Contains(message, "[ERRO]"):
		logtext.Attrs["Severity"] = 3
		logtext.Attrs["SeverityMessage"] = "error conditions"
	case strings.Contains(message, "[WARN]"):
		logtext.Attrs["Severity"] = 4
		logtext.Attrs["SeverityMessage"] = "warning conditions"
	default:
		logtext.Attrs["Severity"] = 6
		logtext.Attrs["SeverityMessage"] = "informational messages"
	}

	s.Manager.GetOpsManager().AddSyslog(logtext)
}

func (s SyslogServer) StartRfc3164() error {
	ip := net.ParseIP(s.Manager.Config.Syslogd.Host)
	port := s.Manager.Config.Syslogd.Rfc3164Port
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
		go s.HandleRfc3164(remoteAddr, logdata)
	}
}

func (s SyslogServer) StartRfc5424() error {
	ip := net.ParseIP(s.Manager.Config.Syslogd.Host)
	port := s.Manager.Config.Syslogd.Rfc5424Port
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return err
	}
	for {
		data := make([]byte, 2048)
		n, remoteAddr, err := listener.ReadFrom(data)
		if err != nil {
			log.Error(err)
		}
		var logdata = data[:n]
		if s.Debug {
			log.Info(string(logdata))
		}
		go s.HandleRfc5424(remoteAddr, logdata)
	}
}

func (s SyslogServer) StartTextlog() error {
	ip := net.ParseIP(s.Manager.Config.Syslogd.Host)
	port := s.Manager.Config.Syslogd.TextlogPort
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		return err
	}
	for {
		data := make([]byte, 8912)
		n, remoteAddr, err := listener.ReadFrom(data)
		if err != nil {
			log.Error(err)
		}
		var logdata = data[:n]
		if s.Debug {
			log.Info(string(logdata))
		}
		go s.HandleText(remoteAddr, logdata)
	}
}
