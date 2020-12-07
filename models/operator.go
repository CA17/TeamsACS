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
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/constant"
)

// Operator
type Operator struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string `bson:"email,omitempty" json:"email,omitempty"`
	Username  string `bson:"username,omitempty" json:"username,omitempty"`
	Level     string `bson:"level" json:"level,omitempty"`
	ApiSecret string `bson:"api_secret,omitempty" json:"api_secret,omitempty"`
	Status    string `bson:"status,omitempty" json:"status,omitempty"`
	Remark    string `bson:"remark,omitempty" json:"remark,omitempty"`
}

func (a *Operator) AddValidate() error {
	switch {
	case common.IsEmptyOrNA(a.Username):
		return fmt.Errorf("invalid username")
	case common.IsEmptyOrNA(a.Level):
		return fmt.Errorf("invalid level")
	}
	return nil
}

// OperatorManager
type OperatorManager struct{ *ModelManager }

func (m *ModelManager) GetOpsManager() *OperatorManager {
	store, _ := m.ManagerMap.Get("OperatorManager")
	return store.(*OperatorManager)
}

// ExistOperator
func (m *OperatorManager) ExistOperator(username string) bool {
	coll := m.GetTeamsAcsCollection(TeamsacsOperator)
	count, _ := coll.CountDocuments(context.TODO(), bson.M{"username": username})
	return count > 0
}

// QueryOperators
func (m *OperatorManager) QueryOperators(params web.RequestParams) (*web.PageResult, error) {
	return m.QueryPagerItems(params, TeamsacsOperator)
}

// GetOperator
func (m *OperatorManager) GetOperator(username string) (*Operator, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsOperator)
	doc := coll.FindOne(context.TODO(), bson.M{"username": username})
	err := doc.Err()
	if err != nil {
		return nil, err
	}
	var result = new(Operator)
	err = doc.Decode(result)
	return result, err
}

// UpdateApiSecret
func (m *OperatorManager) UpdateApiSecret(username string) (string, error) {
	coll := m.GetTeamsAcsCollection(TeamsacsOperator)
	apisecret := common.UUID()
	_, err := coll.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$set": bson.M{"api_secret": apisecret}})
	return apisecret, err
}

// UpdateOperator
// update by username
func (m *OperatorManager) UpdateOperator(operator *Operator) error {
	coll := m.GetTeamsAcsCollection(TeamsacsOperator)
	query := bson.M{"username": operator.Username}
	data := bson.M{}
	if common.InSlice(operator.Level, []string{constant.NBIAdminLevel, constant.NBIOprLevel}) {
		data["level"] = operator.Level
	}
	if common.InSlice(operator.Status, []string{constant.ENABLED, constant.DISABLED}) {
		data["status"] = operator.Status
	}
	data["email"] = operator.Email
	data["remark"] = operator.Remark

	update := bson.M{"$set": data}
	_, err := coll.UpdateOne(context.TODO(), query, update)
	return err
}

// AddOperator
func (m *OperatorManager) AddOperator(operator *Operator) (string, error) {
	if err := operator.AddValidate(); err != nil {
		return "", err
	}
	if m.ExistOperator(operator.Username) {
		return "", fmt.Errorf("operator exists")
	}
	if common.IsEmptyOrNA(operator.ApiSecret) {
		operator.ApiSecret = common.UUID()
	}
	operator.Status = constant.ENABLED
	r, err := m.GetTeamsAcsCollection(TeamsacsOperator).InsertOne(context.TODO(), operator)
	if err != nil {
		return "", err
	}
	return r.InsertedID.(string), err
}

// initSuper
func (m *OperatorManager) InitSuper(username string) (string, error) {
    secret := common.UUID()
	if !m.ExistOperator(username) {
		sopr := new(Operator)
		sopr.Username = username
		sopr.ApiSecret = secret
		sopr.Level = ""
		sopr.Status = constant.ENABLED
		sopr.Level = constant.NBIAdminLevel
		sopr.Email = ""
		sopr.Remark = "init"
		_, err := m.AddOperator(sopr)
		if err != nil {
			return "", err
		}
	} else {
		var err error
		secret, err = m.UpdateApiSecret(username)
		if err != nil {
			return "", err
		}
	}
	return secret, nil
}

// DeleteSubscribe
func (m *OperatorManager) DeleteOperator(username string) error {
	if common.IsEmptyOrNA(username) {
		return fmt.Errorf("username is empty or NA")
	}
	_, err := m.GetTeamsAcsCollection(TeamsacsOperator).DeleteOne(context.TODO(), bson.M{"username": username})
	return err
}
