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
	"go.nanomsg.org/mangos/v3"

	"github.com/ca17/teamsacs/models"
)

const (
	// message define
	TeamsDNSUpdateCpe   = 0x1e
	TeamsDNSRemoveCpe   = 0x1f
	TeamsDNSCleanCpe    = 0x20
	TeamsDNSCpeTunnelIP = "TunnelIP"
	TeamsDNSCpeIspIp    = "IspIP"
	TeamsDNSCPETopic    = "TeamsDNS-CPE"
	TeamsDNSLOGTopic    = "TeamsDNS-LOG"
)


type NnMessage struct {
	Uid     string
	Command byte
	Attrs   map[string]interface{}
}

type PubSubService struct {
	Manager *models.ModelManager
	PubSock mangos.Socket
	SubSock mangos.Socket
}
