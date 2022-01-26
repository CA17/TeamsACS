package app

import (
	"time"

	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
)

// Sched 计划任务管理
var Sched *cron.Cron

// Init 初始化任务计划
func setupJobs() {
	loc, _ := time.LoadLocation(Config.System.Location)
	Sched = cron.New(cron.WithLocation(loc))

	// 清理日志任务
	Sched.AddFunc("@every 3600s", func() {

	})

	// 每天凌晨的第30分钟执行小时性能统计
	Sched.AddFunc("0 30 0 * * *", func() {

	})

	Sched.Start()
}
