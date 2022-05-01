package app

import (
	"github.com/ca17/teamsacs/models"
)

/**

CREATE USER teamsacs WITH PASSWORD 'teamsacs'

CREATE DATABASE teamsacs_osc OWNER postgres;
CREATE DATABASE teamsacs OWNER teamsacs;

GRANT ALL PRIVILEGES ON DATABASE teamsacs TO teamsacs;

*/

func initSettings() {
	GormDB.Where("1 = 1").Delete(&models.SysConfig{})
}

// 初始化时序表
func initTimescaleTable() {
	GormDB.Exec("SELECT create_hypertable('tx_device_load', 'timestamp');")
}
