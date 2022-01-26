package app

import (
	"github.com/ca17/teamsacs/models"
)

func initSettings() {
	GormDB.Where("1 = 1").Delete(&models.SysConfig{})
	GormDB.Create(&models.SysConfig{ID: 0, Sort: 1, Type: "system", Name: "SystemTitle", Value: "MetasLink", Remark: "系统标题"})
	GormDB.Create(&models.SysConfig{ID: 0, Sort: 2, Type: "system", Name: "SystemTimeZone", Value: "Asia/Shanghai", Remark: "系统时区"})
	GormDB.Create(&models.SysConfig{ID: 0, Sort: 3, Type: "system", Name: "SystemLoginRemark", Value: "Recommended:   Chrome/Edge", Remark: "登录页面描述"})
}

// 初始化时序表
func initTimescaleTable() {
	GormDB.Exec("SELECT create_hypertable('tx_device_log', 'timestamp');")
	GormDB.Exec("SELECT create_hypertable('tx_device_flows', 'timestamp');")
	GormDB.Exec("SELECT create_hypertable('tx_device_hour_rate', 'timestamp');")
}
