package app

import (
	"github.com/ca17/teamsacs/models"
)

func initSettings() {
	GormDB.Where("1 = 1").Delete(&models.SysConfig{})
}

// 初始化时序表
func initTimescaleTable() {
	GormDB.Exec("SELECT create_hypertable('tx_device_load', 'timestamp');")
}
