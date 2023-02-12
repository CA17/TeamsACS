package dashboard

import (
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/echarts"
	"github.com/ca17/teamsacs/common/zaplog"
	"github.com/ca17/teamsacs/webserver"
	"github.com/labstack/echo/v4"
)

func initSystemMetricsRouter() {

	webserver.GET("/admin/metrics/cpuuse/line", func(c echo.Context) error {
		var items []echarts.MetricLineItem

		points, err := zaplog.TSDB().Select("teamsacs_cpuuse", nil,
			time.Now().Add(-24*time.Hour).Unix(), time.Now().Unix())
		if err != nil {
			return c.JSON(200, common.EmptyList)
		}
		for i, p := range points {
			items = append(items, echarts.MetricLineItem{
				Id:    i + 1,
				Time:  time.Unix(p.Timestamp, 0).Format("2006-01-02 15:04"),
				Value: p.Value,
			})
		}

		result := echarts.AvgMetricLine(items)
		tsdata := echarts.NewTimeValues()
		for _, item := range result {
			timestamp, err := time.Parse("2006-01-02 15:04", item.Time)
			if err != nil {
				continue
			}
			tsdata.AddData(timestamp.Unix()*1000, item.Value)
		}
		so := echarts.NewSeriesObject("line")
		so.SetAttr("showSymbol", false)
		so.SetAttr("smooth", true)
		so.SetAttr("areaStyle", echarts.Dict{})
		so.SetAttr("data", tsdata)

		return c.JSON(200, echarts.Series(so))
	})

	webserver.GET("/admin/metrics/memuse/line", func(c echo.Context) error {

		var items []echarts.MetricLineItem

		points, err := zaplog.TSDB().Select("teamsacs_memuse", nil,
			time.Now().Add(-24*time.Hour).Unix(), time.Now().Unix())
		if err != nil {
			return c.JSON(200, common.EmptyList)
		}
		for i, p := range points {
			items = append(items, echarts.MetricLineItem{
				Id:    i + 1,
				Time:  time.Unix(p.Timestamp, 0).Format("2006-01-02 15:04"),
				Value: p.Value,
			})
		}

		result := echarts.AvgMetricLine(items)
		tsdata := echarts.NewTimeValues()
		for _, item := range result {
			timestamp, err := time.Parse("2006-01-02 15:04", item.Time)
			if err != nil {
				continue
			}
			tsdata.AddData(timestamp.Unix()*1000, item.Value)
		}
		so := echarts.NewSeriesObject("line")
		so.SetAttr("showSymbol", false)
		so.SetAttr("smooth", true)
		so.SetAttr("areaStyle", echarts.Dict{})
		so.SetAttr("data", tsdata)
		return c.JSON(200, echarts.Series(so))
	})
}
