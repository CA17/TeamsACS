package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"

	"github.com/ca17/teamsacs/models"
)

func Start(manager *models.ModelManager) error {
	s := gocron.NewScheduler(time.UTC)
	// 日志分发任务
	_, _ = s.Every(300).Second().StartImmediately().Do(SyncAcsDeviceInfo, manager)
	// _, _ = s.Every(10).Minute().Do(CheckAcsScriptTask, manager, "10m")
	// _, _ = s.Every(30).Minute().Do(CheckAcsScriptTask, manager, "30m")
	// _, _ = s.Every(60).Minute().Do(CheckAcsScriptTask, manager, "60m")
	// _, _ = s.Every(2).Hour().Do(CheckAcsScriptTask, manager, "2h")
	// _, _ = s.Every(4).Hour().Do(CheckAcsScriptTask, manager, "4h")
	// _, _ = s.Every(8).Hour().Do(CheckAcsScriptTask, manager, "8h")
	// _, _ = s.Every(12).Hour().Do(CheckAcsScriptTask, manager, "12h")
	// _, _ = s.Every(24).Hour().Do(CheckAcsScriptTask, manager, "24")
	// _, _ = s.Every(60).Seconds().StartImmediately().Do(UpdateSysinfoTask, manager)
	<-s.Start()
	return nil
}
