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

package azureblob

import (
	"os"
	"testing"
)

func TestAzureBlob_UploadFile(t *testing.T) {
	ab := NewAzureBlob(os.Getenv("TEAMSACS_AZURE_STORAGE_ACCOUNT"), os.Getenv("TEAMSACS_AZURE_STORAGE_ACCESS_KEY"))
	r, err := ab.UploadFile("test", "","/Users/wangjuntao/github/goproject/src/TeamsACS/common/azureblob/azureblob.go")
	t.Log(r.Response())
	t.Log(err)
}

func TestAzureBlob_UploadFile2(t *testing.T) {
	ab := NewAzureBlob(os.Getenv("TEAMSACS_AZURE_STORAGE_ACCOUNT"), os.Getenv("TEAMSACS_AZURE_STORAGE_ACCESS_KEY"))
	r, err := ab.UploadFile("test", "20201101/test.go","/Users/wangjuntao/github/goproject/src/TeamsACS/common/azureblob/azureblob.go")
	t.Log(r.Response())
	t.Log(err)
}

func TestAzureBlob_ListFiles(t *testing.T) {
	ab := NewAzureBlob(os.Getenv("TEAMSACS_AZURE_STORAGE_ACCOUNT"), os.Getenv("TEAMSACS_AZURE_STORAGE_ACCESS_KEY"))
	r , err := ab.ListFiles("test")
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range r {
		t.Log(item)
	}
}
