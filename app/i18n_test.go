package app

import (
	"testing"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/config"
)

func TestApplication_Translate(t *testing.T) {
	InitGlobalApplication(config.LoadConfig("../teamsacs.yml"))
	app.Translate(ZhCN, "global", "Create", "Create")
	app.Translate(ZhCN, "global", "Remove", "Remove")
	app.Translate(ZhCN, "global", "Edit", "编辑")
	app.Translate(ZhCN, "global", "Node", "节点")
	app.Translate(ZhCN, "global", "Exit", "退出")
	rets := app.LoadTranslateDict(ZhCN)
	t.Log(common.ToJson(rets))
	t.Log(common.ToJson(app.QueryTranslateTable(ZhCN, "", "")))
	app.RenderTranslateFiles()
}
