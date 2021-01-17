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
	"time"

	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/mikrotik_api"
	"github.com/ca17/teamsacs/models/elastic"
)

type DeviceManager struct{ *ModelManager }

func (m *ModelManager) GetDeviceManager() *DeviceManager {
	store, _ := m.ManagerMap.Get("DeviceManager")
	return store.(*DeviceManager)
}

func _ifLtZeroInt64(s, defval int64) int64 {
	if s < 0 {
		return defval
	}
	return s
}

// SyncMikrotikDeviceToElastic
func (m *DeviceManager) SyncMikrotikDeviceToElastic(devtype string) error {
	var devices = new(web.QueryResult)
	var err error
	switch devtype {
	case "cpe":
		devices, err = m.QueryItems(web.EmptyRequestParams, TeamsacsCpe)
	case "vpe":
		devices, err = m.QueryItems(web.EmptyRequestParams, TeamsacsVpe)
	default:
		err = errors.New("unsupported device types")
	}
	if err != nil {
		return err
	}
	result := make([]elastic.TeamsacsLog, 0)
	for _, dev := range *devices {
		sn := maputils.GetStringValue(dev, "sn", "")
		if sn == "" {
			continue
		}
		stats, err := m.GetMikrotikDeviceStat(dev)
		if err != nil {
			continue
		}
		for _, s := range stats {
			var nstat = elastic.DeviceNetstat{
				Interface:   s.Interface,
				SendBytes:   0,
				RecvBytes:   0,
				SendPackets: 0,
				RecvPackets: 0,
			}
			cachekey := fmt.Sprintf("%s-%s", sn, s.Interface)
			m.DeviceStatCache.Set(cachekey, s)
			scache, ok := m.DeviceStatCache.Get(cachekey)
			if ok {
				var cacheStat = scache.(elastic.DeviceNetstat)
				nstat = elastic.DeviceNetstat{
					Interface:   s.Interface,
					SendBytes:   _ifLtZeroInt64(s.SendBytes-cacheStat.SendBytes, 0),
					RecvBytes:   _ifLtZeroInt64(s.RecvBytes-cacheStat.RecvBytes, 0),
					SendPackets: _ifLtZeroInt64(s.SendPackets-cacheStat.SendPackets, 0),
					RecvPackets: _ifLtZeroInt64(s.RecvPackets-cacheStat.RecvBytes, 0),
				}
			}

			tlog := elastic.TeamsacsLog{
				Timestamp: time.Now().Format(time.RFC3339),
				Source:    "teamsacs",
				Sn:        sn,
				Name:      maputils.GetStringValue(dev, "name", ""),
				Tags:      maputils.GetStringValue(dev, "tags", ""),
				Model:     maputils.GetStringValue(dev, "model", ""),
				Version:   maputils.GetStringValue(dev, "sversion", ""),
				Devtype:   "vpe",
				Netstat:   &nstat,
				Sysstat: &elastic.DeviceSysstat{
					UpTime:     maputils.GetInt64Value(dev, "uptime", 0),
					MemPercent: maputils.GetInt64Value(dev, "memuse", 0),
					CpuPercent: maputils.GetInt64Value(dev, "cpuuse", 0),
				},
			}
			result = append(result, tlog)
		}
	}
	_, err = m.Elastic.BulkTeamslog(result...)
	if err != nil {
		return err
	}
	log.Info("SyncMikrotikDeviceToElastic done")
	return nil
}

// GetMikrotikDeviceStat
func (m *DeviceManager) GetMikrotikDeviceStat(devmap map[string]interface{}) ([]elastic.DeviceNetstat, error) {
	apiUser := maputils.GetStringValue(devmap, "api_user", "")
	_apiPwd := maputils.GetStringValue(devmap, "api_pwd", "")
	apiPwd, err := aes.DecryptFromB64(_apiPwd, m.Config.System.Aeskey)
	if err != nil {
		return nil, errors.New("StatMikrotikInterface error, passwd Decrypt failure "+err.Error())
	}
	apiAddr := maputils.GetStringValue(devmap, "api_addr", "")
	if apiUser == "" || apiPwd == "" || apiAddr == "" {
		return nil, errors.New("StatMikrotikInterface error, auth data is invalid")
	}
	api, err := mikrotik_api.GetConnection(apiUser, apiPwd, apiAddr, false)
	if err != nil {
		return nil, errors.New("StatMikrotikInterface error, device conn error "+err.Error())
	}
	r, err := api.GetInterfaceStats()
	if err != nil {
		return nil, err
	}
	result := make([]elastic.DeviceNetstat, 0)
	for _, s := range r {
		inf := s["name"]
		tx := maputils.GetSInt64Value(s, "tx-byte", 0)
		rx := maputils.GetSInt64Value(s, "rx-byte", 0)
		txp := maputils.GetSInt64Value(s, "tx-packet", 0)
		rxp := maputils.GetSInt64Value(s, "rx-packet", 0)
		stat := elastic.DeviceNetstat{
			Interface:   inf,
			SendBytes:   tx,
			RecvBytes:   rx,
			SendPackets: txp,
			RecvPackets: rxp,
		}
		result = append(result, stat)
	}
	return result, nil
}
