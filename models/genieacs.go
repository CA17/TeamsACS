/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/validutil"
	"github.com/ca17/teamsacs/models/mikrotik"
)

type GenieacsManager struct{ *ModelManager }

func (m *ModelManager) GetGenieacsManager() *GenieacsManager {
	store, _ := m.ManagerMap.Get("GenieacsManager")
	return store.(*GenieacsManager)
}

// query device info by cpesn
func (m *GenieacsManager) QueryMikrotikDeviceInfo() ([]mikrotik.DeviceInfo, error) {
	items, err := m.QueryMikrotikSourceData("")
	if err != nil {
		return nil, err
	}
	result := make([]mikrotik.DeviceInfo, 0)
	linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   c.(mikrotik.TMap)["_id"],
				Value: mikrotik.GetObject(c, "Device.DeviceInfo"),
			}
		}).
		ForEachT(func(i linq.KeyValue) {
			var info = new(mikrotik.DeviceInfo)
			info.ParseBson(i.Value.(mikrotik.TMap))
			info.DeviceId = i.Key.(string)
			result = append(result, *info)
		})
	return result, nil
}

// query DNSServer by cpe sn
func (m *GenieacsManager) QueryMikrotikDnsServer(sn string) (*mikrotik.DnsClientServer, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.DNS.Client.Server"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.DnsClientServer {
			var r = new(mikrotik.DnsClientServer)
			r.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.T)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					r.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					rItem := new(mikrotik.DnsClientServerItem)
					rItem.Key = fmt.Sprintf("Device.DNS.Client.Server.%s", it.Key)
					rItem.ParseBson(it.Value.(mikrotik.TMap))
					r.Items = append(r.Items, *rItem)
				}
			})
			return r
		}).
		First()

	if result == nil {
		return nil, nil
	}
	return result.(*mikrotik.DnsClientServer), nil
}

// query routers bu cpe sn
func (m *GenieacsManager) QueryMikrotikRouters(sn string) (*mikrotik.DeviceRouter, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.Routing.Router.1.IPv4Forwarding"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.DeviceRouter {
			var r = new(mikrotik.DeviceRouter)
			r.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.T)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					r.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					rItem := new(mikrotik.DeviceRouterItem)
					rItem.Key = fmt.Sprintf("Routing.Router.1.IPv4Forwarding.%s", it.Key)
					rItem.ParseBson(it.Value.(mikrotik.TMap))
					r.Items = append(r.Items, *rItem)
				}
			})
			return r
		}).
		First()

	if result == nil {
		return nil, nil
	}
	return result.(*mikrotik.DeviceRouter), nil
}

// query all ethernet by cpe sn
func (m *GenieacsManager) QueryMikrotikEthernetInterface(sn string) (*mikrotik.EthernetInterface, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.Ethernet.Interface"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.EthernetInterface {
			var gei = new(mikrotik.EthernetInterface)
			gei.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					gei.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					eifItem := new(mikrotik.EthernetInterfaceItem)
					eifItem.Key = fmt.Sprintf("Ethernet.Interface.%s", it.Key)
					eifItem.ParseBson(it.Value.(mikrotik.TMap))
					gei.Items = append(gei.Items, *eifItem)
				}
			})
			return gei
		}).
		First()

	if result == nil {
		return nil, nil
	}
	return result.(*mikrotik.EthernetInterface), nil
}

// query all ppp interface by cpe sn
func (m *GenieacsManager) QueryMikrotikPPPInterface(sn string) (*mikrotik.PPPInterface, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.PPP.Interface"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.PPPInterface {
			var gei = new(mikrotik.PPPInterface)
			gei.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					gei.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					eifItem := new(mikrotik.PPPInterfaceItem)
					eifItem.Key = fmt.Sprintf("PPP.Interface.%s", it.Key)
					eifItem.ParseBson(it.Value.(mikrotik.TMap))
					gei.Items = append(gei.Items, *eifItem)
				}
			})
			return gei
		}).
		First()

	if result == nil {
		return nil, nil
	}
	return result.(*mikrotik.PPPInterface), nil
}

// query all ethernet by cpe sn
func (m *GenieacsManager) QueryMikrotikIpInterface(sn string) (*mikrotik.IpInterface, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.IP.Interface"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.IpInterface {
			var ipif = new(mikrotik.IpInterface)
			ipif.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					ipif.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					st := mikrotik.GetSubValue(it.Value.(mikrotik.TMap), "Status", "_value", false).(string)
					if st == "Up" {
						eifItem := new(mikrotik.IpInterfaceItem)
						eifItem.Key = fmt.Sprintf("IP.Interface.%s", it.Key)
						eifItem.ParseBson(it.Value.(mikrotik.TMap))
						ipif.Items = append(ipif.Items, *eifItem)
					}
				}
			})
			return ipif
		}).
		First()

	if result == nil {
		return nil, nil
	}
	return result.(*mikrotik.IpInterface), nil
}

// query all cpe data
func (m *GenieacsManager) QueryMikrotikSourceData(sn string) ([]map[string]interface{}, error) {
	findOptions := options.Find()
	findOptions.SetLimit(100)
	coll := m.GetGenieAcsCollection(GenieacsDevices)
	var q = bson.M{}
	if sn != "" {
		q["Device.DeviceInfo.SerialNumber._value"] = sn
	}
	cur, err := coll.Find(context.TODO(), q, findOptions)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0)
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		} else {
			items = append(items, elem)
		}
	}
	return items, nil
}

type GenieacsTask struct {
	Id        string    `json:"id"`
	TId       string    `json:"_id"`
	FileName  string    `json:"fileName"`
	FileType  string    `json:"fileType"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
}

// GetAcsTaskDeviceIdList
// Get a list of device IDs for currently running tasks
func (m *GenieacsManager) GetAcsTaskDeviceIdList() ([]string, error) {
	result := make([]string, 0)
	resp := make([]GenieacsTask, 0)

	hresp, err := http.Get(common.UrlJoin(m.Config.Genieacs.NbiUrl, "/tasks"))
	if err != nil {
		return nil, err
	}

	defer hresp.Body.Close()
	body, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tasks resp error %s", err.Error())
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("read tasks json resp error %s", err.Error())
	}
	for _, r := range resp {
		result = append(result, r.Device)
	}
	return result, nil
}

func (m *GenieacsManager) GetAcsTaskDataList() ([]GenieacsTask, error) {
	resp := make([]GenieacsTask, 0)

	hresp, err := http.Get(common.UrlJoin(m.Config.Genieacs.NbiUrl, "/tasks"))
	if err != nil {
		return nil, err
	}

	defer hresp.Body.Close()
	body, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tasks resp error %s", err.Error())
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("read tasks json resp error %s", err.Error())
	}
	return resp, nil
}

func (m *GenieacsManager) RetryAcsTaskData(ids string) error {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	for _, tid := range strings.Split(ids, ",") {
		url := common.UrlJoin(m.Config.Genieacs.NbiUrl, "/tasks/"+tid+"/retry")
		req, err := http.NewRequest(http.MethodPost, url, nil)
		common.Must(err)
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("RetryAcsTaskData error %s", err.Error())
			continue
		}
		log.Infof("RetryAcsTaskData device:%s %d", tid, resp.StatusCode)
	}
	return nil
}

func (m *GenieacsManager) DeleteAcsTaskData(ids string) error {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	for _, tid := range strings.Split(ids, ",") {
		url := common.UrlJoin(m.Config.Genieacs.NbiUrl, "/tasks/"+tid+"")
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		common.Must(err)
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("DeleteAcsTaskData error %s", err.Error())
			continue
		}
		log.Infof("DeleteAcsTaskData device:%s %d", tid, resp.StatusCode)
	}
	return nil
}

func (m *GenieacsManager) GetAcsTaskData(taskid string) (*GenieacsTask, error) {
	resp := new(GenieacsTask)

	hresp, err := http.Get(common.UrlJoin(m.Config.Genieacs.NbiUrl, "/tasks/"+taskid))
	if err != nil {
		return nil, err
	}

	defer hresp.Body.Close()
	body, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tasks resp error %s", err.Error())
	}

	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, fmt.Errorf("read tasks json resp error %s", err.Error())
	}
	return resp, nil
}

var picmap = map[string]string{
	"hEX S":             "hEX-S.png",
	"hAP ac²":           "hAP-ac2.png",
	"Audience LTE6 kit": "Audience.png",
	"Audience":          "Audience.png",
	"RB4011iGS+":        "RB4011iGS+.png",
}

// Sync all device info to teamsacs colls
func (m GenieacsManager) SyncMikrotikDeviceInfo(devinfos []mikrotik.DeviceInfo) {
	ctime := time.Now()
	for _, dev := range devinfos {
		sn := dev.SerialNumber
		if sn == "" {
			continue
		}
		log.Infof("Process Device sn=%s", sn)

		picture, ok := picmap[dev.ProductClass]
		if !ok {
			picture = "cpe.png"
		}
		existVpe := m.GetVpeManager().ExistVpe(sn)

		onlineStatus := "off"
		if time.Now().Sub(dev.Timestamp.Time()).Minutes() < 5 {
			onlineStatus = "on"
		}

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
			"sversion":      dev.SoftwareVersion,
			"timestamp":     dev.Timestamp,
			"picture":       picture,
			"update_time":   ctime,
			"last_inform":   ctime,
			"online_status": onlineStatus,
		}
		valmap := make(map[string]interface{})
		for k, v := range _valmap {
			if v != "" {
				valmap[k] = v
			}
		}
		if existVpe {
			err := m.GetVpeManager().UpdateVpeBySn(sn, valmap)
			if err != nil {
				log.Errorf("SyncAcsDeviceInfo update vpe:sn=%s error %s", sn, err.Error())
				continue
			}
			log.Infof("SyncAcsDeviceInfo update vpe:sn=%s", sn)
			continue
		}

		existCpe := m.GetCpeManager().ExistCpe(sn)

		if existCpe {
			err := m.GetCpeManager().UpdateCpeBySn(sn, valmap)
			if err != nil {
				log.Errorf("SyncAcsDeviceInfo update cpe:sn=%s error %s", sn, err.Error())
				continue
			}
			log.Infof("SyncAcsDeviceInfo update cpe:sn=%s", sn)
			continue
		} else {
			cpe := Cpe{}
			cpe.Set("_id", common.UUID())
			cpe.Set("name", common.EmptyToNA(dev.X_MIKROTIK_SystemIdentity))
			cpe.Set("sn", sn)
			cpe.Set("device_id", common.EmptyToNA(dev.DeviceId))
			cpe.Set("product_class", common.EmptyToNA(dev.ProductClass))
			cpe.Set("manufacturer", common.EmptyToNA(dev.Manufacturer))
			cpe.Set("version", common.EmptyToNA(dev.HardwareVersion))
			cpe.Set("sversion", common.EmptyToNA(dev.SoftwareVersion))
			cpe.Set("oui", common.EmptyToNA(dev.ManufacturerOUI))
			cpe.Set("model", common.EmptyToNA(dev.ModelName))
			cpe.Set("cpuuse", int(dev.CPUUsage))
			cpe.Set("memuse", int(dev.MemoryUsage))
			cpe.Set("rd_ipaddr", "")
			cpe.Set("remark", "tr069 auto join")
			cpe.Set("status", common.ENABLED)
			cpe.Set("picture", picture)
			cpe.Set("online_status", onlineStatus)
			cpe.Set("create_time", ctime.Format("2006-01-02 15:04:05"))

			err := m.GetCpeManager().AddCpeDataMap(cpe)
			if err != nil {
				log.Errorf("SyncAcsDeviceInfo add cpe:sn=%s  error, %s", sn, err)
				continue
			}
		}

	}

}
