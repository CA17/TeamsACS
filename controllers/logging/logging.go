package logging

import (
	"net/http"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/webserver"
	"github.com/labstack/echo/v4"
)

func InitRouter() {

	webserver.GET("/admin/logging", func(c echo.Context) error {
		return c.Render(http.StatusOK, "logging", map[string]interface{}{})
	})

	webserver.GET("/admin/logging/query", func(c echo.Context) error {
		var count, start int
		web.NewParamReader(c).
			ReadInt(&start, "start", 0).
			ReadInt(&count, "count", 40)
		var data []models.SysOprLog
		prequery := web.NewPreQuery(c).
			DefaultOrderBy("opt_time desc").
			DateRange2("starttime", "endtime", "opt_time", time.Now().Add(-time.Hour*8), time.Now()).
			KeyFields("opr_name", "opt_action", "opr_ip", "opt_desc")

		var total int64
		common.Must(prequery.Query(app.GDB().Model(&models.SysOprLog{})).Count(&total).Error)

		query := prequery.Query(app.GDB().Debug().Model(&models.SysOprLog{})).Offset(start).Limit(count)
		if query.Find(&data).Error != nil {
			return c.JSON(http.StatusOK, common.EmptyList)
		}
		return c.JSON(http.StatusOK, &web.PageResult{TotalCount: total, Pos: int64(start), Data: data})
	})

}
