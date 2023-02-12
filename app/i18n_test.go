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

package app

import (
	"testing"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/config"
)

func TestApplication_Translate(t *testing.T) {
	InitGlobalApplication(config.LoadConfig("../teamsacs.yml"))
	app.appConfig.System.Workdir = "/tmp"
	app.Translate(ZhCN, "global", "Create", "Create")
	app.Translate(ZhCN, "global", "Remove", "Remove")
	app.Translate(ZhCN, "global", "Edit", "编辑")
	app.Translate(ZhCN, "global", "Node", "节点")
	app.Translate(ZhCN, "global", "Exit", "退出")
	rets := app.LoadTranslateDict(ZhCN)
	t.Log(common.ToJson(rets))
	t.Log(common.ToJson(app.QueryTranslateTable(ZhCN, "", "")))
	app.RenderTranslateFiles()
}
