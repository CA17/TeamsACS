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

package mikrotik

import (
	"encoding/json"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const _test_jsonstr = `{
	"AccessPoint": {
      "1": {
        "_object": true,
        "_writable": true,
        "APN": {
          "_object": false,
          "_writable": true
        },
        "Password": {
          "_object": false,
          "_writable": true
        },
        "Username": {
          "_object": false,
          "_writable": true
        },
        "_timestamp": {"$date": "2020-09-12T14:49:48.260Z"}
      },
      "_object": true,
      "_timestamp": {"$date": "2020-09-12T14:49:48.260Z"},
      "_writable": true
    }
}`

func TestGetObject(t *testing.T) {
	var tmap = make(map[string]interface{})
	err := json.Unmarshal([]byte(_test_jsonstr), &tmap)
	if err != nil {
		t.Fatal(err)
	}
	val := GetObject(tmap, "AccessPoint.1.Username._writable")
	t.Log(val)
	if val != true {
		t.Fatal()
	}
	val2 := GetObject(tmap, "AccessPoint.1._writable")
	t.Log(val2)
	if val2 != true {
		t.Fatal()
	}
	val3 := GetObject(tmap, "AccessPoint.1.Password._writable")
	t.Log(val3)
	if val3 != true {
		t.Fatal()
	}
}

func Benchmark(b *testing.B) {
	var tmap = make(map[string]interface{})
	err := json.Unmarshal([]byte(_test_jsonstr), &tmap)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		GetObject(tmap, "AccessPoint.1.Username._writable")
	}
}

func TestTimestamp(t *testing.T) {
	tt := primitive.NewDateTimeFromTime(time.Now())
	jsons, err := tt.MarshalJSON()
	t.Log(string(jsons), err)
	val := map[string]interface{}{
		"_timestamp": tt,
	}
	vv := val["_timestamp"].(primitive.DateTime)
	t.Log(vv)
}
