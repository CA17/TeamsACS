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
	"encoding/json"
	"testing"

	"github.com/ca17/teamsacs/common"
)

func TestOperator2json(t *testing.T) {
	item := &Operator{
		ID:       common.UUID(),
		Email:    "test@teamsacs.com",
		Username: "opr",
		Level:    "opr",
		Remark:   "opr",
	}
	bs, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))

}

func Benchmark(b *testing.B) {
	vpe := make(Vpe)
	vpe["secret"] = "aaa"
	for i := 0; i < b.N; i++ {
		vpe.GetSecret()
	}
}
