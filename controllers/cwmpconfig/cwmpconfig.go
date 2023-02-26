package cwmpconfig

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/webserver"
	"github.com/labstack/echo/v4"
)

func initTemplateRouter() {

	webserver.GET("/admin/cwmp/config", func(c echo.Context) error {
		return c.Render(http.StatusOK, "cwmp_config", map[string]interface{}{
			"oprlevel": webserver.GetCurrUserlevel(c),
		})
	})

	webserver.GET("/admin/cwmp/config/options", func(c echo.Context) error {
		var data []models.CwmpConfig
		common.Must(app.GDB().Find(&data).Error)
		var opts = make([]web.JsonOptions, 0)
		for _, d := range data {
			opts = append(opts, web.JsonOptions{
				Id:    d.ID,
				Value: d.Name,
			})
		}
		return c.JSON(http.StatusOK, opts)
	})

	webserver.GET("/admin/cwmp/config/get", func(c echo.Context) error {
		var id string
		web.NewParamReader(c).ReadRequiedString(&id, "id")
		var data models.CwmpConfig
		err := app.GDB().Where("id=?", id).First(&data).Error
		if err != nil {
			return c.JSON(http.StatusOK, common.EmptyData)
		}
		return c.JSON(http.StatusOK, data)
	})

	webserver.GET("/admin/cwmp/config/query", func(c echo.Context) error {
		prequery := web.NewPreQuery(c).
			DefaultOrderBy("updated_at desc").
			KeyFields("oid", "name", "software_version",
				"product_class", "oui", "task_tags")

		result, err := web.QueryPageResult[models.CwmpConfig](c, app.GDB(), prequery)
		if err != nil {
			return c.JSON(http.StatusOK, common.EmptyList)
		}
		return c.JSON(http.StatusOK, result)
	})

	webserver.POST("/admin/cwmp/config/add", func(c echo.Context) error {
		form := new(models.CwmpConfig)
		common.Must(c.Bind(form))
		form.ID, _ = common.UUIDBase32()
		common.MustNotEmpty("名称", form.Name)
		common.Must(app.GDB().Create(form).Error)
		webserver.PubOpLog(c, fmt.Sprintf("Create CWMP config information：%v", form))
		return c.JSON(http.StatusOK, web.RestSucc("success"))
	})

	webserver.POST("/admin/cwmp/config/update", func(c echo.Context) error {
		form := new(models.CwmpConfig)
		common.Must(c.Bind(form))
		common.MustNotEmpty("名称", form.Name)
		common.Must(app.GDB().Save(form).Error)
		webserver.PubOpLog(c, fmt.Sprintf("Update CWMP config information：%v", form))
		return c.JSON(http.StatusOK, web.RestSucc("success"))
	})

	webserver.GET("/admin/cwmp/config/delete", func(c echo.Context) error {
		ids := c.QueryParam("ids")
		common.Must(app.GDB().Delete(models.CwmpConfig{}, strings.Split(ids, ",")).Error)
		webserver.PubOpLog(c, fmt.Sprintf("Delete CWMP config information：%s", ids))
		return c.JSON(http.StatusOK, web.RestSucc("success"))
	})

	webserver.POST("/admin/cwmp/config/import", func(c echo.Context) error {
		datas, err := webserver.ImportData(c, "cwmpconfig")
		common.Must(err)
		common.Must(app.GDB().Model(models.CwmpConfig{}).Create(datas).Error)
		return c.JSON(http.StatusOK, web.RestSucc("Success"))
	})

	webserver.GET("/admin/cwmp/config/export", func(c echo.Context) error {
		var data []models.CwmpConfig
		common.Must(app.GDB().Find(&data).Error)
		datas := make([]map[string]interface{}, 0)
		for _, d := range data {
			mitem, err := common.StructToMap(&d)
			if err == nil {
				datas = append(datas, mitem)
			}
		}
		return webserver.ExportData(c, datas, "cwmpconfig")
	})

}
