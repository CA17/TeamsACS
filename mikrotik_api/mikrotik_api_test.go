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
	"testing"
)

func getTestConn()  (*MikrotikApi, error){
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		return nil, err
	}
	return api, nil
}

func TestMikrotikApi_RemoveSocksUser(t *testing.T) {

	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}
	err = api.RemoveSocksUser("xxx")
	if err != nil {
		t.Error(err)
	}
}

func TestMikrotikApi_AddSocksUser(t *testing.T) {
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}

	err = api.AddSocksUser("xxx", "xxx", "10m/10m")
	if err != nil {
		t.Error(err)
	}
}

func TestMikrotikApi_ExecuteCommand(t *testing.T) {
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}
	r, err := api.ExecuteCommand("/interface/getall", "?running=true", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestMikrotikApi_GetSystemResource(t *testing.T) {
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}
	r, err := api.GetSystemResource()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestMikrotikApi_GetInterfaceStats(t *testing.T) {
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}
	r, err := api.GetInterfaceStats()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestMikrotikApi_CheckConnection(t *testing.T) {
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(api.CheckConnection())
}

func TestMikrotikApi_ReConnect(t *testing.T) {
	api := NewMikrotikApi("apiuser", "apipwd", "192.168.100.1:8728", false)
	err := api.Connect()
	if err != nil {
		t.Fatal(err)
	}
	api.Client.Close()
	t.Log(api.ReConnect())
	t.Log(api.CheckConnection())
}

func TestMikrotikApi_GetInterfaceList(t *testing.T) {
	api, err := getTestConn()
	if err != nil {
		t.Fatal(err)
	}
	ifs, err := api.GetInterfaceList()
	for _, l := range ifs {
		t.Logf("%+v", l)
	}
}

func TestMikrotikApi_GetIpRoutes(t *testing.T) {
	api, err := getTestConn()
	if err != nil {
		t.Fatal(err)
	}
	ifs, err := api.GetIpRoutes()
	for _, l := range ifs {
		t.Logf("%+v", l)
	}
}

func TestMikrotikApi_GetIpDnsInfos(t *testing.T) {
	api, err := getTestConn()
	if err != nil {
		t.Fatal(err)
	}
	info, err := api.GetIpDnsInfo()
	t.Logf("%+v", info)
}
