package scheduler

import (
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

var picmap = map[string]string{
	"hEX S":             "hEX-S.png",
	"hAP ac²":           "hAP-ac2.png",
	"Audience LTE6 kit": "Audience.png",
	"Audience":          "Audience.png",
	"RB4011iGS+":        "RB4011iGS+.png",
}

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

	ctime := time.Now()
	for _, dev := range devinfos {
		sn := dev.SerialNumber
		if sn == "" {
			continue
		}

		picture, ok := picmap[dev.ProductClass]
		if !ok {
			picture = "cpe.png"
		}
		existVpe := manager.GetVpeManager().ExistVpe(sn)
		_valmap := map[string]interface{}{
			"identifier":    dev.X_MIKROTIK_SystemIdentity,
			"manufacturer":  dev.Manufacturer,
			"device_id":     dev.DeviceId,
			"product_class": dev.ProductClass,
			"oui":           dev.ManufacturerOUI,
			"model":         dev.ModelName,
			"uptime":        dev.UpTime,
			"cpuuse":        dev.CPUUsage,
			"memuse":        dev.MemoryUsage,
			"version":       dev.HardwareVersion,
			"picture":       picture,
			"update_time":   ctime,
			"last_inform":   ctime,
		}
		valmap := make(map[string]interface{})
		for k, v := range _valmap {
			if v != "" {
				valmap[k] = v
			}
		}
		if existVpe {
			err = manager.GetVpeManager().UpdateVpeBySn(sn, valmap)
			if err != nil {
				log.Errorf("SyncAcsDeviceInfo update vpe:sn=%s error %s", sn, err.Error())
				continue
			}
			log.Infof("SyncAcsDeviceInfo update vpe:sn=%s", sn)
			continue
		}

		existCpe := manager.GetCpeManager().ExistCpe(sn)

		if existCpe {
			delete(valmap, "identifier")
			err = manager.GetCpeManager().UpdateCpeBySn(sn, valmap)
			if err != nil {
				log.Errorf("SyncAcsDeviceInfo update cpe:sn=%s error %s", sn, err.Error())
				continue
			}
			log.Infof("SyncAcsDeviceInfo update cpe:sn=%s", sn)
			continue
		} else {
			cpe := models.Cpe{}
			cpe.Set("_id", common.UUID())
			cpe.Set("name", common.EmptyToNA(dev.X_MIKROTIK_SystemIdentity))
			cpe.Set("sn", sn)
			cpe.Set("device_id", common.EmptyToNA(dev.DeviceId))
			cpe.Set("product_class", common.EmptyToNA(dev.ProductClass))
			cpe.Set("manufacturer", common.EmptyToNA(dev.Manufacturer))
			cpe.Set("version", common.EmptyToNA(dev.HardwareVersion))
			cpe.Set("oui", common.EmptyToNA(dev.ManufacturerOUI))
			cpe.Set("model", common.EmptyToNA(dev.ModelName))
			cpe.Set("cpuuse", int(dev.CPUUsage))
			cpe.Set("memuse", int(dev.MemoryUsage))
			cpe.Set("rd_ipaddr", "")
			cpe.Set("remark", "tr069 auto join")
			cpe.Set("status",  common.ENABLED)
			cpe.Set("picture",  picture)
			cpe.Set("create_time",  ctime.Format("2006-01-02 15:04:05"))


			err = manager.GetCpeManager().AddCpeDataMap(cpe)
			if err != nil {
				log.Errorf("SyncAcsDeviceInfo add cpe:sn=%s  error, %s", sn, err)
				continue
			}
		}

	}

}
