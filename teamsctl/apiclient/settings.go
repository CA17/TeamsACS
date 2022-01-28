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
 */

package apiclient

import (
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
	"github.com/guonaihong/gout"
)

func FindSettings(ctype string) ([]models.SysConfig, error) {
	var resp []models.SysConfig
	var url = common.UrlJoin2(api.Apiurl, "/settings/list")
	err := gout.
		GET(url).
		SetHeader(api.CreateAuthorization()).
		Debug(api.Debug).
		SetTimeout(time.Second * 5).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func RemoveSettings(ctype, name string) (*web.WebRestResult, error) {
	var resp web.WebRestResult
	var url = common.UrlJoin2(api.Apiurl, "/settings/remove")
	err := gout.
		GET(url).
		SetHeader(api.CreateAuthorization()).
		Debug(api.Debug).
		SetTimeout(time.Second * 5).
		SetQuery(gout.H{"type": ctype, "name": name}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func UpdateSettings(cfgs ...models.SysConfig) (*web.WebRestResult, error) {
	var resp web.WebRestResult
	err := gout.
		POST(common.UrlJoin(api.Apiurl, "/settings/update")).
		Debug(api.Debug).
		SetHeader(api.CreateAuthorization()).
		SetTimeout(time.Second * 5).
		SetJSON(&cfgs).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
