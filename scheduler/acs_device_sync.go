package scheduler

import (
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)



func SyncAcsDeviceInfo(manager *models.ModelManager) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	devinfos, err := manager.GetGenieacsManager().QueryMikrotikDeviceInfo()
	if err != nil {
		log.Error("SyncAcsDeviceInfo error, query deviceInfo error", err)
		return
	}

	log.Infof("fetch device num %d", len(devinfos))

	manager.GetGenieacsManager().SyncDeviceInfo(devinfos)

}
