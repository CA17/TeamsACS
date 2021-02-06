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

// Mongodb general query method

package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/web"
)

type QueryForm struct {
	Sort string `query:"sort" form:"sort"`
	Size int64  `query:"size" form:"size"`
}

type SubscribeQueryForm struct {
	QueryForm
	Username string `query:"username" form:"username"`
}

func logQueryParams(q interface{}) {
	bs, _ := json.MarshalIndent(q, "", "\t")
	log.Info(string(bs))
}

// Processing query parameters
func processQueryParams(params web.RequestParams, findOptions *options.FindOptions) bson.D {
	var q = bson.D{}
	for qname, val := range params.GetParamMap("filtermap") {
		filter := bson.D{{"$regex", primitive.Regex{Pattern: val.(string), Options: "i"}}}
		q = append(q, bson.E{Key: qname, Value: filter})
	}
	for qname, val := range params.GetParamMap("equalmap") {
		q = append(q, bson.E{Key: qname, Value: val})
	}
	for qname, val := range params.GetParamMap("filterinmap") {
		var varray = bson.A{}
		for _, v := range strings.Split(val.(string), ",") {
			varray = append(varray, v)
		}
		q = append(q, bson.E{Key: qname, Value: bson.M{"$in": varray}})
	}
	for sname, sval := range params.GetParamMap("sortmap") {
		if sval == "asc" {
			findOptions.SetSort(bson.D{{sname, 1}})
		} else if sval == "desc" {
			findOptions.SetSort(bson.D{{sname, -1}})
		}
	}

	timerangemap := params.GetParamMap("timerangemap")
	_start := timerangemap["start"]
	startValue := timerangemap["start_value"]
	_end := timerangemap["end"]
	endValue := timerangemap["end_value"]
	switch {
	case _start != nil && _end != nil && startValue != nil && endValue != nil:
		timefilterVal := bson.M{"$gte": startValue, "$lte": endValue}
		timefilter := bson.E{Key: _start.(string), Value: timefilterVal}
		q = append(q, timefilter)
	case _start != nil && startValue != nil && (_end == nil || endValue == nil):
		timefilterVal := bson.M{"$gte": startValue}
		timefilter := bson.E{Key: _start.(string), Value: timefilterVal}
		q = append(q, timefilter)
	case (startValue == nil || _start == nil) && endValue != nil && _end != nil:
		timefilterVal := bson.M{"$lte": endValue}
		timefilter := bson.E{Key: _end.(string), Value: timefilterVal}
		q = append(q, timefilter)
	}

	return q
}

// General query, no paging, suitable for small data sets
func (m *ModelManager) QueryItems(params web.RequestParams, collatiion string) (*web.QueryResult, error) {
	var findOptions = options.Find()
	limit := params.GetInt64WithDefval("limit", 0)
	if limit > 0 {
		findOptions.SetLimit(limit)
	}
	coll := m.GetTeamsAcsCollection(collatiion)
	q := processQueryParams(params, findOptions)
	logQueryParams(q)
	cur, err := coll.Find(context.TODO(), q, findOptions)
	if err != nil {
		return nil, err
	}
	items := make(web.QueryResult, 0)
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Error(err)
		} else {
			items = append(items, elem)
		}
	}
	return &items, nil
}

func (m *ModelManager) QueryItemOptions(params web.RequestParams, collatiion string) ([]web.JsonOptions, error) {
	jsonoptions := make([]web.JsonOptions, 0)
	var findOptions = options.Find()
	coll := m.GetTeamsAcsCollection(collatiion)
	querymap := params.GetParamMap("querymap")
	optionName := querymap.GetString("optionname")
	optionId := querymap.GetString("optionid")
	if optionName == "" {
		return jsonoptions, fmt.Errorf("option name is empty")
	}
	if optionId == "" {
		optionId = "_id"
	}
	q := processQueryParams(params, findOptions)
	logQueryParams(q)
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
			optionId := elem[optionId].(string)
			optionValue := elem[optionName]
			if optionValue == nil || optionValue == "" {
				continue
			}
			jsonoptions = append(jsonoptions, web.JsonOptions{
				Id:    optionId,
				Value: optionValue.(string),
			})
		}
	}
	return jsonoptions, nil
}

// Paging query, suitable for large data sets
func (m *ModelManager) QueryPagerItems(params web.RequestParams, collatiion string) (*web.PageResult, error) {
	var findOptions = options.Find()
	var pos = params.GetInt64WithDefval("start", 0)
	findOptions.SetSkip(pos)
	findOptions.SetLimit(params.GetInt64WithDefval("count", 40))
	coll := m.GetTeamsAcsCollection(collatiion)
	q := processQueryParams(params, findOptions)
	logQueryParams(q)
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
