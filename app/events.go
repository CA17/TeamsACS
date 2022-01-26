package app

import (
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/ca17/teamsacs/common/cwmp"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

var (
	Bus = evbus.New()
)

const (
	EventSyslog           = "EventSyslog"
	EventDeviceEventLog   = "EventDeviceEventLog"
	EventCwmpInformUpdate = "EventCwmpInformUpdate"
	EventCwmpParamsUpdate = "EventCwmpParamsUpdate"
)

func PubSyslog(logdata interface{}) {
	Bus.Publish(EventSyslog, logdata)
}

func PubDeviceEventLog(sn, logtype, session, action, level, message string) {
	Bus.Publish(EventDeviceEventLog, sn, logtype, session, action, level, message)
}

// setupEventSubscribe 事件订阅
func setupEventSubscribe() {

	// Syslog 订阅
	_ = Bus.SubscribeAsync(EventSyslog, func(data interface{}) {
		log.ErrorIf(GormDB.Create(data).Error)
	}, false)

	Bus.SubscribeAsync(EventCwmpInformUpdate, func(sn string, msg *cwmp.Inform) {
		valmap := map[string]interface{}{}
		setMapValue(valmap, "manufacturer", msg.Manufacturer)
		setMapValue(valmap, "product_class", msg.ProductClass)
		setMapValue(valmap, "oui", msg.OUI)
		setMapValue(valmap, "cwmp_status", "online")
		setMapValue(valmap, "cwmp_last_inform", time.Now())
		setMapValue(valmap, "cwmp_url", msg.GetParam("Device.ManagementServer.ConnectionRequestURL"))
		setMapValue(valmap, "software_version", msg.GetParam("Device.DeviceInfo.SoftwareVersion"))
		setMapValue(valmap, "hardware_version", msg.GetParam("Device.DeviceInfo.HardwareVersion"))
		setMapValue(valmap, "model", msg.GetParam("Device.DeviceInfo.ModelName"))
		setMapValue(valmap, "uptime", msg.GetParam("Device.DeviceInfo.UpTime"))
		setMapValue(valmap, "cpu_usage", msg.GetParam("Device.DeviceInfo.ProcessStatus.CPUUsage"))
		setMapValue(valmap, "memory_total", msg.GetParam("Device.DeviceInfo.MemoryStatus.Free"))
		setMapValue(valmap, "memory_free", msg.GetParam("Device.DeviceInfo.MemoryStatus.Total"))
		// mikrotik
		setMapValue(valmap, "arch_name", msg.GetParam("Device.DeviceInfo.X_MIKROTIK_ArchName"))
		setMapValue(valmap, "system_name", msg.GetParam("Device.DeviceInfo.X_MIKROTIK_SystemIdentity"))

		if len(valmap) > 0 {
			GormDB.Model(&models.NetDevice{}).Where("sn=?", sn).Updates(valmap)
		}
	}, false)

	Bus.SubscribeAsync(EventCwmpParamsUpdate, func(sn string, params map[string]string) {
		var getParam = func(name string) string {
			v, ok := params[name]
			if ok {
				return v
			}
			return ""
		}
		valmap := map[string]interface{}{}
		setMapValue(valmap, "cwmp_last_inform", time.Now())
		setMapValue(valmap, "cwmp_status", "online")
		setMapValue(valmap, "cwmp_url", getParam("Device.ManagementServer.ConnectionRequestURL"))
		setMapValue(valmap, "software_version", getParam("Device.DeviceInfo.SoftwareVersion"))
		setMapValue(valmap, "hardware_version", getParam("Device.DeviceInfo.HardwareVersion"))
		setMapValue(valmap, "model", getParam("Device.DeviceInfo.ModelName"))
		setMapValue(valmap, "uptime", getParam("Device.DeviceInfo.UpTime"))
		setMapValue(valmap, "cpu_usage", getParam("Device.DeviceInfo.ProcessStatus.CPUUsage"))
		setMapValue(valmap, "memory_total", getParam("Device.DeviceInfo.MemoryStatus.Free"))
		setMapValue(valmap, "memory_free", getParam("Device.DeviceInfo.MemoryStatus.Total"))
		// mikrotik
		setMapValue(valmap, "arch_name", getParam("Device.DeviceInfo.X_MIKROTIK_ArchName"))
		setMapValue(valmap, "system_name", getParam("Device.DeviceInfo.X_MIKROTIK_SystemIdentity"))

		if len(valmap) > 0 {
			GormDB.Model(&models.NetDevice{}).Where("sn=?", sn).Updates(valmap)
		}

	}, false)

}
