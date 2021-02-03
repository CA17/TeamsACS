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
	"strconv"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ca17/teamsacs/common/web"
)

// Config
type Config struct {
	ID    string `bson:"_id,omitempty" json:"id,omitempty"`
	Type  string `bson:"type" json:"type,omitempty"`
	Name  string `bson:"name" json:"name,omitempty"`
	Value string `bson:"value" json:"value,omitempty"`
}

type ConfigManager struct{ *ModelManager }

func (m *ModelManager) GetConfigManager() *ConfigManager {
	store, _ := m.ManagerMap.Get("ConfigManager")
	return store.(*ConfigManager)
}

func (m *ConfigManager) QueryConfig(params web.RequestParams) (*web.QueryResult, error) {
	return m.QueryItems(params, TeamsacsConfig)
}

func (m *ConfigManager) GetConfigValue(ctype, name string) string {
	coll := m.GetTeamsAcsCollection(TeamsacsConfig)
	doc := coll.FindOne(context.TODO(), bson.M{"type": ctype, "name": name})
	err := doc.Err()
	if err != nil {
		return ""
	}
	var result = new(Config)
	err = doc.Decode(&result)
	return result.Value
}

func (m *ConfigManager) GetRadiusConfigValue(name string) string {
	coll := m.GetTeamsAcsCollection(TeamsacsConfig)
	doc := coll.FindOne(context.TODO(), bson.M{"type": "radius", "name": name})
	err := doc.Err()
	if err != nil {
		return ""
	}
	var result = ""
	err = doc.Decode(&result)
	return result
}

func (m *ConfigManager) GetRadiusConfigStringValue(name string, defval string) string {
	val := m.GetRadiusConfigValue(name)
	if val == "" {
		return defval
	}
	return val
}

func (m *ConfigManager) GetRadiusConfigIntValue(name string, defval int64) int64 {
	val := m.GetRadiusConfigValue(name)
	if val == "" {
		return defval
	}
	v, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defval
	}
	return v
}

func (m *ConfigManager) UpdateConfigValue(ctype, name, value string) error {
	coll := m.GetTeamsAcsCollection(TeamsacsConfig)
	query := bson.M{"type": ctype, "name": name}
	update := bson.M{"$set": bson.M{"value": value}}
	_, err := coll.UpdateOne(context.TODO(), query, update)
	return err
}
