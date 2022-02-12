package app

import (
	"errors"
	"sync"
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/cwmp"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/timeutil"
	"github.com/ca17/teamsacs/models"
)

// CwmpEventTable CPE 缓存表
type CwmpEventTable struct {
	cpeTable map[string]*CwmpCpe
	cpeLock  sync.Mutex
}

type CwmpEventData struct {
	Session string       `json:"session"`
	Sn      string       `json:"sn"`
	Message cwmp.Message `json:"message"`
}

type CwmpCpe struct {
	Sn             string `json:"sn"`
	cwmpQueueMap   chan CwmpEventData
	cwmpHPQueueMap chan CwmpEventData
	LastInform     *cwmp.Inform `json:"latest_message"`
	LastUpdate     time.Time    `json:"last_update"`
	LastDataNotify time.Time    `json:"last_data_notify"`
	IsRegister     bool         `json:"is_register"`
}

func NewCwmpEventTable() *CwmpEventTable {
	et := &CwmpEventTable{
		cpeTable: make(map[string]*CwmpCpe),
	}
	return et
}

// GetCwmpCpe 查询或者创建一个缓存 CPE 数据
func (c *CwmpEventTable) GetCwmpCpe(key string) *CwmpCpe {
	if common.IsEmptyOrNA(key) {
		panic(errors.New("key is empty"))
	}
	c.cpeLock.Lock()
	defer c.cpeLock.Unlock()
	cpe, ok := c.cpeTable[key]
	if !ok {
		var count int64 = 0
		GormDB.Model(models.NetDevice{}).Where("sn=?", key).Count(&count)
		cpe = &CwmpCpe{
			Sn:             key,
			LastUpdate:     timeutil.EmptyTime,
			LastDataNotify: timeutil.EmptyTime,
			cwmpQueueMap:   make(chan CwmpEventData, 32),
			cwmpHPQueueMap: make(chan CwmpEventData, 8),
			LastInform:     nil,
			IsRegister:     count > 0,
		}
		c.cpeTable[key] = cpe
	}
	return cpe
}

func (c *CwmpEventTable) UpdateCwmpCpe(key string, msg *cwmp.Inform) {
	cpe := c.GetCwmpCpe(key)
	cpe.UpdateStatus(msg)
}

func (c *CwmpCpe) UpdateStatus(msg *cwmp.Inform) {
	c.LastInform = msg
	c.LastUpdate = time.Now()
}

// CheckRegister 检测并自动注册CPE信息
func (c *CwmpCpe) CheckRegister(ip string, msg *cwmp.Inform) {
	if !c.IsRegister {
		var ctime = time.Now()
		err := GormDB.Create(&models.NetDevice{
			ID:              common.UUIDint64(),
			DevType:         "general",
			Sn:              msg.Sn,
			Name:            "Device-" + msg.Sn,
			SystemName:      "",
			ArchName:        "",
			PublicIpaddr:    ip,
			TunnelIpaddr:    "N/A",
			Model:           msg.ProductClass,
			VendorCode:      "0",
			SoftwareVersion: "unknow",
			HardwareVersion: "unknow",
			Oui:             msg.OUI,
			Manufacturer:    msg.Manufacturer,
			ProductClass:    msg.ProductClass,
			CwmpUrl:         msg.GetParam("Device.ManagementServer.ConnectionRequestURL"),
			CwmpLastInform:  ctime,
			Tags:            "",
			Uptime:          0,
			MemoryTotal:     0,
			MemoryFree:      0,
			CPUUsage:        0,
			Remark:          "first register",
			CreatedAt:       ctime,
			UpdatedAt:       ctime,
		}).Error
		if err == nil {
			c.IsRegister = true
		}
	}
}

// NotifyDataUpdate 通知 CPE 更新数据（发布通知 event ）
func (c *CwmpCpe) NotifyDataUpdate(force bool) {
	var ctime = time.Now()
	updateFlag := ctime.Sub(c.LastDataNotify).Seconds() > 60
	if force {
		updateFlag = true
	}
	if updateFlag {
		Bus.Publish(EventCwmpInformUpdate, c.Sn, c.LastInform)
		c.LastDataNotify = time.Now()
	}
}

func (c *CwmpCpe) UpdateParamValues(values map[string]string) {
	Bus.Publish(EventCwmpParamsUpdate, c.Sn, values)
}

func (c *CwmpCpe) getQueue(hp bool) chan CwmpEventData {
	var que = c.cwmpQueueMap
	if hp {
		que = c.cwmpHPQueueMap
	}
	return que
}

func setMapValue(vmap map[string]interface{}, name string, value interface{}) {
	if name != "" && value != "" {
		vmap[name] = value
	}
}

func (c *CwmpCpe) RecvCwmpEventData(timeoutMsec int, hp bool) (data *CwmpEventData, err error) {
	select {
	case _data := <-c.getQueue(hp):
		return &_data, nil
	case <-time.After(time.Millisecond * time.Duration(timeoutMsec)):
		return nil, errors.New("read cwmp event channel timeout")
	}
}

func (c *CwmpCpe) SendCwmpEventData(data CwmpEventData, timeoutMsec int, hp bool) error {
	select {
	case c.getQueue(hp) <- data:
		return nil
	case <-time.After(time.Microsecond * time.Duration(timeoutMsec)):
		return errors.New("cwmp event channel full, write timeout")
	}
}

func (c *CwmpCpe) UpdateManagementAuthInfo(session string, timeout int, hp bool) error {
	return c.SendCwmpEventData(CwmpEventData{
		Session: session,
		Sn:      c.Sn,
		Message: &cwmp.SetParameterValues{
			ID:     session,
			Name:   "",
			NoMore: 0,
			Params: map[string]cwmp.ValueStruct{
				"Device.ManagementServer.ConnectionRequestUsername": {
					Type:  "xsd:string",
					Value: c.Sn,
				},
				"Device.ManagementServer.ConnectionRequestPassword": {
					Type:  "xsd:string",
					Value: GetCwmpSettingsStringValue("CpeConnectionRequestPassword"),
				},
			},
		},
	}, timeout, hp)
}

func (c *CwmpCpe) UpdateParameterNames(session string, timeout int, hp bool) error {
	return c.SendCwmpEventData(CwmpEventData{
		Session: session,
		Sn:      c.Sn,
		Message: &cwmp.GetParameterNames{
			ID:            session,
			Name:          "获取参数名列表",
			NoMore:        0,
			ParameterPath: "Device.",
			NextLevel:     "false",
		},
	}, timeout, hp)
}

func (c *CwmpCpe) ProcessParameterNamesResponse(msg *cwmp.GetParameterNamesResponse) {
	var items []models.CwmpParamNameList
	for _, param := range msg.Params {
		items = append(items, models.CwmpParamNameList{Sn: c.Sn, Name: param.Name, Writable: param.Writable})
	}
	log.ErrorIf(GormDB.Save(&items).Error)
}
