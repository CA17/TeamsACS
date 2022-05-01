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

func FindCwmpFile(query string) ([]models.CwmpFile, error) {
	var resp []models.CwmpFile
	var url = common.UrlJoin2(api.Apiurl, "/cwmpfile/list")
	err := gout.
		GET(url).
		SetHeader(api.CreateAuthorization()).
		Debug(api.Debug).
		SetTimeout(time.Second * 5).
		SetQuery(gout.H{"query": query}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func FindCwmpFileTask(id, status string) ([]models.CwmpFile, error) {
	var resp []models.CwmpFile
	var url = common.UrlJoin2(api.Apiurl, "/cwmpfile/task/list")
	err := gout.
		GET(url).
		SetHeader(api.CreateAuthorization()).
		Debug(api.Debug).
		SetTimeout(time.Second * 5).
		SetQuery(gout.H{"id": id, "status": status}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func UploadCwmpFile(file models.CwmpFile) (*web.WebRestResult, error) {
	var hd = api.CreateAuthorization()
	hd["product_class"] = file.ProductClass
	hd["oui"] = file.Oui
	hd["manufacturer"] = file.Manufacturer
	hd["version"] = file.Version
	hd["file_type"] = file.FileType
	var resp web.WebRestResult
	err := gout.
		POST(common.UrlJoin(api.Apiurl, "/cwmpfile/upload")).
		Debug(api.Debug).
		SetHeader(hd).
		SetTimeout(time.Second * 5).
		SetForm(gout.H{"upload": gout.FormFile(file.FilePath)}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func RemoveCwmpFile(id string) (*web.WebRestResult, error) {
	var resp web.WebRestResult
	var url = common.UrlJoin2(api.Apiurl, "/cwmpfile/remove")
	err := gout.
		DELETE(url).
		SetHeader(api.CreateAuthorization()).
		Debug(api.Debug).
		SetTimeout(time.Second * 5).
		SetQuery(gout.H{"id": id}).
		BindJSON(&resp).
		Do()
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
