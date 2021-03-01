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

package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
)

type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// WEB 参数
type WebForm struct {
	FormItem interface{}
	Posts    url.Values        `json:"-" form:"-" query:"-"`
	Gets     url.Values        `json:"-" form:"-" query:"-"`
	Params   map[string]string `json:"-" form:"-" query:"-"`
}

func EmptyWebForm() *WebForm {
	v := &WebForm{}
	v.Params = make(map[string]string, 0)
	v.Posts = make(url.Values, 0)
	v.Gets = make(url.Values, 0)
	return v
}

func NewWebForm(c echo.Context) *WebForm {
	v := &WebForm{}
	v.Params = make(map[string]string)
	v.Posts, _ = c.FormParams()
	v.Gets = c.QueryParams()
	for _, p := range c.ParamNames() {
		v.Params[p] = c.Param(p)
	}
	return v
}

func (f *WebForm) Set(name string, value string) {
	f.Gets.Set(name, value)
}

func (f *WebForm) Param(name string) string {
	return f.Param(name)
}

func (f *WebForm) Param2(name string, defval string) string {
	if val, ok := f.Params[name]; ok {
		return val
	}
	return defval
}

func (f *WebForm) GetDateRange(name string) (DateRange, error) {
	var dr = DateRange{Start: "", End: ""}
	val := f.GetVal(name)
	if val == "" {
		return dr, nil
	}
	err := json.Unmarshal([]byte(val), &dr)
	if err != nil {
		return dr, err
	}
	return dr, nil
}

func (f *WebForm) GetVal(name string) string {
	val := f.Posts.Get(name)
	if val != "" {
		return val
	}
	val = f.Gets.Get(name)
	if val != "" {
		return val
	}
	return ""
}

func (f *WebForm) GetMustVal(name string) (string, error) {
	val := f.Posts.Get(name)
	if val != "" {
		return val, nil
	}
	val = f.Gets.Get(name)
	if val != "" {
		return val, nil
	}
	return "", errors.New(name + " 不能为空")
}

func (f *WebForm) GetVal2(name string, defval string) string {
	val := f.Posts.Get(name)
	if val != "" {
		return val
	}
	val = f.Gets.Get(name)
	if val != "" {
		return val
	}
	return defval
}

func (f *WebForm) GetIntVal(name string, defval int) int {
	val := f.GetVal(name)
	if val == "" {
		return defval
	}
	v, _ := strconv.Atoi(val)
	return v
}

func (f *WebForm) GetInt64Val(name string, defval int64) int64 {
	val := f.GetVal(name)
	if val == "" {
		return defval
	}
	v, _ := strconv.ParseInt(val, 10, 64)
	return v
}

type PageResult struct {
	TotalCount int64       `json:"total_count,omitempty"`
	Pos        int64       `json:"pos"`
	Data       interface{} `json:"data"`
}

var EmptyPageResult = &PageResult{
	TotalCount: 0,
	Pos:        0,
	Data:       common.EmptyList,
}

type JsonOptions struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type QueryResult []map[string]interface{}

// 通用查询参数
type RequestParams map[string]interface{}

var EmptyRequestParams = func() RequestParams {
	return make(RequestParams)
}

func (jp RequestParams) Set(key string, value interface{}) {
	jp[key] = value
}

// 获取单个字符串
func (jp RequestParams) GetString(key string) string {
	v, ok := jp[key]
	if !ok {
		return ""
	}
	vv, ok := v.(string)
	if ok {
		return vv
	}
	return ""
}

// 获取子查询
func (jp RequestParams) GetParamMap(key string) RequestParams {
	v, ok := jp[key]
	if !ok {
		return EmptyRequestParams()
	}
	vv, ok := v.(map[string]interface{})
	if ok {
		return common.DeepCopy(vv).(map[string]interface{})
	}
	return EmptyRequestParams()
}

// CopyIMap copy to interface map
func (jp RequestParams) CopyIMap() (map[string]interface{}, error) {
	newMap := make(map[string]interface{})
	bs, err := common.JsonMarshal(jp)
	if err != nil {
		return nil, err
	}
	err = common.JsonUnmarshal(bs, &newMap)
	if err != nil {
		return nil, err
	}
	return newMap, nil
}

// GetStringWithDefval ..
func (jp RequestParams) GetStringWithDefval(key string, defval string) string {
	v, ok := jp[key]
	if !ok {
		return defval
	}
	vv, ok := v.(string)
	if ok {
		return vv
	}
	return defval
}

func (jp RequestParams) GetInt64(key string) int64 {
	v, ok := jp[key]
	if !ok {
		return 0
	}
	switch v.(type) {
	case int64:
		return v.(int64)
	case string:
		vvv, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return 0
		}
		return vvv
	}
	return 0
}

func (jp RequestParams) GetInt64WithDefval(key string, devfval int64) int64 {
	v, ok := jp[key]
	if !ok {
		return devfval
	}
	switch v.(type) {
	case int64:
		return v.(int64)
	case string:
		vvv, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return devfval
		}
		return vvv
	}
	return devfval
}

var NoValueError = fmt.Errorf("no value")

func (jp RequestParams) GetMustString(key string) string {
	v, ok := jp[key]
	if !ok {
		common.Must(fmt.Errorf(key + " attr no value"))
	}
	vv, ok := v.(string)
	if !ok {
		common.Must(fmt.Errorf(key + " attr not string"))
	}
	return vv
}

func (jp RequestParams) GetMustInt64(key string) int64 {
	v, ok := jp[key]
	if !ok {
		common.Must(fmt.Errorf(key + " attr not int64"))
	}
	vv, ok := v.(int64)
	if !ok {
		common.Must(NoValueError)
	}
	return vv
}

func (jp RequestParams) GetQueryMap() RequestParams {
	return jp.GetParamMap("querymap")
}

func (jp RequestParams) GetFilterMap() RequestParams {
	return jp.GetParamMap("filtermap")
}

func (jp RequestParams) GetFilterInMap() RequestParams {
	return jp.GetParamMap("filterinmap")
}

func (jp RequestParams) GetEqualMap() RequestParams {
	return jp.GetParamMap("equalmap")
}

func (jp RequestParams) GetTimeRangeMap() RequestParams {
	return jp.GetParamMap("timerangemap")
}

func (jp RequestParams) GetSortMap() RequestParams {
	return jp.GetParamMap("sortmap")
}

func (jp RequestParams) ToJson() string {
	bs, _ := json.MarshalIndent(jp, "", "\t")
	return string(bs)
}
