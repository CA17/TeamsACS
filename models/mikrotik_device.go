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
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/mikrotik_api"
	"github.com/ca17/teamsacs/models/elastic"
)

// Mikrotik Api method collection

type MikrotikDeviceManager struct{ *ModelManager }

func (m *ModelManager) GetMikrotikDeviceManager() *MikrotikDeviceManager {
	store, _ := m.ManagerMap.Get("MikrotikDeviceManager")
	return store.(*MikrotikDeviceManager)
}

func _ifLtZeroInt64(s, defval int64) int64 {
	if s < 0 {
		return defval
	}
	return s
}

// GetMikrotikApiBySN Get API according to SN
func (m *MikrotikDeviceManager) GetMikrotikApiBySN(sn, devtype string) (*mikrotik_api.MikrotikApi, error) {
	if sn == "" {
		return nil, fmt.Errorf("sn param  empty")
	}
	switch devtype {
	case "cpe":
		cpe, err := m.GetCpeManager().GetCpeBySn(sn)
		if err != nil {
			return nil, err
		}
		return m.GetMikrotikApi(devtype, *cpe)
	case "vpe":
		vpe, err := m.GetVpeManager().GetVpeBySn(sn)
		if err != nil {
			return nil, err
		}
		return m.GetMikrotikApi(devtype, *vpe)
	}
	return nil, fmt.Errorf("not support")

}

// GetMikrotikApi Get an API long connection, support automatic reconnection
// parammap: {api_addr:"xxx", api_user:"xxx", api_pwd:"xxx"}
func (m *MikrotikDeviceManager) GetMikrotikApi(devtype string, devmap map[string]interface{}) (*mikrotik_api.MikrotikApi, error) {

	apiAddr := maputils.GetStringValue(devmap, "api_addr", "")
	if common.InSlice(apiAddr, []string{"", "N/A"}) {
		switch devtype {
		case "cpe":
			apiAddr = maputils.GetStringValue(devmap, "rd_ipaddr", "") + ":8728"
		case "vpe":
			apiAddr = maputils.GetStringValue(devmap, "ipaddr", "") + ":8728"
		}
	}

	apiUser := maputils.GetStringValue(devmap, "api_user", "")
	apiPwd := ""
	if common.InSlice(apiUser, []string{"", "N/A"}) {
		apiUser = os.Getenv("TEAMSACS_MIKROTIK_APIUSER")
		apiPwd = os.Getenv("TEAMSACS_MIKROTIK_APIPWD")
	} else {
		var err error
		apiPwd, err = aes.DecryptFromB64(maputils.GetStringValue(devmap, "api_pwd", ""), m.Config.System.Aeskey)
		if err != nil {
			return nil, errors.New("GetMikrotikApi error, api passwd is invalid")
		}
	}

	if apiUser == "" || apiPwd == "" {
		return nil, errors.New("GetMikrotikApi error, auth data is invalid")
	}

	if _api, ok := m.DeviceConnPool.Get(apiAddr); ok {
		api := _api.(*mikrotik_api.MikrotikApi)
		if !api.CheckConnection() {
			log.Infof("reconnect to mikrotik device %s", apiAddr)
			api.ApiUser = apiUser
			api.ApiPwd = apiPwd
			api.ApiAddr = apiAddr
			rerr := api.ReConnect()
			if rerr != nil {
				return nil, fmt.Errorf("reconnect to mikrotik device %s fail %s", apiAddr, rerr.Error())
			}
		}
		return api, nil
	} else {
		api, err := mikrotik_api.GetConnection(apiUser, apiPwd, apiAddr, false)
		if err != nil {
			return nil, errors.New("GetMikrotikApi error, device conn error " + err.Error())
		}
		log.Infof("new connect to mikrotik device %s success", apiAddr)
		m.DeviceConnPool.Set(apiAddr, api)
		return api, nil
	}
}

// SyncMikrotikDeviceSysstatToElastic
// Mikrotik device  data statistics are transferred to Elasticsearch
func (m *MikrotikDeviceManager) SyncMikrotikDeviceStatToElastic(devtypes ...string) {
	for _, devtype := range devtypes {
		var devices = new(web.QueryResult)
		var err error
		switch devtype {
		case "cpe":
			devices, err = m.QueryItems(web.EmptyRequestParams(), TeamsacsCpe)
		case "vpe":
			devices, err = m.QueryItems(web.EmptyRequestParams(), TeamsacsVpe)
		default:
			err = errors.New("unsupported device types")
		}
		if err != nil {
			log.Error(err)
		}
		go m.SyncMikrotikDeviceSysstatToElastic(devtype, devices)
		go m.SyncMikrotikDeviceNetstatToElastic(devtype, devices)
	}
}

// SyncMikrotikDeviceSysstatToElastic
// Mikrotik device system data statistics are transferred to Elasticsearch
func (m *MikrotikDeviceManager) SyncMikrotikDeviceSysstatToElastic(devtype string, devices *web.QueryResult) {
	result := make([]elastic.TeamsacsLog, 0)
	for _, dev := range *devices {
		sn := maputils.GetStringValue(dev, "sn", "")
		if sn == "" {
			continue
		}
		sysstatlog := elastic.TeamsacsLog{
			Timestamp: time.Now().Format(time.RFC3339),
			Source:    m.Config.System.Appid,
			Sn:        sn,
			Name:      maputils.GetStringValue(dev, "name", ""),
			Tags:      maputils.GetStringValue(dev, "tags", ""),
			Model:     maputils.GetStringValue(dev, "model", ""),
			Version:   maputils.GetStringValue(dev, "sversion", ""),
			Devtype:   devtype,
			Sysstat: &elastic.DeviceSysstat{
				UpTime:     maputils.GetInt64Value(dev, "uptime", 0),
				MemPercent: maputils.GetInt64Value(dev, "memuse", 0),
				CpuPercent: maputils.GetInt64Value(dev, "cpuuse", 0),
			},
		}
		result = append(result, sysstatlog)
	}
	_, err := m.Elastic.BulkTeamslog(result...)
	if err != nil {
		log.Error(err)
	}
}

// SyncMikrotikDeviceNetstatToElastic
// Mikrotik device network data statistics are transferred to Elasticsearch
func (m *MikrotikDeviceManager) SyncMikrotikDeviceNetstatToElastic(devtype string, devices *web.QueryResult) {
	for _, dev := range *devices {
		// async stat interface
		sn := maputils.GetStringValue(dev, "sn", "")
		if sn == "" {
			continue
		}
		go func(gdev map[string]interface{}) {
			api, err := m.GetMikrotikApi(devtype, gdev)
			if err != nil {
				return
			}
			ifstats, err := api.GetInterfaceStats()
			if err != nil {
				return
			}

			nsesult := make([]elastic.TeamsacsLog, 0)
			for _, stmap := range ifstats {
				_interface := stmap["name"]
				nsesult = append(nsesult, elastic.TeamsacsLog{
					Timestamp: time.Now().Format(time.RFC3339),
					Source:    m.Config.System.Appid,
					Sn:        sn,
					Name:      maputils.GetStringValue(gdev, "name", ""),
					Tags:      maputils.GetStringValue(gdev, "tags", ""),
					Model:     maputils.GetStringValue(gdev, "model", ""),
					Version:   maputils.GetStringValue(gdev, "sversion", ""),
					Devtype:   devtype,
					// Conversion unit byte to bit
					Netstat: &elastic.DeviceNetstat{
						Interface:   _interface,
						SendBytes:   maputils.GetSInt64Value(stmap, "tx-byte", 0) * 8,
						RecvBytes:   maputils.GetSInt64Value(stmap, "rx-byte", 0) * 8,
						SendPackets: maputils.GetSInt64Value(stmap, "tx-packet", 0) * 8,
						RecvPackets: maputils.GetSInt64Value(stmap, "rx-packet", 0) * 8,
					},
				})
			}
			_, err = m.Elastic.BulkTeamslog(nsesult...)
			if err != nil {
				log.Error(err)
			}
		}(dev)
	}
}

func (m *MikrotikDeviceManager) QueryDeviceInterfaceList(sn, devtype string) ([]map[string]string, error) {
	api, err := m.GetMikrotikApiBySN(sn, devtype)
	if err != nil {
		return nil, err
	}
	return api.GetInterfaceList()
}

func (m *MikrotikDeviceManager) QueryDeviceRoutes(sn, devtype string) ([]map[string]string, error) {
	api, err := m.GetMikrotikApiBySN(sn, devtype)
	if err != nil {
		return nil, err
	}
	return api.GetIpRoutes()
}

func (m *MikrotikDeviceManager) QueryDeviceDnsinfo(sn, devtype string) (map[string]string, error) {
	api, err := m.GetMikrotikApiBySN(sn, devtype)
	if err != nil {
		return nil, err
	}
	return api.GetIpDnsInfo()
}
