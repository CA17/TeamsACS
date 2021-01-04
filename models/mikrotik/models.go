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

package mikrotik

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ca17/teamsacs/common/validutil"
)

type T = interface{}
type TMap = map[string]interface{}

var _emptydate, _ = time.Parse("2006-01-02 15:04:05 Z0700 MST", "1979-11-30 00:00:00 +0000 GMT")
var EmptyDate = primitive.NewDateTimeFromTime(_emptydate)

// GetObject
// Read nested objects
func GetObject(m interface{}, key string) interface{} {
	return linq.From(strings.Split(key, ".")).AggregateWithSeed(m, func(i T, i2 T) interface{} {
		switch i.(type) {
		case TMap:
			v, ok := i.(TMap)[i2.(string)]
			if !ok {
				return nil
			}
			return v
		default:
			return nil
		}
	})
}

func GetSubValue(m TMap, key string, secondKey string, defval interface{}) interface{} {
	v, ok := m[key]
	if !ok {
		return defval
	}
	sv, ok := v.(TMap)
	if !ok {
		return defval
	}
	subv, ok := sv[secondKey]
	if !ok {
		return defval
	}
	return subv
}

func GetFloat64SubValue(m TMap, key string, secondKey string, defval interface{}) float64 {
	val := GetSubValue(m, key, secondKey, defval)
	return ParseFloat64(val)
}

func GetInt64SubValue(m TMap, key string, secondKey string, defval interface{}) int64 {
	val := GetSubValue(m, key, secondKey, defval)
	return ParseInt64(val)
}

func ParseInt64(val interface{}) int64 {
	switch val.(type) {
	case int:
		return int64(val.(int))
	case int32:
		return int64(val.(int32))
	case int64:
		return val.(int64)
	case float32:
		return int64(val.(float32))
	case float64:
		return int64(val.(float64))
	case string:
		v, _ := strconv.ParseInt(val.(string), 10, 64)
		return v
	default:
		return 0
	}
}

func ParseFloat64(val interface{}) float64 {
	switch val.(type) {
	case int:
		return float64(val.(int))
	case int32:
		return float64(val.(int32))
	case int64:
		return float64(val.(int64))
	case float32:
		return float64(val.(float32))
	case float64:
		return val.(float64)
	case string:
		v, _ := strconv.ParseFloat(val.(string), 64)
		return v
	default:
		return 0
	}
}

func ParseDateTime(val interface{}) primitive.DateTime {
	switch val.(type) {
	case primitive.DateTime:
		return val.(primitive.DateTime)
	default:
		return EmptyDate
	}
}

type DeviceInfo struct {
	DeviceId                  string
	Description               string
	HardwareVersion           string
	Manufacturer              string
	ModelName                 string
	ManufacturerOUI           string
	ProductClass              string
	ProvisioningCode          string
	SoftwareVersion           string
	SerialNumber              string
	X_MIKROTIK_ArchName       string
	X_MIKROTIK_SystemIdentity string
	UpTime                    int64
	MemoryUsage               int64
	CPUUsage                  int64
	Timestamp                 primitive.DateTime
}

// 解析设备信息
func (info *DeviceInfo) ParseBson(val map[string]interface{}) {
	info.Description = GetSubValue(val, "Description", "_value", "").(string)
	info.HardwareVersion = GetSubValue(val, "HardwareVersion", "_value", "").(string)
	info.Manufacturer = GetSubValue(val, "Manufacturer", "_value", "").(string)
	info.ProductClass = GetSubValue(val, "ProductClass", "_value", "").(string)
	info.ModelName = GetSubValue(val, "ModelName", "_value", "").(string)
	info.ProvisioningCode = GetSubValue(val, "ProvisioningCode", "_value", "").(string)
	info.ManufacturerOUI = GetSubValue(val, "ManufacturerOUI", "_value", "").(string)
	info.SoftwareVersion = GetSubValue(val, "SoftwareVersion", "_value", "").(string)
	info.SerialNumber = GetSubValue(val, "SerialNumber", "_value", "").(string)
	info.X_MIKROTIK_ArchName = GetSubValue(val, "X_MIKROTIK_ArchName", "_value", "").(string)
	info.X_MIKROTIK_SystemIdentity = GetSubValue(val, "X_MIKROTIK_SystemIdentity", "_value", "").(string)
	info.UpTime = GetInt64SubValue(val, "UpTime", "_value", "")
	info.Timestamp = ParseDateTime(val["_timestamp"])

	MemoryStatus, ok := val["MemoryStatus"]
	if ok {
		Free := GetFloat64SubValue(MemoryStatus.(TMap), "Free", "_value", "")
		Total := GetFloat64SubValue(MemoryStatus.(TMap), "Total", "_value", "")
		info.MemoryUsage = int64(math.Round((Total - Free) / Total * 100))
	}

	ProcessStatus, ok := val["ProcessStatus"]
	if ok {
		info.CPUUsage = GetInt64SubValue(ProcessStatus.(TMap), "CPUUsage", "_value", "")
	}
}

type DeviceEthernet struct {
	Interfaces []map[string]EthernetInterface
	Timestamp  time.Time
}

type EthernetInterface struct {
	Sn        string                  `json:"sn"`
	Items     []EthernetInterfaceItem `json:"items"`
	Timestamp primitive.DateTime      `json:"timestamp"`
}

type EthernetInterfaceItem struct {
	Key         string                     `json:"key"`
	Enable      bool                       `json:"enable"`
	LowerLayers string                     `json:"lower_layers"`
	MACAddress  string                     `json:"mac_address"`
	Stats       EthernetInterfaceItemStats `json:"stats"`
	Status      string                     `json:"status"`
	Timestamp   primitive.DateTime         `json:"timestamp"`
}

func (v *EthernetInterfaceItem) ParseBson(val map[string]interface{}) {
	v.Enable = GetSubValue(val, "Enable", "_value", false).(bool)
	v.LowerLayers = GetSubValue(val, "LowerLayers", "_value", "").(string)
	v.MACAddress = GetSubValue(val, "MACAddress", "_value", "").(string)
	v.Status = GetSubValue(val, "Status", "_value", "").(string)
	v.Timestamp = ParseDateTime(val["_timestamp"])
	var stats = new(EthernetInterfaceItemStats)
	stats.ParseBson(val["Stats"].(TMap))
	v.Stats = *stats
}

type EthernetInterfaceItemStats struct {
	BytesReceived          int64              `json:"bytes_received"`
	BytesSent              int64              `json:"bytes_sent"`
	DiscardPacketsReceived int64              `json:"discard_packets_received"`
	DiscardPacketsSent     int64              `json:"discard_packets_sent"`
	ErrorsReceived         int64              `json:"errors_received"`
	ErrorsSent             int64              `json:"errors_sent"`
	PacketsReceived        int64              `json:"packets_received"`
	PacketsSent            int64              `json:"packets_sent"`
	Timestamp              primitive.DateTime `json:"timestamp"`
}

func (v *EthernetInterfaceItemStats) ParseBson(valmap map[string]interface{}) {
	v.BytesReceived = GetInt64SubValue(valmap, "BytesReceived", "_value", 0)
	v.BytesSent = GetInt64SubValue(valmap, "BytesSent", "_value", 0)
	v.DiscardPacketsReceived = GetInt64SubValue(valmap, "DiscardPacketsReceived", "_value", 0)
	v.DiscardPacketsSent = GetInt64SubValue(valmap, "DiscardPacketsSent", "_value", 0)
	v.ErrorsReceived = GetInt64SubValue(valmap, "ErrorsReceived", "_value", 0)
	v.ErrorsSent = GetInt64SubValue(valmap, "ErrorsSent", "_value", 0)
	v.PacketsReceived = GetInt64SubValue(valmap, "PacketsReceived", "_value", 0)
	v.PacketsSent = GetInt64SubValue(valmap, "PacketsSent", "_value", 0)
	v.Timestamp = ParseDateTime(valmap["_timestamp"])
}

// 路由表
type DeviceRouter struct {
	Sn        string             `json:"sn"`
	Items     []DeviceRouterItem `json:"items"`
	Timestamp primitive.DateTime `json:"timestamp"`
}

// 路由表项目
type DeviceRouterItem struct {
	Key              string             `json:"key"`
	Enable           bool               `json:"enable"`
	DestIPAddress    string             `json:"dest_ip_address"`
	DestSubnetMask   string             `json:"dest_subnet_mask"`
	GatewayIPAddress string             `json:"gateway_ip_address"`
	Interface        string             `json:"interface"`
	Origin           string             `json:"origin"`
	StaticRoute      bool               `json:"static_route"`
	Status           string             `json:"status"`
	Timestamp        primitive.DateTime `json:"timestamp"`
}

// 解析路由表
func (v *DeviceRouterItem) ParseBson(val map[string]interface{}) {
	v.Enable = GetSubValue(val, "Enable", "_value", false).(bool)
	v.StaticRoute = GetSubValue(val, "StaticRoute", "_value", false).(bool)
	v.DestIPAddress = GetSubValue(val, "DestIPAddress", "_value", "").(string)
	v.DestSubnetMask = GetSubValue(val, "DestSubnetMask", "_value", "").(string)
	v.GatewayIPAddress = GetSubValue(val, "GatewayIPAddress", "_value", "").(string)
	v.Interface = GetSubValue(val, "Interface", "_value", "").(string)
	v.Origin = GetSubValue(val, "Origin", "_value", "").(string)
	v.Status = GetSubValue(val, "Status", "_value", "").(string)
	v.Timestamp = ParseDateTime(val["_timestamp"])
}

type DnsClientServer struct {
	Sn        string                `json:"sn"`
	Items     []DnsClientServerItem `json:"items"`
	Timestamp primitive.DateTime    `json:"timestamp"`
}

type DnsClientServerItem struct {
	Key       string             `json:"key"`
	DnsServer string             `json:"dns_server"`
	Enable    bool               `json:"enable"`
	Status    string             `json:"status"`
	Type      string             `json:"type"`
	Timestamp primitive.DateTime `json:"timestamp"`
}

func (v *DnsClientServerItem) ParseBson(val map[string]interface{}) {
	v.DnsServer = GetSubValue(val, "DNSServer", "_value", false).(string)
	v.Status = GetSubValue(val, "Status", "_value", false).(string)
	v.Type = GetSubValue(val, "Type", "_value", false).(string)
	v.Enable = GetSubValue(val, "Enable", "_value", false).(bool)
	v.Timestamp = ParseDateTime(val["_timestamp"])

}

type IpInterface struct {
	Sn        string             `json:"sn"`
	Items     []IpInterfaceItem  `json:"items"`
	Timestamp primitive.DateTime `json:"timestamp"`
}

type IpInterfaceItem struct {
	Key         string                           `json:"key"`
	Enable      bool                             `json:"enable"`
	LowerLayers string                           `json:"lower_layers"`
	Status      string                           `json:"status"`
	Type        string                           `json:"type"`
	Ipv4Items   []IpInterfaceItemIpv4AddressItem `json:"ipv4_items"`
	Timestamp   primitive.DateTime               `json:"timestamp"`
}

type IpInterfaceItemIpv4AddressItem struct {
	Key       string             `json:"key"`
	Type      string             `json:"type"`
	Enable    bool               `json:"enable"`
	IpAddress string             `json:"ip_address"`
	Status    string             `json:"status"`
	Netmask   string             `json:"netmask"`
	Timestamp primitive.DateTime `json:"timestamp"`
}

func (v *IpInterfaceItem) ParseBson(val map[string]interface{}) {
	v.Status = GetSubValue(val, "Status", "_value", "").(string)
	v.LowerLayers = GetSubValue(val, "LowerLayers", "_value", "").(string)
	v.Type = GetSubValue(val, "Type", "_value", "").(string)
	v.Enable = GetSubValue(val, "Enable", "_value", false).(bool)
	v.Timestamp = ParseDateTime(val["_timestamp"])

	ipv4addrs, ok := val["IPv4Address"]
	if !ok {
		return
	}

	for ik, iv := range ipv4addrs.(TMap) {
		if validutil.IsInt(ik) {
			var item = new(IpInterfaceItemIpv4AddressItem)
			item.Key = fmt.Sprintf("IP.Interface.%s.IPv4Address.%s", v.Key, ik)
			item.ParseBson(iv.(TMap))
			v.Ipv4Items = append(v.Ipv4Items, *item)
		}
	}

}

func (v *IpInterfaceItemIpv4AddressItem) ParseBson(val map[string]interface{}) {
	v.IpAddress = GetSubValue(val, "IPAddress", "_value", "").(string)
	v.Status = GetSubValue(val, "Status", "_value", "").(string)
	v.Type = GetSubValue(val, "AddressingType", "_value", "").(string)
	v.Netmask = GetSubValue(val, "SubnetMask", "_value", "").(string)
	v.Enable = GetSubValue(val, "Enable", "_value", false).(bool)
	v.Timestamp = ParseDateTime(val["_timestamp"])
}

type PPPInterface struct {
	Sn        string             `json:"sn"`
	Items     []PPPInterfaceItem `json:"items"`
	Timestamp primitive.DateTime `json:"timestamp"`
}

type PPPInterfaceItem struct {
	Key                string                `json:"key"`
	Enable             bool                  `json:"enable"`
	LowerLayers        string                `json:"lower_layers"`
	Stats              PPPInterfaceItemStats `json:"stats"`
	Status             string                `json:"status"`
	ConnectionStatus   string                `json:"connection_status"`
	ConnectionTrigger  string                `json:"connection_trigger"`
	EncryptionProtocol string                `json:"encryption_protocol"`
	RemoteIPAddress    string                `json:"remote_ip_address"`
	LocalIPAddress     string                `json:"local_ip_address"`
	Username           string                `json:"username"`
	Timestamp          primitive.DateTime    `json:"timestamp"`
}

func (v *PPPInterfaceItem) ParseBson(val map[string]interface{}) {
	v.Enable = GetSubValue(val, "Enable", "_value", false).(bool)
	v.LowerLayers = GetSubValue(val, "LowerLayers", "_value", "").(string)
	v.ConnectionStatus = GetSubValue(val, "ConnectionStatus", "_value", "").(string)
	v.ConnectionTrigger = GetSubValue(val, "ConnectionTrigger", "_value", "").(string)
	v.EncryptionProtocol = GetSubValue(val, "EncryptionProtocol", "_value", "").(string)
	v.RemoteIPAddress = GetObject(val, "IPCP.RemoteIPAddress._value").(string)
	v.LocalIPAddress = GetObject(val, "IPCP.LocalIPAddress._value").(string)
	v.Username = GetSubValue(val, "Username", "_value", "").(string)
	v.Status = GetSubValue(val, "Status", "_value", "").(string)
	v.Timestamp = ParseDateTime(val["_timestamp"])
	var stats = new(PPPInterfaceItemStats)
	stats.ParseBson(val["Stats"].(TMap))
	v.Stats = *stats
}

type PPPInterfaceItemStats struct {
	BytesReceived          int64              `json:"bytes_received"`
	BytesSent              int64              `json:"bytes_sent"`
	DiscardPacketsReceived int64              `json:"discard_packets_received"`
	DiscardPacketsSent     int64              `json:"discard_packets_sent"`
	ErrorsReceived         int64              `json:"errors_received"`
	ErrorsSent             int64              `json:"errors_sent"`
	PacketsReceived        int64              `json:"packets_received"`
	PacketsSent            int64              `json:"packets_sent"`
	Timestamp              primitive.DateTime `json:"timestamp"`
}

func (v *PPPInterfaceItemStats) ParseBson(valmap map[string]interface{}) {
	v.BytesReceived = GetInt64SubValue(valmap, "BytesReceived", "_value", 0)
	v.BytesSent = GetInt64SubValue(valmap, "BytesSent", "_value", 0)
	v.DiscardPacketsReceived = GetInt64SubValue(valmap, "DiscardPacketsReceived", "_value", 0)
	v.DiscardPacketsSent = GetInt64SubValue(valmap, "DiscardPacketsSent", "_value", 0)
	v.ErrorsReceived = GetInt64SubValue(valmap, "ErrorsReceived", "_value", 0)
	v.ErrorsSent = GetInt64SubValue(valmap, "ErrorsSent", "_value", 0)
	v.PacketsReceived = GetInt64SubValue(valmap, "PacketsReceived", "_value", 0)
	v.PacketsSent = GetInt64SubValue(valmap, "PacketsSent", "_value", 0)
}
