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

// NetDevice Device 数据模型
type NetDevice struct {
	ID              int64     `json:"id,string" form:"id"`                      // 主键ID
	DevType         string    `json:"dev_type" form:"dev_type"`                 // 设备类型 server | switch | router
	Sn              string    `gorm:"uniqueIndex" json:"sn" form:"sn"`          // 设备序列号
	Name            string    `json:"name" form:"name"`                         // 设备名称
	SystemName      string    `json:"system_name" form:"system_name"`           // 设备系统名称
	ArchName        string    `json:"arch_name" form:"arch_name"`               // 设备架构
	PublicIpaddr    string    `json:"public_ipaddr" form:"public_ipaddr"`       // 设备公网IP
	TunnelIpaddr    string    `json:"tunnel_ipaddr" form:"tunnel_ipaddr"`       // 设备隧道IP
	Model           string    `json:"model" form:"model"`                       // 设备型号
	VendorCode      string    `json:"vendor_code" form:"vendor_code"`           // 设备厂商代码
	SoftwareVersion string    `json:"software_version" form:"software_version"` // 设备固件版本
	HardwareVersion string    `json:"hardware_version" form:"hardware_version"` // 设备版本
	Oui             string    `json:"oui" form:"oui"`                           // 设备OUI
	Manufacturer    string    `json:"manufacturer" form:"manufacturer"`         // 设备制造商
	ProductClass    string    `json:"product_class" form:"product_class"`       // 设备类型
	CwmpUrl         string    `json:"cwmp_url" form:"cwmp_url"`
	CwmpLastInform  time.Time `json:"cwmp_last_inform" `                // CWMP 最后检测时间
	Tags            string    `json:"tags" form:"tags"`                 // 标签
	Uptime          int64     `json:"uptime" form:"uptime"`             // UpTime
	MemoryTotal     int64     `json:"memory_total" form:"memory_total"` // 内存总量
	MemoryFree      int64     `json:"memory_free" form:"memory_free"`   // 内存可用
	CPUUsage        int64     `json:"cpu_usage" form:"cpu_usage"`       // CPE 百分比
	Remark          string    `json:"remark" form:"remark"`             // 备注
	CreatedAt       time.Time `json:"created_at" `
	UpdatedAt       time.Time `json:"updated_at"`
}

// NetDeviceParam 设备参数信息
type NetDeviceParam struct {
	ID        int64     `json:"id" form:"id"`       // 主键 ID
	Sn        string    `json:"sn" form:"sn"`       // 设备序列号
	Name      string    `json:"name" form:"name"`   // 参数名
	Value     string    `json:"value" form:"value"` // 参数值
	LastFetch time.Time `json:"last_fetch" `        // 最后采集时间
}

type CwmpParamFetchTask struct {
	ID           int64     `json:"id"`                               // 主键 ID
	Oui          string    `json:"oui" form:"oui"`                   // 设备OUI
	Manufacturer string    `json:"manufacturer" form:"manufacturer"` // 设备制造商
	ProductClass string    `json:"productClass" form:"productClass"` // 设备类型
	Name         string    `json:"name" form:"name"`                 // 参数名
	Interval     int64     `json:"interval" form:"interval"`         // 间隔时间
	LastFetch    time.Time `json:"last_fetch" `                      // 最后采集时间
}

// CwmpFile  Cwmp 设备配置文件或固件文件
type CwmpFile struct {
	ID           string    `gorm:"primaryKey" json:"id" form:"id"`     // 主键 ID
	Name         string    `json:"name" form:"name"`                   // 名称
	Oui          string    `json:"oui" form:"oui"`                     // 设备OUI
	Manufacturer string    `json:"manufacturer" form:"manufacturer"`   // 设备制造商
	ProductClass string    `json:"product_class" form:"product_class"` // 设备类型
	Version      string    `json:"version" form:"version"`             // 设备版本
	FileType     string    `json:"file_type" form:"file_type"`         // 文件保存路径
	FilePath     string    `json:"file_path" form:"file_path"`         // 文件保存路径
	Format       string    `json:"format" form:"format"`               // 格式 xml | alter | json
	Timeout      int64     `json:"timeout" form:"timeout"`             // 超时时间 秒
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CwmpDownloadTask Cwmp 下载任务
type CwmpDownloadTask struct {
	ID           int64     `gorm:"primaryKey" json:"id,string"` // 主键 ID
	Sn           string    `json:"sn"`                          // 设备序列号
	SessionId    string    `json:"session_id" `                 // 会话
	Name         string    `json:"name"`                        // 名称
	FileType     string    `json:"file_type"`                   // 文件保存路径
	FilePath     string    `json:"file_path"`                   // 文件保存路径
	Filename     string    `json:"filename"`                    // 文件名
	Status       string    `json:"status"`                      // 格式 initialize | success | error | timeout
	FaultCode    int       `json:"fault_code"`                  // FaultCode
	FaultString  string    `json:"fault_string"`                // FaultString
	Timeout      int64     `json:"timeout" `                    // 执行时间超时 秒
	StartTime    time.Time `json:"start_time"`                  // 执行时间
	CompleteTime time.Time `json:"complete_time"`               // 完成时间
	CreatedAt    time.Time `json:"created_at"`
}

// CwmpMessageTask Cwmp 下载任务
type CwmpMessageTask struct {
	ID        int64     `gorm:"primaryKey" json:"id,string"` // 消息 ID
	Sn        string    `json:"sn"`                          // 设备序列号
	Name      string    `json:"name"`                        // 消息名
	Message   string    `json:"message"`                     // 消息内容
	CreatedAt time.Time `json:"created_at"`
}

type CwmpParamNameList struct {
	Sn        string    `gorm:"primaryKey" json:"sn" form:"sn"`     // 设备序列号
	Name      string    `gorm:"primaryKey" json:"name" form:"name"` // 参数名
	Writable  string    `json:"writable" form:"writable"`           // 参数名
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// timescaleDB 超表定义

// TsDeviceLoad 设备负载
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
