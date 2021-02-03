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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ca17/teamsacs/common/azureblob"
)

const (
	Version = "1.0.0"
)

var (
	h         = flag.Bool("h", false, "help usage")
	container = flag.String("c", "default", "Container Name")
	target    = flag.String("t", "", "target filepath")
	src       = flag.String("s", "", "src file")
	// get    = flag.Bool("get", false, "下载模块版本")
	// put    = flag.Bool("put", false, "上传模块版本")
	// del    = flag.Bool("del", false, "删除模块版本")
	// list   = flag.Bool("list", false, "列出所有版本")
)

func main() {

	flag.Parse()

	if *h == true {
		ustr := fmt.Sprintf("bssr version: abfs/%s, Usage:\nabfs -h\nOptions:", Version)
		fmt.Fprintf(os.Stderr, ustr)
		flag.PrintDefaults()
		return
	}

	if *container == "" {
		fmt.Fprintf(os.Stderr, "container is empty")
		return
	}

	if *target == "" {
		fmt.Fprintf(os.Stderr, "target path is empty")
		return
	}

	if *target == "" {
		fmt.Fprintf(os.Stderr, "target path is empty")
		return
	}

	ab := azureblob.NewAzureBlob(os.Getenv("TEAMSACS_AZURE_STORAGE_ACCOUNT"), os.Getenv("TEAMSACS_AZURE_STORAGE_ACCESS_KEY"))
	r, err := ab.UploadFile(*container, *target, *src)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	fmt.Fprintln(os.Stdout, "upload status "+r.Response().Status)

}
