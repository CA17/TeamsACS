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

package message

import (
	"strconv"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/constant"
)

// Set all event message processing functions
func (t *PubSubService) SetupEventBus() {
	// Post new CPE IP， Notify mosdns to add an IP
	_ = t.Manager.Bus.Subscribe(constant.EventMosdnsCpeUpdate, func(tunnelIP, ispIP string) {
		log.Infof("Receive event: %s (%s,%s)", constant.EventMosdnsCpeUpdate, tunnelIP, ispIP)
		_ = t.PublishToMosdns(&NnMessage{
			Uid:     strconv.FormatInt(common.UUIDint64(), 10),
			Command: MosdnsUpdateTeamsacsCpe,
			Attrs: map[string]interface{}{
				MosdnsTeamsacsCpeTunnelIP: tunnelIP,
				MosdnsTeamsacsCpeIspIp:    ispIP,
			},
		})
	})

	// Delete CPE IP， Notify mosdns to update IP
	_ = t.Manager.Bus.Subscribe(constant.EventMosdnsCpeRemove, func(tunnelIP string) {
		log.Infof("Receive event: %s (%s)", constant.EventMosdnsCpeUpdate, tunnelIP)
		_ = t.PublishToMosdns(&NnMessage{
			Uid:     strconv.FormatInt(common.UUIDint64(), 10),
			Command: MosdnsRemoveTeamsacsCpe,
			Attrs: map[string]interface{}{
				MosdnsTeamsacsCpeTunnelIP: tunnelIP,
			},
		})
	})

	// Clean CPE IP， Notify mosdns to clear the IP List
	_ = t.Manager.Bus.Subscribe(constant.EventMosdnsCpeClean, func() {
		log.Infof("Receive event: %s", constant.EventMosdnsCpeUpdate)
		_ = t.PublishToMosdns(&NnMessage{
			Uid:     strconv.FormatInt(common.UUIDint64(), 10),
			Command: MosdnsCleanTeamsacsCpe,
			Attrs: map[string]interface{}{},
		})
	})

}
