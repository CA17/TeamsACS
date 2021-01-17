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
	"errors"
	"fmt"
	"time"

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

func GetConnection(apiUser string, apiPwd string, apiAddr string, TLS bool) (*MikrotikApi,error) {
	api := &MikrotikApi{ApiUser: apiUser, ApiPwd: apiPwd, ApiAddr: apiAddr, TLS: TLS}
	err := api.Connect()
	if err != nil {
		return nil ,err
	}
	return api, nil
}

func (a *MikrotikApi) Connect() error {
	var err error
	a.Client, err = routeros.DialTimeout(a.ApiAddr, a.ApiUser, a.ApiPwd, time.Second*3)
	if err != nil {
		return fmt.Errorf("connect mikrotik device error %s", err.Error())
	}
	return nil
}

// ExecuteCommand
// command: /interface/print
// params: "?xx=a?,yy=b"
func (a *MikrotikApi) ExecuteCommand(command string, params string, props string) (*routeros.Reply, error) {
	args := make([]string, 0)
	args = append(args, command)
	if params != "" {
		args = append(args, params)
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

// GetInterfaceStats
// item map:  "Map": { "name": "ether1",  "rx-byte": "176418199832", "rx-packet": "154145141","tx-byte": "26099616808","tx-packet": "93550768" }
func (a *MikrotikApi) GetInterfaceStats() ([]map[string]string, error) {
	var result = make([]map[string]string, 0)
	reply, err := a.ExecuteCommand("/interface/print", "stats", "name,rx-byte,tx-byte,rx-packet,tx-packet")
	if err != nil {
		return nil, fmt.Errorf("GetInterfaceStats error %s", err.Error())
	}
	for _, re := range reply.Re {
		result = append(result, re.Map)
	}
	return result, nil
}

// GetSystemResource
func (a *MikrotikApi) GetSystemResource() (map[string]string, error) {
	reply, err := a.ExecuteCommand("/system/resource/print", "",
		"uptime,version,free-memory,total-memory,cpu-load,free-hdd-space,total-hdd-space")
	if err != nil {
		return nil, fmt.Errorf("GetSystemResource error %s", err.Error())
	}
	if len(reply.Re) == 0 {
		return nil, errors.New("no result")
	}
	return reply.Re[0].Map, nil
}

// AddApiUser
// add  api user
func (a *MikrotikApi) AddApiUser(name, password string) error {
	_, err := a.Client.Run("/user/add", "=name="+name, "=password="+password, "=group=write", "=comment=ApiUser")
	if err != nil {
		return fmt.Errorf("AddApiUser Execute Api error %s", err.Error())
	}
	return nil
}
