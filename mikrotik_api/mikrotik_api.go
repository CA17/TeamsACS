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

package mikrotik_api

import (
	"fmt"
	"strings"

	"gopkg.in/routeros.v2"

	"github.com/ca17/teamsacs/common/log"
)

type MikrotikApi struct {
	ApiUser string
	ApiPwd  string
	ApiAddr string
	TLS     bool
	Debug   bool
	Client  *routeros.Client
}

func NewMikrotikApi(apiUser string, apiPwd string, apiAddr string, TLS bool) *MikrotikApi {
	return &MikrotikApi{ApiUser: apiUser, ApiPwd: apiPwd, ApiAddr: apiAddr, TLS: TLS}
}

func (a *MikrotikApi) Connect() error {
	var err error
	a.Client, err = routeros.Dial(a.ApiAddr, a.ApiUser, a.ApiPwd)
	if err != nil {
		return fmt.Errorf("connect mikrotik device error %s", err.Error())
	}
	return nil
}

// ExecuteCommand
func (a *MikrotikApi) ExecuteCommand(command string, params string, props string) (*routeros.Reply, error) {
	args := make([]string, 0)
	args = append(args, command)
	for _, p := range strings.Split(params, ",") {
		if p == "" {
			continue
		}
		args = append(args, p)
	}
	if props != "" {
		args = append(args, "=.proplist="+props)
	}
	if a.Debug {
		log.Infof("Mikrotik Command Exec: %s %s %s", command, params, props)
	}
	reply, err := a.Client.Run(args...)
	if err != nil {
		return nil, fmt.Errorf("command[%s %s %s] exec error %s", command, params, props, err.Error())
	}
	if a.Debug {
		log.Infof("Mikrotik Command Exec Reply: %s", reply.String())
	}
	return reply, nil
}



// AddSocksUser
func (a *MikrotikApi) AddSocksUser(name, password, rateLimit string) error {
	_, err := a.Client.Run("/ip/socks/users/add", "=name="+name, "=password="+password, "=only-one=yes", "=rate-limit="+rateLimit)
	if err != nil {
		return fmt.Errorf("AddSocksUser Execute Api error %s", err.Error())
	}
	return nil
}

// RemoveSocksUser
func (a *MikrotikApi) RemoveSocksUser(name string) error {
	reply, err := a.Client.Run("/ip/socks/users/getall", "?name="+name, "=.proplist=.id")
	if err != nil {
		return fmt.Errorf("RemoveSocksUser find error %s", err.Error())
	}
	if a.Debug {
		log.Info(reply.String())
	}
	for _, re := range reply.Re {
		_, err := a.Client.Run("/ip/socks/users/remove", "=.id="+re.Map[".id"])
		if err != nil {
			return fmt.Errorf("RemoveSocksUser  error %s", err.Error())
		}
	}
	return nil
}


// AddPPPUser
// add  ppp user with fix ip
func (a *MikrotikApi) AddPPPUser(name, password, ip, gateway, remark string) error {
	_, err := a.Client.Run("/ppp/secret/add", "=name="+name, "=password="+password, "=local-address="+gateway, "=comment="+remark, "=remote-address="+ip)
	if err != nil {
		return fmt.Errorf("AddPPPUser Execute Api error %s", err.Error())
	}
	return nil
}

// RemovePPPUser
func (a *MikrotikApi) RemovePPPUser(name string) error {
	reply, err := a.Client.Run("/ppp/secret/getall", "?name="+name, "=.proplist=.id")
	if err != nil {
		return fmt.Errorf("RemovePPPUser find error %s", err.Error())
	}
	if a.Debug {
		log.Info(reply.String())
	}
	for _, re := range reply.Re {
		_, err := a.Client.Run("/ppp/secret/remove", "=.id="+re.Map[".id"])
		if err != nil {
			return fmt.Errorf("RemovePPPUser  error %s", err.Error())
		}
	}
	return nil
}


