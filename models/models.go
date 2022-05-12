package models

import (
	"time"
)

type SysConfig struct {
	Type   string `gorm:"primaryKey" json:"type"`
	Name   string `gorm:"primaryKey" json:"name"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// NetDevice Device Data Model
type NetDevice struct {
	ID              int64     `json:"id,string" form:"id"`                      // primary key ID
	DevType         string    `json:"dev_type" form:"dev_type"`                 // device type: server | switch | router
	Sn              string    `gorm:"uniqueIndex" json:"sn" form:"sn"`          // device sn
	Name            string    `json:"name" form:"name"`                         // device name
	SystemName      string    `json:"system_name" form:"system_name"`           // Device system name
	ArchName        string    `json:"arch_name" form:"arch_name"`               // Device Architecture
	PublicIpaddr    string    `json:"public_ipaddr" form:"public_ipaddr"`       // Device public IP
	TunnelIpaddr    string    `json:"tunnel_ipaddr" form:"tunnel_ipaddr"`       // Device Tunnel IP
	Model           string    `json:"model" form:"model"`                       // Device model
	VendorCode      string    `json:"vendor_code" form:"vendor_code"`           // device vendor code
	SoftwareVersion string    `json:"software_version" form:"software_version"` // Device firmware version
	HardwareVersion string    `json:"hardware_version" form:"hardware_version"` // Device hardware version
	Oui             string    `json:"oui" form:"oui"`                           // device oui
	Manufacturer    string    `json:"manufacturer" form:"manufacturer"`         // device manufactory
	ProductClass    string    `json:"product_class" form:"product_class"`       // device type
	CwmpUrl         string    `json:"cwmp_url" form:"cwmp_url"`
	CwmpLastInform  time.Time `json:"cwmp_last_inform" `                // CWMP Last detection time
	Tags            string    `json:"tags" form:"tags"`                 // device tags
	Uptime          int64     `json:"uptime" form:"uptime"`             // device UpTime
	MemoryTotal     int64     `json:"memory_total" form:"memory_total"` // total memory
	MemoryFree      int64     `json:"memory_free" form:"memory_free"`   // memory available
	CPUUsage        int64     `json:"cpu_usage" form:"cpu_usage"`       // CPE percentage
	Remark          string    `json:"remark" form:"remark"`             // remark
	CreatedAt       time.Time `json:"created_at" `
	UpdatedAt       time.Time `json:"updated_at"`
}

// NetDeviceParam Device parameter information
type NetDeviceParam struct {
	ID        int64     `json:"id" form:"id"`       // primary key ID
	Sn        string    `json:"sn" form:"sn"`       // device serial number
	Name      string    `json:"name" form:"name"`   // parameter name
	Value     string    `json:"value" form:"value"` // parameter value
	LastFetch time.Time `json:"last_fetch" `        // Last fetch time
}

type CwmpParamFetchTask struct {
	ID           int64     `json:"id"`                               // primary key ID
	Oui          string    `json:"oui" form:"oui"`                   // Device OUI
	Manufacturer string    `json:"manufacturer" form:"manufacturer"` // device manufactory
	ProductClass string    `json:"productClass" form:"productClass"` // device type
	Name         string    `json:"name" form:"name"`                 // param name
	Interval     int64     `json:"interval" form:"interval"`         // Intervals time : seconds
	LastFetch    time.Time `json:"last_fetch" `                      //  Last fetch time
}

// CwmpFile  Cwmp Device configuration file or firmware file
type CwmpFile struct {
	ID           string    `gorm:"primaryKey" json:"id" form:"id"`     // primary key ID
	Name         string    `json:"name" form:"name"`                   // Cwmp file name
	Oui          string    `json:"oui" form:"oui"`                     // device OUI
	Manufacturer string    `json:"manufacturer" form:"manufacturer"`   // device manufactory
	ProductClass string    `json:"product_class" form:"product_class"` // device product class
	Version      string    `json:"version" form:"version"`             // device firmware version
	FileType     string    `json:"file_type" form:"file_type"`         // file type
	FilePath     string    `json:"file_path" form:"file_path"`         // file save path
	Format       string    `json:"format" form:"format"`               // format:  xml | alter | json
	Timeout      int64     `json:"timeout" form:"timeout"`             // timeout: second
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CwmpDownloadTask Cwmp Download tasks
type CwmpDownloadTask struct {
	ID           int64     `gorm:"primaryKey" json:"id,string"` // primary key ID
	Sn           string    `json:"sn"`                          // device sn
	SessionId    string    `json:"session_id" `                 // session id
	Name         string    `json:"name"`                        // task name
	FileType     string    `json:"file_type"`                   // cwmp file type
	FilePath     string    `json:"file_path"`                   // file save path
	Filename     string    `json:"filename"`                    // file name
	Status       string    `json:"status"`                      // statue: initialize | success | error | timeout
	FaultCode    int       `json:"fault_code"`                  // FaultCode
	FaultString  string    `json:"fault_string"`                // FaultString
	Timeout      int64     `json:"timeout" `                    // execution time timeout: second
	StartTime    time.Time `json:"start_time"`                  // execution start time
	CompleteTime time.Time `json:"complete_time"`               // Complete time
	CreatedAt    time.Time `json:"created_at"`
}

// CwmpMessageTask Cwmp message tasks
type CwmpMessageTask struct {
	ID        int64     `gorm:"primaryKey" json:"id,string"` // message ID
	Sn        string    `json:"sn"`                          // device sn
	Name      string    `json:"name"`                        // message name
	Message   string    `json:"message"`                     // message content
	CreatedAt time.Time `json:"created_at"`
}

type CwmpParamNameList struct {
	Sn        string    `gorm:"primaryKey" json:"sn" form:"sn"`     // device sn
	Name      string    `gorm:"primaryKey" json:"name" form:"name"` // param name
	Writable  string    `json:"writable" form:"writable"`           // param perm
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// timescaleDB hypertable definition

// TsDeviceLoad device load
type TsDeviceLoad struct {
	Timestamp   time.Time `gorm:"primaryKey" json:"timestamp"`
	Sn          string    `gorm:"primaryKey" json:"sn"`
	Name        string    `json:"name"`
	CpuLoad     int64     `json:"cpu_load"`
	MemUsage    int64     `json:"mem_usage"`
	MemPercent  float64   `json:"mem_percent"`
	DiskUsage   int64     `json:"disk_usage"`
	DiskPercent float64   `json:"disk_percent"`
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
