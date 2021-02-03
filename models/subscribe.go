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
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/constant"
)

type Subscribe map[string]interface{}

// AuthorizationProfile Method
func (a Subscribe) GetExpireTime() time.Time {
	return maputils.GetDateObject(a, "expire_time", time.Now().Add(time.Second*60))
}

func (a Subscribe) GetInterimInterval() int {
	return maputils.GetIntValue(a, "interim_interval", 120)
}

func (a Subscribe) GetAddrPool() string {
	return maputils.GetStringValue(a, "addr_pool", constant.NA)
}

func (a Subscribe) GetIpaddr() string {
	return maputils.GetStringValue(a, "ipaddr", constant.NA)
}

func (a Subscribe) GetUpRateKbps() int {
	return maputils.GetIntValue(a, "up_rate", 0)
}

func (a Subscribe) GetDownRateKbps() int {
	return maputils.GetIntValue(a, "down_rate", 0)
}

func (a Subscribe) GetDomain() string {
	return maputils.GetStringValue(a, "domain", constant.NA)
}

func (a Subscribe) GetLimitPolicy() string {
	return maputils.GetStringValue(a, "limit_policy", constant.NA)
}

func (a Subscribe) GetUpLimitPolicy() string {
	return maputils.GetStringValue(a, "up_limit_policy", constant.NA)
}

func (a Subscribe) GetDownLimitPolicy() string {
	return maputils.GetStringValue(a, "down_limit_policy", constant.NA)
}

func (a Subscribe) GetMacAddr() string {
	return maputils.GetStringValue(a, "mac_addr", constant.NA)
}

func (a Subscribe) GetPassword() string {
	return maputils.GetStringValue(a, "password", constant.NA)
}

func (a Subscribe) GetUsername() string {
	return maputils.GetStringValue(a, "username", constant.NA)
}

func (a Subscribe) GetActiveNum() int {
	return maputils.GetIntValue(a, "active_num", 0)
}

func (a Subscribe) GetStatus() string {
	return maputils.GetStringValue(a, "status", constant.ENABLED)
}

func (a Subscribe) GetUserType() string {
	return maputils.GetStringValue(a, "user_type", constant.NA)
}

func (a Subscribe) GetBindVlan() int {
	return maputils.GetIntValue(a, "bind_vlan", 0)
}

func (a Subscribe) GetVlan1() int {
	return maputils.GetIntValue(a, "vlanid1", 0)
}

func (a Subscribe) GetVlan2() int {
	return maputils.GetIntValue(a, "vlanid2", 0)
}

// SubscribeManager
type SubscribeManager struct{ *ModelManager }

func (m *ModelManager) GetSubscribeManager() *SubscribeManager {
	store, _ := m.ManagerMap.Get("SubscribeManager")
	return store.(*SubscribeManager)
}

// GetSubscribeByAttr
func (m *SubscribeManager) GetSubscribeByAttr(name string, value interface{}) (*Subscribe, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsSubscribe)
	doc := coll.FindOne(context.TODO(), bson.M{name: value})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Subscribe)
	err = doc.Decode(result)
	return result, err
}

// QuerySubscribes
func (m *SubscribeManager) QuerySubscribes(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsSubscribe)
}

// GetSubscribeByUser
func (m *SubscribeManager) GetSubscribeByUser(username string) (*Subscribe, error) {
	return m.GetSubscribeByAttr("username", username)
}

// GetSubscribeByMac
func (m *SubscribeManager) GetSubscribeByMac(mac string) (*Subscribe, error) {
	return m.GetSubscribeByAttr("macaddr", mac)
}

// GetSubscribeById
func (m *SubscribeManager) GetSubscribeById(id string) (*Subscribe, error) {
	return m.GetSubscribeByAttr("_id", id)
}

// GetSubscribeByCpeid
func (m *SubscribeManager) GetSubscribeByCpeid(cpeid string) (*Subscribe, error) {
	filter := bson.D{{"$regex", primitive.Regex{Pattern: cpeid, Options: "i"}}}
	return m.GetSubscribeByAttr("cpe_ids", filter)
}

// UpdateSubscribeByUsername
func (m *SubscribeManager) UpdateSubscribeByUsername(username string, valmap map[string]interface{}) error {
	coll := m.GetTeamsAcsCollection(TeamsacsSubscribe)
	_, err := coll.UpdateOne(context.TODO(), bson.M{"username": username}, valmap)
	return err
}

// AddSubscribeData
// Create an account while the view is synchronized to Elastic
// If no Elastic configuration is available, the sync will be ignored
func (m *SubscribeManager) AddSubscribeData(params web.RequestParams) error {
	data := params.GetParamMap("data")
	_id := data.GetString("_id")
	if common.IsEmptyOrNA(_id) {
		data["_id"] = common.UUID()
	}
	go func() {
		_ = m.Elastic.AddData("teamsacs_subscribe", data)
	}()
	coll := m.GetTeamsAcsCollection(TeamsacsSubscribe)
	_, err := coll.InsertOne(context.TODO(), data)
	return err
}
