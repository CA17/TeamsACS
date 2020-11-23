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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ahmetb/go-linq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/validutil"
	"github.com/ca17/teamsacs/models/mikrotik"
)

type GenieacsManager struct{ *ModelManager }

func (m *ModelManager) GetGenieacsManager() *GenieacsManager {
	store, _ := m.ManagerMap.Get("GenieacsManager")
	return store.(*GenieacsManager)
}

// query device info by cpesn
func (m *GenieacsManager) QueryMikrotikDeviceInfo() ([]mikrotik.DeviceInfo, error) {
	items, err := m.QueryMikrotikSourceData("")
	if err != nil {
		return nil, err
	}

	result := make([]mikrotik.DeviceInfo, 0)
	linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   c.(mikrotik.TMap)["_id"],
				Value: mikrotik.GetObject(c, "Device.DeviceInfo"),
			}
		}).
		ForEachT(func(i linq.KeyValue) {
			var info = new(mikrotik.DeviceInfo)
			info.ParseBson(i.Value.(mikrotik.TMap))
			info.DeviceId = i.Key.(string)
			result = append(result, *info)
		})
	return result, nil
}

// query DNSServer by cpe sn
func (m *GenieacsManager) QueryMikrotikDnsServer(sn string) (*mikrotik.DnsClientServer, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}

	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.DNS.Client.Server"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.DnsClientServer {
			var r = new(mikrotik.DnsClientServer)
			r.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.T)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					r.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					rItem := new(mikrotik.DnsClientServerItem)
					rItem.Key = fmt.Sprintf("Device.DNS.Client.Server.%s", it.Key)
					rItem.ParseBson(it.Value.(mikrotik.TMap))
					r.Items = append(r.Items, *rItem)
				}
			})
			return r
		}).
		First()

	return result.(*mikrotik.DnsClientServer), nil
}

// query routers bu cpe sn
func (m *GenieacsManager) QueryMikrotikRouters(sn string) (*mikrotik.DeviceRouter, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}

	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.Routing.Router.1.IPv4Forwarding"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.DeviceRouter {
			var r = new(mikrotik.DeviceRouter)
			r.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.T)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					r.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					rItem := new(mikrotik.DeviceRouterItem)
					rItem.Key = fmt.Sprintf("Routing.Router.1.IPv4Forwarding.%s", it.Key)
					rItem.ParseBson(it.Value.(mikrotik.TMap))
					r.Items = append(r.Items, *rItem)
				}
			})
			return r
		}).
		First()

	return result.(*mikrotik.DeviceRouter), nil
}

// query all ethernet by cpe sn
func (m *GenieacsManager) QueryMikrotikEthernetInterface(sn string) (*mikrotik.EthernetInterface, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.Ethernet.Interface"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.EthernetInterface {
			var gei = new(mikrotik.EthernetInterface)
			gei.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					gei.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					eifItem := new(mikrotik.EthernetInterfaceItem)
					eifItem.Key = fmt.Sprintf("Ethernet.Interface.%s", it.Key)
					eifItem.ParseBson(it.Value.(mikrotik.TMap))
					gei.Items = append(gei.Items, *eifItem)
				}
			})
			return gei
		}).
		First()

	return result.(*mikrotik.EthernetInterface), nil
}

// query all ppp interface by cpe sn
func (m *GenieacsManager) QueryMikrotikPPPInterface(sn string) (*mikrotik.PPPInterface, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.PPP.Interface"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.PPPInterface {
			var gei = new(mikrotik.PPPInterface)
			gei.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					gei.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					eifItem := new(mikrotik.PPPInterfaceItem)
					eifItem.Key = fmt.Sprintf("PPP.Interface.%s", it.Key)
					eifItem.ParseBson(it.Value.(mikrotik.TMap))
					gei.Items = append(gei.Items, *eifItem)
				}
			})
			return gei
		}).
		First()

	return result.(*mikrotik.PPPInterface), nil
}

// query all ethernet by cpe sn
func (m *GenieacsManager) QueryMikrotikIpInterface(sn string) (*mikrotik.IpInterface, error) {
	items, err := m.QueryMikrotikSourceData(sn)
	if err != nil {
		return nil, err
	}
	var result = linq.From(items).
		Select(func(c mikrotik.T) mikrotik.T {
			return linq.KeyValue{
				Key:   mikrotik.GetObject(c, "Device.DeviceInfo.SerialNumber._value"),
				Value: mikrotik.GetObject(c, "Device.IP.Interface"),
			}
		}).
		SelectT(func(i linq.KeyValue) *mikrotik.IpInterface {
			var ipif = new(mikrotik.IpInterface)
			ipif.Sn = i.Key.(string)
			linq.From(i.Value.(mikrotik.TMap)).ForEachT(func(it linq.KeyValue) {
				if it.Key == "_timestamp" {
					ipif.Timestamp = mikrotik.ParseDateTime(it.Value)
				} else if validutil.IsInt(it.Key) {
					st := mikrotik.GetSubValue(it.Value.(mikrotik.TMap), "Status", "_value", false).(string)
					if st == "Up" {
						eifItem := new(mikrotik.IpInterfaceItem)
						eifItem.Key = fmt.Sprintf("IP.Interface.%s", it.Key)
						eifItem.ParseBson(it.Value.(mikrotik.TMap))
						ipif.Items = append(ipif.Items, *eifItem)
					}
				}
			})
			return ipif
		}).
		First()

	return result.(*mikrotik.IpInterface), nil
}

// query all cpe data
func (m *GenieacsManager) QueryMikrotikSourceData(sn string) ([]map[string]interface{}, error) {
	findOptions := options.Find()
	findOptions.SetLimit(100)
	coll := m.GetGenieAcsCollection(GenieacsDevices)
	var q = bson.M{}
	if sn != "" {
		q["Device.DeviceInfo.SerialNumber._value"] = sn
	}
	cur, err := coll.Find(context.TODO(), q, findOptions)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0)
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println(err)
		} else {
			items = append(items, elem)
		}
	}
	return items, nil
}

type GenieacsTask struct {
	Id        string    `json:"id"`
	TId       string    `json:"_id"`
	FileName  string    `json:"fileName"`
	FileType  string    `json:"fileType"`
	Name      string    `json:"name"`
	Device    string    `json:"device"`
	Timestamp time.Time `json:"timestamp"`
}

// GetAcsTaskDeviceIdList
// Get a list of device IDs for currently running tasks
func (m *GenieacsManager) GetAcsTaskDeviceIdList() ([]string, error) {
	result := make([]string, 0)
	resp := make([]GenieacsTask, 0)

	hresp, err := http.Get(common.UrlJoin(m.Config.Genieacs.NbiUrl, "/tasks"))
	if err != nil {
		return nil, err
	}

	defer hresp.Body.Close()
	body, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		return nil, fmt.Errorf("read tasks resp error %s", err.Error())
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("read tasks json resp error %s", err.Error())
	}
	for _, r := range resp {
		result = append(result, r.Device)
	}
	return result, nil
}
