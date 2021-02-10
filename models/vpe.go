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
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/constant"
)

// Vpe
// VPE is also a BRAS system
type Vpe map[string]interface{}

func (v Vpe) GetSecret() string {
	return maputils.GetStringValue(v, "secret", constant.NA)
}

func (v Vpe) GetVendorCode() string {
	return maputils.GetStringValue(v, "vendor_code", constant.NA)
}

func (v Vpe) GetIpaddr() string {
	return maputils.GetStringValue(v, "ipaddr", constant.NA)
}

func (v Vpe) GetCoaPort() int {
	return maputils.GetIntValue(v, "coa_port", 3799)
}

func (v Vpe) GetApiUser() (string, error) {
	return maputils.GetStringValueWithErr(v, "api_user")
}

func (v Vpe) GetApiPwd() (string, error) {
	return maputils.GetStringValueWithErr(v, "api_pwd")
}

func (v Vpe) GetApiAddr() (string, error) {
	return maputils.GetStringValueWithErr(v, "api_addr")
}

// VpeManager
type VpeManager struct{ *ModelManager }

func (m *ModelManager) GetVpeManager() *VpeManager {
	store, _ := m.ManagerMap.Get("VpeManager")
	return store.(*VpeManager)
}

func (m *VpeManager) QueryVpes(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsVpe)
}

// GetVpeByIpaddr
func (m *VpeManager) GetVpeByIpaddr(ip string) (*Vpe, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsVpe)
	doc := coll.FindOne(context.TODO(), bson.M{"ipaddr": ip})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Vpe)
	err = doc.Decode(result)
	return result, err
}

// GetVpeByIdentifier
func (m *VpeManager) GetVpeByIdentifier(identifier string) (*Vpe, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsVpe)
	doc := coll.FindOne(context.TODO(), bson.M{"identifier": identifier})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Vpe)
	err = doc.Decode(result)
	return result, err
}

func (m *VpeManager) ExistVpe(sn string) bool {
	coll := m.GetTeamsAcsCollection(TeamsacsVpe)
	count, _ := coll.CountDocuments(context.TODO(), bson.M{"sn": sn})
	return count > 0
}

func (m *VpeManager) GetVpeBySn(sn string) (*Cpe, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsVpe)
	doc := coll.FindOne(context.TODO(), bson.M{"sn": sn})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Cpe)
	err = doc.Decode(result)
	return result, err
}

// AddVpeData
func (m *VpeManager) AddVpeData(params web.RequestParams) error {
	data := params.GetParamMap("data")
	data["data_ver"] = common.GenerateDataVer()
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05 Z0700 MST")
	data["_id"] = common.UUID()
	coll := m.GetTeamsAcsCollection(TeamsacsVpe)
	var err error
	// If an api password is set, use aes encryption.
	apiPwd := data.GetString("api_pwd")
	if common.IsNotEmptyAndNA(apiPwd) {
		data["api_pwd"], err = aes.EncryptToB64(apiPwd, m.Config.System.Aeskey)
		if err != nil {
			return err
		}
	}
	go func() {
		_ = m.Elastic.AddData("teamsacs_vpe", data)
	}()
	_, err = coll.InsertOne(context.TODO(), data)
	return err
}

// UpdateVpeData
func (m *VpeManager) UpdateVpeData(params web.RequestParams) error {
	data := params.GetParamMap("data")
	data["data_ver"] = common.GenerateDataVer()
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05 Z0700 MST")
	_id := data.GetMustString("_id")
	var err error
	// If an api password is set, use aes encryption.
	apiPwd := data.GetString("api_pwd")
	if common.IsNotEmptyAndNA(apiPwd) {
		data["api_pwd"], err = aes.EncryptToB64(apiPwd, m.Config.System.Aeskey)
		if err != nil {
			return err
		}
	} else {
		delete(data, "api_pwd")
	}
	query := bson.M{"_id": _id}
	update := bson.M{"$set": data}
	go func() {
		_ = m.Elastic.UpdateData("teamsacs_vpe", data)
	}()
	_, err = m.GetTeamsAcsCollection(TeamsacsVpe).UpdateOne(context.TODO(), query, update)
	return err
}

// UpdateVpeBySn
func (m *VpeManager) UpdateVpeBySn(sn string, valmap map[string]interface{}) error {
	coll := m.GetTeamsAcsCollection(TeamsacsVpe)
	valmap["update_time"] = time.Now().Format("2006-01-02 15:04:05 Z0700 MST")
	update := bson.M{"$set": valmap}
	_, err := coll.UpdateOne(context.TODO(), bson.M{"sn": sn}, update)
	return err
}
