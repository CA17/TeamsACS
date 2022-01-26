package models

import (
	"time"

	"github.com/ca17/teamsacs/common/cwmp"
)

type CwmpEventData struct {
	Session string       `json:"session"`
	Sn      string       `json:"sn"`
	Message cwmp.Message `json:"message"`
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

// database

// CwmpVendorConfig  Cwmp 设备配置文件
type CwmpVendorConfig struct {
	ID           string    `gorm:"primaryKey" json:"id" form:"id"`     // 主键 ID
	Name         string    `json:"name" form:"name"`                   // 名称
	Oui          string    `json:"oui" form:"oui"`                     // 设备OUI
	Manufacturer string    `json:"manufacturer" form:"manufacturer"`   // 设备制造商
	ProductClass string    `json:"product_class" form:"product_class"` // 设备类型
	Level        string    `json:"level" form:"level"`                 // 脚本级别  normal｜major
	Content      string    `json:"content" form:"content"`             // 内容
	Format       string    `json:"format" form:"format"`               // 格式 xml | alter | json
	Timeout      int64     `json:"timeout" form:"timeout"`             // 超时时间 秒
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CwmpVendorConfigHistory  Cwmp 设备配置文件历史版本
type CwmpVendorConfigHistory struct {
	ID        string    `gorm:"primaryKey" json:"id" form:"id"` // 主键 ID
	ConfigId  string    `json:"config_id"`                      // 设备配置 ID
	Name      string    `json:"name" form:"name"`               // 名称
	Level     string    `json:"level" form:"level"`             // 脚本级别  normal｜major
	Content   string    `json:"content" form:"content"`         // 内容
	Format    string    `json:"format" form:"format"`           // 格式 xml | alter | json
	Timeout   int64     `json:"timeout" form:"timeout"`         // 超时时间 秒
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CwmpDownloadTask Cwmp 下载任务
type CwmpDownloadTask struct {
	ID           int64     `gorm:"primaryKey" json:"id,string"`       // 主键 ID
	DeviceId     int64     `json:"device_id,string" form:"device_id"` // 设备 ID
	ConfigId     string    `json:"config_id"`                         // 设备配置 ID
	SessionId    string    `json:"session_id" `                       // 会话
	Name         string    `json:"name"`                              // 名称
	Level        string    `json:"level"`                             // 脚本级别  normal｜major
	Content      string    `json:"content"`                           // 内容
	Filename     string    `json:"filename"`                          // 文件名
	Status       string    `json:"status"`                            // 格式 initialize | success | error | timeout
	FaultCode    int       `json:"fault_code"`                        // FaultCode
	FaultString  string    `json:"fault_string"`                      // FaultString
	Timeout      int64     `json:"timeout" `                          // 执行时间超时 秒
	StartTime    time.Time `json:"start_time"`                        // 执行时间
	CompleteTime time.Time `json:"complete_time"`                     // 完成时间
	CreatedAt    time.Time `json:"created_at"`
}

// CwmpSetParamTask Cwmp 配置修改任务
type CwmpSetParamTask struct {
	ID          int64     `gorm:"primaryKey" json:"id,string"`       // 主键 ID
	DeviceId    int64     `json:"device_id,string" form:"device_id"` // 设备 ID
	ConfigId    int64     `json:"config_id,string"`                  // 设备配置 ID
	SessionId   string    `json:"session_id" `                       // 会话
	Name        string    `json:"name"`                              // 名称
	Content     string    `json:"content"`                           // 内容
	FaultCode   string    `json:"fault_code"`                        // FaultCode
	FaultString string    `json:"fault_string"`                      // FaultString
	Status      string    `json:"status"`                            // 格式 initialize | success | error | timeout
	Timeout     int64     `json:"timeout" `                          // 执行时间超时 秒
	StartTime   time.Time `json:"start_time"`                        // 执行时间
	RespTime    time.Time `json:"resp_time"`                         // 响应时间
	CreatedAt   time.Time `json:"created_at"`
}

type CwmpParamNameList struct {
	Sn        string    `gorm:"primaryKey" json:"sn" form:"sn"`     // 设备序列号
	Name      string    `gorm:"primaryKey" json:"name" form:"name"` // 参数名
	Writable  string    `json:"writable" form:"writable"`           // 参数名
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
