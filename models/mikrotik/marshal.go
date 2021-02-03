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

import "encoding/json"

func (d EthernetInterfaceItem) MarshalJSON() ([]byte, error) {
	type Alias EthernetInterfaceItem
	return json.Marshal(&struct {
		Alias
		Timestamp string `db:"timestamp" json:"timestamp" form:"-" query:"-"  `
	}{
		Alias:     (Alias)(d),
		Timestamp: d.Timestamp.Time().Format("2006-01-02 15:04:05"),
	})
}

func (d DnsClientServerItem) MarshalJSON() ([]byte, error) {
	type Alias DnsClientServerItem
	return json.Marshal(&struct {
		Alias
		Timestamp string `db:"timestamp" json:"timestamp" form:"-" query:"-"  `
	}{
		Alias:     (Alias)(d),
		Timestamp: d.Timestamp.Time().Format("2006-01-02 15:04:05"),
	})
}

func (d DeviceRouterItem) MarshalJSON() ([]byte, error) {
	type Alias DeviceRouterItem
	return json.Marshal(&struct {
		Alias
		Timestamp string `db:"timestamp" json:"timestamp" form:"-" query:"-"  `
	}{
		Alias:     (Alias)(d),
		Timestamp: d.Timestamp.Time().Format("2006-01-02 15:04:05"),
	})
}

func (d IpInterfaceItemIpv4AddressItem) MarshalJSON() ([]byte, error) {
	type Alias IpInterfaceItemIpv4AddressItem
	return json.Marshal(&struct {
		Alias
		Timestamp string `db:"timestamp" json:"timestamp" form:"-" query:"-"  `
	}{
		Alias:     (Alias)(d),
		Timestamp: d.Timestamp.Time().Format("2006-01-02 15:04:05"),
	})
}

func (d IpInterfaceItem) MarshalJSON() ([]byte, error) {
	type Alias IpInterfaceItem
	return json.Marshal(&struct {
		Alias
		Timestamp string `db:"timestamp" json:"timestamp" form:"-" query:"-"  `
	}{
		Alias:     (Alias)(d),
		Timestamp: d.Timestamp.Time().Format("2006-01-02 15:04:05"),
	})
}

func (d PPPInterfaceItem) MarshalJSON() ([]byte, error) {
	type Alias PPPInterfaceItem
	return json.Marshal(&struct {
		Alias
		Timestamp string `db:"timestamp" json:"timestamp" form:"-" query:"-"  `
	}{
		Alias:     (Alias)(d),
		Timestamp: d.Timestamp.Time().Format("2006-01-02 15:04:05"),
	})
}
