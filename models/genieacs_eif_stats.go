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
	"fmt"
	"time"

	"github.com/ahmetb/go-linq"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/validutil"
	"github.com/ca17/teamsacs/models/mikrotik"
)

// StatMikrotikEthernetInterface
func (m *GenieacsManager) StatMikrotikEthernetInterface(sn string) error {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return err
	}
	logs := make([]TeamsacsLog, 0)
	linq.From(items).ForEach(func(t mikrotik.T) {
		infoObj := mikrotik.GetObject(t, "Device.DeviceInfo")
		if infoObj == nil {
			return
		}
		// parse device info
		var info = new(mikrotik.DeviceInfo)
		info.ParseBson(infoObj.(mikrotik.TMap))

		// parse interface stat
		eifObj := mikrotik.GetObject(t, "Device.Ethernet.Interface")

		if eifObj == nil {
			log.Errorf("device %s stat data is nil", info.SerialNumber)
			return
		}
		var eif = new(mikrotik.EthernetInterface)
		eif.Sn = info.SerialNumber
		linq.From(eifObj.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
			if it.Key == "_timestamp" {
				eif.Timestamp = mikrotik.ParseDateTime(it.Value)
			} else if validutil.IsInt(it.Key) {
				eifItem := new(mikrotik.EthernetInterfaceItem)
				eifItem.Key = fmt.Sprintf("Ethernet.Interface.%s", it.Key)
				eifItem.ParseBson(it.Value.(mikrotik.TMap))
				// eif.Items = append(eif.Items, *eifItem)
				stat := eifItem.Stats

				getSysstat := func() * DeviceSysstat{
					if time.Now().Sub(info.Timestamp.Time()) >  time.Minute *2 {
						return nil
					}
					return &DeviceSysstat{
						Stattime: info.Timestamp.Time().Format(time.RFC3339),
						MemPercent:  info.MemoryUsage,
						CpuPercent:   info.CPUUsage,
						UpTime: info.UpTime,
					}
				}
				getNetstat := func() * DeviceNetstat{
					if time.Now().Sub(stat.Timestamp.Time()) >  time.Minute *2 {
						return nil
					}
					return &DeviceNetstat{
						Interface: eifItem.Key,
						Mac: eifItem.MACAddress,
						Stattime: stat.Timestamp.Time().Format(time.RFC3339),
						SendBytes:   stat.BytesSent,
						RecvBytes:   stat.BytesReceived,
						SendDrops:   stat.DiscardPacketsSent,
						RecvDrops:   stat.DiscardPacketsReceived,
						SendErrors:  stat.ErrorsSent,
						RecvErrors:  stat.ErrorsReceived,
						SendPackets: stat.PacketsSent,
						RecvPackets: stat.PacketsReceived,
					}
				}

				log := TeamsacsLog{
					Timestamp: time.Now().Format(time.RFC3339),
					Source:    "teamsacs",
					Sn:        info.SerialNumber,
					Name:      info.X_MIKROTIK_SystemIdentity,
					Tags:      "",
					Model:     info.ModelName,
					Version:   info.SoftwareVersion,
					Devtype:   "cpe",
					Sysstat:   getSysstat(),
					Netstat:   getNetstat(),
					Radiuslog: nil,
				}
				if log.Sysstat == nil && log.Netstat == nil {
					return
				}
				logs = append(logs, log)
			}
		})
	})

	_, err = m.Elastic.BulkTeamslog(logs...)
	if err != nil {
		return err
	}
	return nil
}
