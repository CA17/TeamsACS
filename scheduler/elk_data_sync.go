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

package scheduler

import (
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)


func ElkDataSync(manager *models.ModelManager) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	synclist := []string{
		models.TeamsacsCpe,
		models.TeamsacsVpe,
		models.TeamsacsSubscribe,
		"product",
		"anode",
		"pnode",
		"application",
		"channel",
		"customer",
		"workload",
	}

	for _, name := range synclist {
		items, err := manager.QueryItems(map[string]interface{}{}, name)
		if err != nil {
			continue
		}
		if items != nil {
			_, err := manager.Elastic.BulkData("teamsacs_"+name, *items, true)
			if err != nil {
				log.Errorf("sync elk data %s error %s, %s", name, err.Error())
			}
		}
	}
}
