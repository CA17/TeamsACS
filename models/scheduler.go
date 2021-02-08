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
	"time"

	"github.com/go-co-op/gocron"

	"github.com/ca17/teamsacs/common/log"
)

func (m *ModelManager) StartSched() error {
	m.Sched = gocron.NewScheduler(m.Location)

	_, _ = m.Sched.Every(1).Day().StartImmediately().Do(m.elkDataSync)
	time.Sleep(time.Second * 10)

	// Update mikrotik basic data regularly, read the latest data from Genieacs
	_, _ = m.Sched.Every(60).Second().StartImmediately().Do(func() {
		devinfos, err := m.GetGenieacsManager().QueryMikrotikDeviceInfo()
		if err != nil {
			log.Errorf("SyncAcsDeviceInfo error, query deviceInfo error %s", err.Error())
			return
		}
		log.Infof("fetch device num %d", len(devinfos))
		m.GetGenieacsManager().SyncMikrotikDeviceInfo(devinfos)
	})

	// Regularly count Mikrotik traffic data and system performance data to Elasticsearch
	_, _ = m.Sched.Every(60).Second().StartImmediately().Do(func() {
		m.GetMikrotikDeviceManager().SyncMikrotikDeviceStatToElastic("cpe", "vpe")
	})

	// Regularly clean up expired online users
	_, _ = m.Sched.Every(120).Second().StartImmediately().Do(func() {
		count, err := m.GetRadiusManager().ClearExpireOnline()
		if err != nil {
			log.Error(err)
		}
		log.Infof("clear expire online count %d", count)
	})

	<-m.Sched.Start()
	return nil
}


// elkDataSync
// Synchronize TeamsACS data to Elasticsearch for analysis
func (m *ModelManager) elkDataSync() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	synclist := []string{
		TeamsacsCpe,
		TeamsacsVpe,
		TeamsacsSubscribe,
		"product",
		"anode",
		"pnode",
		"application",
		"channel",
		"customer",
		"workload",
	}

	for _, name := range synclist {
		items, err := m.QueryItems(map[string]interface{}{}, name)
		if err != nil {
			continue
		}
		if items != nil {
			_, err := m.Elastic.BulkData("teamsacs_"+name, *items, true)
			if err != nil {
				log.Errorf("sync elk data %s error %s, %s", name, err.Error())
			}
		}
	}
}
