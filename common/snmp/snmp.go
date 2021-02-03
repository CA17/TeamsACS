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

package snmp

import (
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"

	"github.com/ca17/teamsacs/common/snmp/mibs/ifmib"
)

type DeviceInterface struct {
	Name      string `json:"name"`
	Index     int64  `json:"index"`
	Type      int64  `json:"type"`
	InOctets  int64  `json:"in_octets"`
	OutOctets int64  `json:"out_octets"`
}

type SnmpV2Client struct {
	Snmpc *gosnmp.GoSNMP
}

func NewSnmpV2Client(target string, port int, community string) *SnmpV2Client {
	return &SnmpV2Client{
		Snmpc: &gosnmp.GoSNMP{
			Port:               uint16(port),
			Transport:          "udp",
			Community:          community,
			Version:            gosnmp.Version2c,
			Timeout:            time.Duration(3) * time.Second,
			Retries:            3,
			ExponentialTimeout: true,
			MaxOids:            gosnmp.MaxOids,
			Target:             target,
		},
	}

}

// SNMP Get interface definition data
func (c *SnmpV2Client) QueryInterfaces() (map[int64]*DeviceInterface, error) {
	// Get interface index
	indexs, err := c.Snmpc.BulkWalkAll(ifmib.IF_MIB_ifIndex_OID)
	if err != nil {
		return nil, err
	}
	ifmap := make(map[int64]*DeviceInterface)
	for _, idx := range indexs {
		idxval := gosnmp.ToBigInt(idx.Value).Int64()
		ifmap[idxval] = &DeviceInterface{
			Index: gosnmp.ToBigInt(idx.Value).Int64(),
		}
	}
	names, err := c.Snmpc.BulkWalkAll(ifmib.IF_MIB_ifDescr_OID)
	if err != nil {
		return nil, err
	}
	// 获取接口名称
	for _, name := range names {
		nameoid := name.Name
		idxval, err := strconv.ParseInt(nameoid[strings.LastIndex(nameoid, ".")+1:], 10, 64)
		if err != nil {
			continue
		}
		if _, flag := ifmap[idxval]; flag {
			ifmap[idxval].Name = string(name.Value.([]byte))
		}
	}

	types, err := c.Snmpc.BulkWalkAll(ifmib.IF_MIB_ifType_OID)
	if err != nil {
		return nil, err
	}

	// Get interface name
	for _, _type := range types {
		nameoid := _type.Name
		idxval, err := strconv.ParseInt(nameoid[strings.LastIndex(nameoid, ".")+1:], 10, 64)
		if err != nil {
			continue
		}
		if _, flag := ifmap[idxval]; flag {
			ifmap[idxval].Type = gosnmp.ToBigInt(_type.Value).Int64()
		}
	}

	return ifmap, nil

}

func (c *SnmpV2Client) QueryInterfacesInOctets(ifmap map[int64]*DeviceInterface) error {
	bytes, err := c.Snmpc.BulkWalkAll(ifmib.IF_MIB_ifInOctets_OID)
	if err != nil {
		return err
	}
	// Get interface name
	for _, _item := range bytes {
		nameoid := _item.Name
		idxval, err := strconv.ParseInt(nameoid[strings.LastIndex(nameoid, ".")+1:], 10, 64)
		if err != nil {
			continue
		}
		if _, flag := ifmap[idxval]; flag {
			ifmap[idxval].InOctets = gosnmp.ToBigInt(_item.Value).Int64()
		}
	}
	return nil
}

func (c *SnmpV2Client) QueryInterfacesOutOctets(ifmap map[int64]*DeviceInterface) error {
	bytes, err := c.Snmpc.BulkWalkAll(ifmib.IF_MIB_ifOutOctets_OID)
	if err != nil {
		return err
	}
	// Get interface name
	for _, _item := range bytes {
		nameoid := _item.Name
		idxval, err := strconv.ParseInt(nameoid[strings.LastIndex(nameoid, ".")+1:], 10, 64)
		if err != nil {
			continue
		}
		if _, flag := ifmap[idxval]; flag {
			ifmap[idxval].OutOctets = gosnmp.ToBigInt(_item.Value).Int64()
		}
	}
	return nil
}

// Connect to device
func (c *SnmpV2Client) Connect() error {
	return c.Snmpc.Connect()
}

// Close the connection
func (c *SnmpV2Client) Close() error {
	return c.Snmpc.Conn.Close()
}
