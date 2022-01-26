package models

import (
	"time"
)

// 网络模块相关模型

// NetRegion 管理区域
type NetRegion struct {
	ID        int64     `json:"id,string" form:"id"`  // 主键 ID
	Name      string    `json:"name" form:"name"`     // 设备名称
	Tags      string    `json:"tags" form:"tags"`     // 标签
	Remark    string    `json:"remark" form:"remark"` // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NetEdge 边缘节点表
type NetEdge struct {
	ID        int64     `json:"id,string"  form:"id" `             // 主键 ID
	RegionId  int64     `json:"region_id,string" form:"region_id"` // 区域ID
	Eid       string    `json:"eid" form:"eid"`                    // 节点标识
	Sn        string    `json:"sn" form:"sn"`                      // 节点设备SN
	Name      string    `json:"name" form:"name"`                  // 节点名称∂
	IpAddr    string    `json:"ip_addr"`                           // 内部IP
	Status    string    `json:"status" form:"status"`              // 状态 online/offline
	Tags      string    `json:"tags" form:"tags"`                  // 标签
	Remark    string    `json:"remark" form:"remark"`              // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NetEdgeEventLog struct {
	ID        int64     `json:"id,string"`  // 主键 ID
	Eid       string    `json:"eid"`        // 节点标识
	SessionId string    `json:"session_id"` // 会话ID
	LogType   string    `json:"log_type"`   // 日志类型  schedule | supervise
	Action    string    `json:"action"`     // 执行动作, 比如 hostping, speedtest, traceroute
	Level     string    `json:"level"`      // 动作级别  info｜error
	Message   string    `json:"message"`    // 消息内容
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// NetDevice Device 数据模型
type NetDevice struct {
	ID              int64     `json:"id,string" form:"id"`                      // 主键ID
	RegionId        int64     `json:"region_id,string" form:"region_id"`        // 区域ID
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
	CwmpLastInform  time.Time `json:"cwmp_last_inform" `                    // CWMP 最后检测时间
	SnmpPort        int       `json:"snmp_port" form:"snmp_port"`           // snmp 端口
	SnmpCommunity   string    `json:"snmp_community" form:"snmp_community"` // snmp 密码
	SnmpStatus      string    `json:"snmp_status"`                          // snmp 状态
	CwmpStatus      string    `json:"cwmp_status"`                          // cwmp 状态
	Tags            string    `json:"tags" form:"tags"`                     // 标签
	Uptime          int64     `json:"uptime" form:"uptime"`                 // UpTime
	MemoryTotal     int64     `json:"memory_total" form:"memory_total"`     // 内存总量
	MemoryFree      int64     `json:"memory_free" form:"memory_free"`       // 内存可用
	CPUUsage        int64     `json:"cpu_usage" form:"cpu_usage"`           // CPE 百分比
	Remark          string    `json:"remark" form:"remark"`                 // 备注
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

// NetDeviceEventLog 设备事件日志
type NetDeviceEventLog struct {
	ID        int64     `json:"id,string"`
	SessionId string    `json:"session_id"` // 事件会话ID
	Sn        string    `json:"sn"`         // 设备序列号
	LogType   string    `json:"log_type"`   // 日志类型  schedule | supervise
	Level     string    `json:"level"`
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// NetLocalScript 本地脚本
type NetLocalScript struct {
	ID        string    `gorm:"primaryKey" json:"id" form:"id"` // 主键 ID
	Name      string    `json:"name" form:"name"`
	Lang      string    `json:"lang" form:"lang"`
	Stype     string    `json:"stype" form:"stype"` // normal | supervise
	Key       string    `gorm:"unique" json:"key" form:"key"`
	Content   string    `json:"content" form:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NetIpaddress struct {
	Ipaddress string    `gorm:"primaryKey" json:"ipaddress"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Isp       string    `json:"isp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
