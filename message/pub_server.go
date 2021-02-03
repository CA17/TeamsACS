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
	"bytes"

	"github.com/pkg/errors"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"

	"github.com/vmihailenco/msgpack/v5"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

const (
	// mosdns message define
	MosdnsUpdateTeamsacsCpe   = 0x1e
	MosdnsRemoveTeamsacsCpe   = 0x1f
	MosdnsCleanTeamsacsCpe   = 0x20
	MosdnsTeamsacsCpeTunnelIP = "TunnelIP"
	MosdnsTeamsacsCpeIspIp    = "IspIP"
	MosdnsTopic               = "mosdns"
)

type NnMessage struct {
	Uid     string
	Command byte
	Attrs   map[string]interface{}
}

type PubSubService struct {
	Manager *models.ModelManager
	PubSock mangos.Socket
}

func NewPubSubService(manager *models.ModelManager) *PubSubService {
	serv := &PubSubService{Manager: manager}
	return serv
}

// public the message
func (t *PubSubService) PublishToMosdns(msg *NnMessage) error {
	_msg, err := msgpack.Marshal(msg)
	if err != nil {
		return err
	}
	var buff = bytes.NewBuffer([]byte(MosdnsTopic))
	buff.Write(_msg)
	return t.PubSock.Send(buff.Bytes())
}

func (t *PubSubService) StartPubServer() error {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		return err
	}
	if err = sock.Listen(t.Manager.Config.Messaged.Address); err != nil {
		log.Errorf("%+v", errors.WithStack(err))
		return err
	}
	t.PubSock = sock
	log.Infof("PubServer listen %s...", t.Manager.Config.Messaged.Address)
	t.SetupEventBus()
	return nil
}
