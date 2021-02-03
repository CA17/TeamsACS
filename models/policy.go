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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/web"
)

type Policy map[string]interface{}

func (p Policy) GetPid() (string, error) {
	return maputils.GetStringValueWithErr(p, "pid")
}

func (p Policy) GetName() (string, error) {
	return maputils.GetStringValueWithErr(p, "name")
}

func (p Policy) GetClassify() (string, error) {
	return maputils.GetStringValueWithErr(p, "classify")
}

func (p Policy) GetTr069ParamData() (string, error) {
	return maputils.GetStringValueWithErr(p, "data")
}

func (p Policy) GetApiCommand() (string, error) {
	return maputils.GetStringValueWithErr(p, "command")
}

func (p Policy) GetApiParams() (string, error) {
	return maputils.GetStringValueWithErr(p, "params")
}

func (p Policy) GetApiProps() (string, error) {
	return maputils.GetStringValueWithErr(p, "props")
}

func (p Policy) GetScriptContent() (string, error) {
	return maputils.GetStringValueWithErr(p, "content")
}

type PolicyVariable map[string]interface{}

func (p PolicyVariable) GetName() (string, error) {
	return maputils.GetStringValueWithErr(p, "name")
}

func (p PolicyVariable) GetValue() (string, error) {
	return maputils.GetStringValueWithErr(p, "value")
}

func (p PolicyVariable) GetClassify() (string, error) {
	return maputils.GetStringValueWithErr(p, "classify")
}

type PolicyManager struct{ *ModelManager }

func (m *ModelManager) GetPolicyManager() *PolicyManager {
	store, _ := m.ManagerMap.Get("PolicyManager")
	return store.(*PolicyManager)
}

// GetMikrotikApiPolicyByPid
func (m *PolicyManager) GetMikrotikApiPolicyByPid(pid string) (*Policy, error) {
	return m.GetPolicyByPid(TeamsacsMikrotikApiPolicy, pid)
}

// GetMikrotikScriptPolicyByPid
func (m *PolicyManager) GetMikrotikScriptPolicyByPid(pid string) (*Policy, error) {
	return m.GetPolicyByPid(TeamsacsMikrotikScriptPolicy, pid)
}

// GetTr069PolicyByPid
func (m *PolicyManager) GetTr069PolicyByPid(pid string) (*Policy, error) {
	return m.GetPolicyByPid(TeamsacsTr069Policy, pid)
}

// GetPolicyByPid
func (m *PolicyManager) GetPolicyByPid(collname, pid string) (*Policy, error) {
	coll := m.GetTeamsAcsCollection(collname)
	doc := coll.FindOne(context.TODO(), bson.M{"pid": pid})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Policy)
	err = doc.Decode(result)
	return result, err
}

// QueryAllVars
func (m *PolicyManager) QueryAllVars() ([]PolicyVariable, error) {
	var findOptions = options.Find()
	coll := m.GetTeamsAcsCollection(TeamsacsPolicyVariable)
	cur, err := coll.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	items := make([]PolicyVariable, 0)
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Error(err)
		} else {
			items = append(items, elem)
		}
	}
	return items, nil
}

// QueryMikrotikApiPolicyOptions
func (m *PolicyManager) QueryMikrotikApiPolicyOptions() ([]web.JsonOptions, error) {
	return m.QueryPolicyOptions(TeamsacsMikrotikApiPolicy)
}

// QueryMikrotikScriptPolicyOptions
func (m *PolicyManager) QueryMikrotikScriptPolicyOptions() ([]web.JsonOptions, error) {
	return m.QueryPolicyOptions(TeamsacsMikrotikScriptPolicy)
}

// QueryTr069PolicyOptions
func (m *PolicyManager) QueryTr069PolicyOptions() ([]web.JsonOptions, error) {
	return m.QueryPolicyOptions(TeamsacsTr069Policy)
}

// QueryPolicyOptions
func (m *PolicyManager) QueryPolicyOptions(collatiion string) ([]web.JsonOptions, error) {
	jsonoptions := make([]web.JsonOptions, 0)
	var findOptions = options.Find()
	coll := m.GetTeamsAcsCollection(collatiion)
	var q = bson.D{}
	cur, err := coll.Find(context.TODO(), q, findOptions)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Error(err)
		} else {
			optionId := elem["pid"].(string)
			optionValue := elem["name"].(string)
			if optionValue == "" {
				continue
			}
			jsonoptions = append(jsonoptions, web.JsonOptions{
				Id:    optionId,
				Value: optionValue,
			})
		}
	}
	return jsonoptions, nil
}
