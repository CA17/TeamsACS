package models

import (
	"time"
)

//  系统数据模型

type SysConfig struct {
	Type   string `gorm:"primaryKey" json:"type"`
	Name   string `gorm:"primaryKey" json:"name"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type Syslog struct {
	ID              int64     `json:"id,string"`
	Timestamp       time.Time `json:"timestamp"`
	Logtype         string    `json:"logtype"`
	MsgID           string    `json:"msg_id,omitempty"`
	ProcID          string    `json:"proc_id,omitempty"`
	Appname         string    `json:"appname,omitempty"`
	Hostname        string    `json:"hostname,omitempty"`
	Priority        int64     `json:"priority,omitempty"`
	Facility        int64     `json:"facility,omitempty"`
	FacilityMessage string    `json:"facility_message,omitempty"`
	Severity        int64     `json:"severity,omitempty"`
	SeverityMessage string    `json:"severity_message,omitempty"`
	Version         int64     `json:"version,omitempty"`
	Message         string    `json:"message"`
	Tags            string    `json:"tags,omitempty"`
}
