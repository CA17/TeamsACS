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
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/web"
)

// A generic data CRUD management API with no predefined schema,

// DataManager
type DataManager struct{ *ModelManager }

func (m *ModelManager) GetDataManager() *DataManager {
	store, _ := m.ManagerMap.Get("DataManager")
	return store.(*DataManager)
}

// GetDataById
func (m *DataManager) GetData(params web.RequestParams) (*Attributes, error) {
	_id := params.GetParamMap("querymap").GetMustString("_id")
	collname := params.GetMustString("collname")
	coll := m.GetTeamsAcsCollection(collname)
	doc := coll.FindOne(context.TODO(), bson.M{"_id": _id})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Attributes)
	err = doc.Decode(result)
	return result, err
}

// GetDataById
func (m *DataManager) GetDataNameValues(params web.RequestParams) ([]NameValue, error) {
	_id := params.GetParamMap("querymap").GetMustString("_id")
	collname := params.GetMustString("collname")
	coll := m.GetTeamsAcsCollection(collname)
	doc := coll.FindOne(context.TODO(), bson.M{"_id": _id})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var nvs = make([]NameValue, 0)
	var result = new(Attributes)
	err = doc.Decode(result)
	for name, value := range *result {
		nvs = append(nvs, NameValue{
			Name:  name,
			Value: value,
		})
	}
	return nvs, err
}

// AddData
func (m *DataManager) AddData(params web.RequestParams) error {
	data := params.GetParamMap("data")
	_id := data.GetString("_id")
	if common.IsEmptyOrNA(_id) {
		data["_id"] = common.UUID()
	}
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05 Z0700 MST")
	coll := m.GetTeamsAcsCollection(params.GetMustString("collname"))
	_, err := coll.InsertOne(context.TODO(), data)
	return err
}

// AddData
func (m *DataManager) AddBatchData(collname string, datas []interface{}) error {
	coll := m.GetTeamsAcsCollection(collname)
	_, err := coll.InsertMany(context.TODO(), datas)
	return err
}

// UpdateData
func (m *DataManager) UpdateData(params web.RequestParams) error {
	data := params.GetParamMap("data")
	data["update_time"] = time.Now().Format("2006-01-02 15:04:05 Z0700 MST")
	_id := data.GetMustString("_id")
	query := bson.M{"_id": _id}
	update := bson.M{"$set": data}
	_, err := m.GetTeamsAcsCollection(params.GetMustString("collname")).UpdateOne(context.TODO(), query, update)
	return err
}

// DeleteData
func (m *DataManager) DeleteData(params web.RequestParams) error {
	collname := params.GetMustString("collname")
	ids := params.GetParamMap("querymap").GetString("ids")
	if ids != "" {
		idarray := bson.A{}
		for _, id := range strings.Split(ids, ",") {
			idarray = append(idarray, id)
		}
		filter := bson.M{"_id": bson.M{"$in": idarray}}
		_, err := m.GetTeamsAcsCollection(collname).DeleteMany(context.TODO(), filter)
		return err
	}else{
		_id := params.GetParamMap("data").GetString("_id")
		_, err := m.GetTeamsAcsCollection(collname).DeleteMany(context.TODO(), bson.M{"_id": _id})
		return err
	}
}

// SaveData
func (m *DataManager) SaveData(params web.RequestParams) (interface{}, error) {
	op := params.GetString("webix_operation")
	switch op {
	case "insert":
		params["_id"] = common.UUID()
		err := m.AddData(params)
		return map[string]interface{}{"id": params["_id"]}, err
	case "update":
		err := m.UpdateData(params)
		return map[string]interface{}{"status": "updated"}, err
	case "delete":
		err := m.DeleteData(params)
		return make(map[string]interface{}), err
	}
	return nil, errors.New("Unsupported Operations")
}
