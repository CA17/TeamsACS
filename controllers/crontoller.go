package controllers

import (
	"github.com/ca17/teamsacs/controllers/cpe"
	"github.com/ca17/teamsacs/controllers/cwmpconfig"
	"github.com/ca17/teamsacs/controllers/cwmppreset"
	"github.com/ca17/teamsacs/controllers/dashboard"
	"github.com/ca17/teamsacs/controllers/factoryreset"
	"github.com/ca17/teamsacs/controllers/files"
	"github.com/ca17/teamsacs/controllers/firmwareconfig"
	"github.com/ca17/teamsacs/controllers/index"
	"github.com/ca17/teamsacs/controllers/logging"
	"github.com/ca17/teamsacs/controllers/metrics"
	"github.com/ca17/teamsacs/controllers/node"
	"github.com/ca17/teamsacs/controllers/opr"
	"github.com/ca17/teamsacs/controllers/settings"
	"github.com/ca17/teamsacs/controllers/supervise"
	"github.com/ca17/teamsacs/controllers/translate"
)

// Init web 控制器初始化
func Init() {
	index.InitRouter()
	opr.InitRouter()
	settings.InitRouter()
	dashboard.InitRouter()
	cpe.InitRouter()
	logging.InitRouter()
	node.InitRouter()
	factoryreset.InitRouter()
	firmwareconfig.InitRouter()
	cwmpconfig.InitRouter()
	supervise.InitRouter()
	cwmppreset.InitRouter()
	metrics.InitRouter()
	translate.InitRouter()
	files.InitRouter()
}
