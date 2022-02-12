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
	"testing"

	"github.com/ca17/teamsacs/models"
)

func TestFindCwmpFile(t *testing.T) {
	api.Debug = true
	r, err := FindCwmpFile("")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestFindCwmpFileTask(t *testing.T) {
	api.Debug = true
	r, err := FindCwmpFileTask("", "error")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestRemoveCwmpFile(t *testing.T) {
	api.Debug = true
	r, err := RemoveCwmpFile("1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestUploadCwmpFile(t *testing.T) {
	api.Debug = true
	r, err := UploadCwmpFile(models.CwmpFile{
		FileType: "3 Vendor Configuration File",
		FilePath: "/Users/wangjuntao/github/TeamsACS/teamsctl/apiclient/files_test.go",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}
