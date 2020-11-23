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

package maputils

import (
	"fmt"
	"time"

	"github.com/ca17/teamsacs/common"
)

func GetStringValue(d map[string]interface{}, key string, defval string) string {
	val, ok := d[key]
	if !ok || val == nil || val == "" {
		return defval
	}
	val2, err := common.ParseString(val)
	if err != nil {
		return defval
	}
	return val2
}

func GetStringValueWithErr(d map[string]interface{}, key string) (string, error) {
	val, ok := d[key]
	if !ok || val == nil || val == "" {
		return "", fmt.Errorf("%s is empty", key)
	}
	val2, err := common.ParseString(val)
	if err != nil {
		return "", err
	}
	return val2, nil
}

func GetIntValue(d map[string]interface{}, key string, defval int) int {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return defval
		}
		return int(v)
	}
	return defval
}

func GetInt64Value(d map[string]interface{}, key string, defval int64) int64 {
	val, ok := d[key]
	if ok {
		v, err := common.ParseInt64(val)
		if err != nil {
			return defval
		}
		return v
	}
	return defval
}


func GetFloat64Value(d map[string]interface{}, key string, defval float64) float64 {
	val, ok := d[key]
	if ok {
		v, err := common.ParseFloat64(val)
		if err != nil {
			return defval
		}
		return v
	}
	return defval
}


func GetDateObject(d map[string]interface{}, key string, defval time.Time) time.Time {
	val, ok := d[key]
	if ok {
		var result = defval
		val, err := common.ParseString(val)
		if err != nil {
			return defval
		}
		if len(val) == 19 {
			result, err = time.Parse("2006-01-02 15:04:05", val)
		}else{
			result, err = time.Parse("2006-01-02 15:04:05 Z0700 MST", val)
		}
		if err != nil {
			return defval
		}
		return result
	}
	return defval
}


