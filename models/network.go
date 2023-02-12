package models

import "time"

// 网络模块相关模型

// NetNode network node
type NetNode struct {
	ID        int64     `json:"id,string" form:"id"`
	Name      string    `json:"name" form:"name"`
	Remark    string    `json:"remark" form:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NetCpe Cpe 数据模型
type NetCpe struct {
	ID              int64     `json:"id,string" form:"id"`                      // primaryKey ID
	NodeId          int64     `json:"node_id,string" form:"node_id"`            // Node ID
	SystemName      string    `json:"system_name" form:"system_name"`           // Device system name
	Sn              string    `gorm:"uniqueIndex" json:"sn" form:"sn"`          // devise serial number
	Name            string    `json:"name" form:"name"`                         // device name
	ArchName        string    `json:"arch_name" form:"arch_name"`               // device architecture
	SoftwareVersion string    `json:"software_version" form:"software_version"` // Device firmware version
	HardwareVersion string    `json:"hardware_version" form:"hardware_version"` // device version
	Model           string    `json:"model" form:"model"`                       // Device model
	Oui             string    `json:"oui" form:"oui"`                           // Device OUI
	Manufacturer    string    `json:"manufacturer" form:"manufacturer"`         // Device manufactory
	ProductClass    string    `json:"product_class" form:"product_class"`       // Device type
	Status          string    `gorm:"index" json:"status" form:"status"`        // Device status enabled | disabled
	TaskTags        string    `gorm:"index" json:"task_tags" form:"task_tags"`  // task tags
	Uptime          int64     `json:"uptime" form:"uptime"`                     // UpTime
	MemoryTotal     int64     `json:"memory_total" form:"memory_total"`         // total memory
	MemoryFree      int64     `json:"memory_free" form:"memory_free"`           // memory available
	CPUUsage        int64     `json:"cpu_usage" form:"cpu_usage"`               // CPE Percentage
	CwmpStatus      string    `gorm:"index"  json:"cwmp_status"`                // cwmp status
	CwmpUrl         string    `json:"cwmp_url"`
	FactoryresetId  string    `json:"factoryreset_id" form:"factoryreset_id"`
	CwmpLastInform  time.Time `json:"cwmp_last_inform" `    // CWMP last inform time
	Remark          string    `json:"remark" form:"remark"` // Remark
	CreatedAt       time.Time `json:"created_at" `
	UpdatedAt       time.Time `json:"updated_at"`
}

// NetCpeParam CPE 参数
type NetCpeParam struct {
	ID        string    `gorm:"primaryKey" json:"string"` // primaryKey ID
	Sn        string    `gorm:"index" json:"sn"`          // devise serial number
	Tag       string    `gorm:"index" json:"tag" `
	Name      string    `gorm:"index" json:"name" `
	Value     string    `json:"value" `
	Remark    string    `json:"remark"`
	Writable  string    `json:"writable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NetCpeTaskQue struct {
	ID     int64  `json:"id,string"` // primaryKey ID
	Sn     string `json:"sn"`        // devise serial number
	TaskId string `json:"task_id"`
}
