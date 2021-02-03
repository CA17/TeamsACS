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

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/web"
)

type Syslog struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty"`
	Logtype   string     `bson:"logtype,omitempty" json:"logtype,omitempty"`
	Attrs     Attributes `bson:"attrs" json:"attrs,omitempty" `
	Timestamp time.Time  `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (m *OperatorManager) AddSyslog(item *Syslog) {
	coll := m.GetTeamsAcsCollection(TeamsacsSyslog)
	_, err := coll.InsertOne(context.TODO(), item)
	if err != nil {
		log.Error(err)
	}
}

func (m *OperatorManager) QuerySyslog(params web.RequestParams) (*web.PageResult, error) {
	var findOptions = options.Find()
	var pos = params.GetInt64WithDefval("start", 0)
	findOptions.SetSkip(pos)
	findOptions.SetLimit(params.GetInt64WithDefval("count", 40))
	coll := m.GetTeamsAcsCollection(TeamsacsSyslog)
	q := processQueryParams(params, findOptions)
	cur, err := coll.Find(context.TODO(), q, findOptions)
	if err != nil {
		return nil, err
	}
	var countOptions = options.Count()
	total, err := coll.CountDocuments(context.TODO(), q, countOptions)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0)
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Error(err)
		} else {
			items = append(items, elem)
		}
	}
	return &web.PageResult{TotalCount: total, Pos: pos, Data: items}, nil
}
