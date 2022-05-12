package models

import (
	"time"

	"github.com/ca17/teamsacs/common/timeutil"
)

// RADIUS 相关模型

// RadiusProfile RADIUS 策略
type RadiusProfile struct {
	ID         int64     `json:"id,string" form:"id"`             // 主键 ID
	PnodeId    int64     `json:"pnode_id,string" form:"pnode_id"` // PN节点ID
	Name       string    `json:"name" form:"name"`                // 策略名称
	Status     string    `json:"status" form:"status"`            // 策略状态 0：禁用 1：正常
	AddrPool   string    `json:"addr_pool" form:"addr_pool"`      // 策略地址池
	ActiveNum  int       `json:"active_num" form:"active_num"`    // 并发数
	UpRate     int       `json:"up_rate" form:"up_rate"`          // 上行速率
	DownRate   int       `json:"down_rate" form:"down_rate"`      // 下行速率
	Level      string    `json:"level" form:"level"`              // 资费级别
	Tags       string    `json:"tags" form:"tags"`                // 标签
	Network    string    `json:"network" form:"network"`          // 网络IP段
	InternalIp string    `json:"internal_ip" form:"internal_ip"`  // 内部 ip
	AgIp       string    `json:"ag_ip" form:"ag_ip"`              // agip
	Dns        string    `json:"dns" form:"dns"`                  // DNS
	Remark     string    `json:"remark" form:"remark"`            // 备注
	DataVer    string    `json:"data_ver" form:"data_ver"`
	CreatedAt  time.Time `json:"created_at" form:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" form:"updated_at"`
}

// RadiusUser RADIUS 认证帐号
type RadiusUser struct {
	ID         int64 `json:"id,string" form:"id"`                   // 主键 ID
	CustomerId int64 `json:"customer_id,string" form:"customer_id"` // 客户ID
	ProfileId  int64 `json:"profile_id,string" form:"profile_id"`   // RADIUS 策略ID
	// PnIds      string             `json:"pn_ids" form:"pn_ids"`                  // PN 节点关联 IDS
	CpeIds     string             `json:"cpe_ids" form:"cpe_ids"`                // CPE 关联 IDS
	UserType   string             `json:"user_type" form:"user_type"`            // 帐号类型 usr cpe
	Realname   string             `json:"realname" form:"realname"`              // 联系人姓名
	Mobile     string             `json:"mobile" form:"mobile"`                  // 联系人电话
	Username   string             `json:"username" gorm:"index" form:"username"` // 账号名
	Password   string             `json:"password" form:"password"`              // 密码
	AddrPool   string             `json:"addr_pool" form:"addr_pool"`            // 策略地址池
	ActiveNum  int                `json:"active_num" form:"active_num"`          // 并发数
	UpRate     int                `json:"up_rate" form:"up_rate"`                // 上行速率
	DownRate   int                `json:"down_rate" form:"down_rate"`            // 下行速率
	IpAddr     string             `json:"ip_addr" form:"ip_addr"`                // 静态IP
	ExpireTime timeutil.LocalTime `json:"expire_time" form:"expire_time"`        // 过期时间
	Status     string             `json:"status" form:"status"`                  // 状态 enabled | disabled
	Tags       string             `json:"tags" form:"tags"`                      // 标签
	Remark     string             `json:"remark" form:"remark"`                  // 备注
	DataVer    string             `json:"data_ver" form:"data_ver"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type RadiusAuthlog struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username,omitempty"`
	NasAddr   string    `json:"nas_addr,omitempty"`
	Cast      int       `json:"cast,omitempty"`
	Result    string    `json:"result,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// RadiusOnline
// Radius RadiusOnline Recode
type RadiusOnline struct {
	ID                int64     `json:"id,string"` // 主键 ID
	Username          string    `json:"username"`
	NasId             string    `json:"nas_id"`
	NasAddr           string    `json:"nas_addr"`
	NasPaddr          string    `json:"nas_paddr"`
	SessionTimeout    int       `json:"session_timeout"`
	FramedIpaddr      string    `json:"framed_ipaddr"`
	FramedNetmask     string    `json:"framed_netmask"`
	MacAddr           string    `json:"mac_addr"`
	NasPort           int64     `json:"nas_port,string"`
	NasClass          string    `json:"nas_class"`
	NasPortId         string    `json:"nas_port_id"`
	NasPortType       int       `json:"nas_port_type"`
	ServiceType       int       `json:"service_type"`
	AcctSessionId     string    `json:"acct_session_id"`
	AcctSessionTime   int       `json:"acct_session_time"`
	AcctInputTotal    int64     `json:"acct_input_total,string"`
	AcctOutputTotal   int64     `json:"acct_output_total,string"`
	AcctInputPackets  int       `json:"acct_input_packets"`
	AcctOutputPackets int       `json:"acct_output_packets"`
	AcctStartTime     time.Time `json:"acct_start_time"`
	LastUpdate        time.Time `json:"last_update"`
}

// RadiusAccounting
// Radius Accounting Recode
type RadiusAccounting struct {
	ID                int64     `json:"id,string"` // 主键 ID
	Username          string    `json:"username"`
	NasId             string    `json:"nas_id"`
	NasAddr           string    `json:"nas_addr"`
	NasPaddr          string    `json:"nas_paddr"`
	SessionTimeout    int       `json:"session_timeout"`
	FramedIpaddr      string    `json:"framed_ipaddr"`
	FramedNetmask     string    `json:"framed_netmask"`
	MacAddr           string    `json:"mac_addr"`
	NasPort           int64     `json:"nas_port,string"`
	NasClass          string    `json:"nas_class"`
	NasPortId         string    `json:"nas_port_id"`
	NasPortType       int       `json:"nas_port_type"`
	ServiceType       int       `json:"service_type"`
	AcctSessionId     string    `json:"acct_session_id"`
	AcctSessionTime   int       `json:"acct_session_time"`
	AcctInputTotal    int64     `json:"acct_input_total,string"`
	AcctOutputTotal   int64     `json:"acct_output_total,string"`
	AcctInputPackets  int       `json:"acct_input_packets"`
	AcctOutputPackets int       `json:"acct_output_packets"`
	AcctStartTime     time.Time `json:"acct_start_time"`
	LastUpdate        time.Time `json:"last_update"`
	AcctStopTime      time.Time `json:"acct_stop_time"`
}
