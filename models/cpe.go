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
)

// Cpe
type Cpe map[string]interface{}

func (v Cpe) Set(key string, value interface{}) {
	v[key] = value
}

func (v Cpe) GetApiUser() (string, error) {
	return maputils.GetStringValueWithErr(v, "api_user")
}

func (v Cpe) GetApiPwd() (string, error) {
	return maputils.GetStringValueWithErr(v, "api_pwd")
}

func (v Cpe) GetApiAddr() (string, error) {
	return maputils.GetStringValueWithErr(v, "api_addr")
}

func (v Cpe) GetDeviceId() string {
	return maputils.GetStringValue(v, "device_id", "")
}

func (v Cpe) GetID() string {
	return maputils.GetStringValue(v, "_id", "")
}

type CpeManager struct{ *ModelManager }

func (m *ModelManager) GetCpeManager() *CpeManager {
	store, _ := m.ManagerMap.Get("CpeManager")
	return store.(*CpeManager)
}

func (m *CpeManager) GetCpeByDeviceId(device_id string) (*Cpe, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsCpe)
	doc := coll.FindOne(context.TODO(), bson.M{"device_id": device_id})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Cpe)
	err = doc.Decode(result)
	return result, err
}

func (m *CpeManager) GetCpeBySn(sn string) (*Cpe, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsCpe)
	doc := coll.FindOne(context.TODO(), bson.M{"sn": sn})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Cpe)
	err = doc.Decode(result)
	return result, err
}

func (m *CpeManager) QueryCpes(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsCpe)
}

func (m *CpeManager) QueryCpeList(params web.RequestParams) (*web.QueryResult, error) {
	return m.QueryItems(params, TeamsacsCpe)
}

func (m *CpeManager) ExistCpe(sn string) bool {
	coll := m.GetTeamsAcsCollection(TeamsacsCpe)
	count, _ := coll.CountDocuments(context.TODO(), bson.M{"sn": sn})
	return count > 0
}

func (m *CpeManager) AddCpeData(params web.RequestParams) error {
	data := params.GetParamMap("data")
	_id := data.GetString("_id")
	if common.IsEmptyOrNA(_id) {
		data["_id"] = common.UUID()
	}
	coll := m.GetTeamsAcsCollection(TeamsacsCpe)
	var err error
	// If an api password is set, use aes encryption.
	apiPwd := data.GetString("api_pwd")
	if common.IsNotEmptyAndNA(apiPwd) {
		data["api_pwd"], err = aes.EncryptToB64(apiPwd, m.Config.System.Aeskey)
		if err != nil {
			return err
		}
	}
	_, err = coll.InsertOne(context.TODO(), data)
	return err
}

func (m *CpeManager) AddCpeDataMap(cpe Cpe) error {
	_id := cpe.GetID()
	if common.IsEmptyOrNA(_id) {
		cpe["_id"] = common.UUID()
	}
	coll := m.GetTeamsAcsCollection(TeamsacsCpe)
	var err error
	// If an api password is set, use aes encryption.
	apiPwd, _ := cpe.GetApiPwd()
	if common.IsNotEmptyAndNA(apiPwd) {
		cpe["api_pwd"], err = aes.EncryptToB64(apiPwd, m.Config.System.Aeskey)
		if err != nil {
			return err
		}
	}
	_, err = coll.InsertOne(context.TODO(), cpe)
	return err
}

func (m *CpeManager) UpdateCpeData(params web.RequestParams) error {
	data := params.GetParamMap("data")
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
	}else{
		delete(data, "api_pwd")
	}
	query := bson.M{"_id": _id}
	update := bson.M{"$set": data}
	_, err = m.GetTeamsAcsCollection(TeamsacsCpe).UpdateOne(context.TODO(), query, update)
	return err
}

// UpdateVpeBySn
func (m *CpeManager) UpdateCpeBySn(sn string, valmap map[string]interface{}) error {
	coll := m.GetTeamsAcsCollection(TeamsacsCpe)
	update := bson.M{"$set": valmap}
	_, err := coll.UpdateOne(context.TODO(), bson.M{"sn": sn}, update)
	return err
}
