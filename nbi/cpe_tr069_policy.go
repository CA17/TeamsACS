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

package nbi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
)

// RunCpeTr069Policy
// params >>
// sn string : cpe sn
// pid string : policy id
func (h *HttpHandler) RunCpeTr069Policy(c echo.Context) error {
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(c.QueryParam("sn"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetCpeBySn error %s", err.Error())))
	}
	policy, err := h.GetManager().GetPolicyManager().GetTr069PolicyByPid(c.QueryParam("pid"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetTr069PolicyByPid error %s", err.Error())))
	}

	connParamStr := fmt.Sprintf("?timeout=%d&connection_request", 5000)
	if c.QueryParam("async") == "true" {
		connParamStr = ""
	}

	client := &http.Client{}
	client.Timeout = time.Second * 5
	url := common.UrlJoin2(
		h.GetManager().Config.Genieacs.NbiUrl,
		fmt.Sprintf("/devices/%s/tasks%s", cpe.GetDeviceId(), connParamStr))

	data, err := policy.GetTr069ParamData()
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Tr069 Param data error %s", err.Error())))
	}
	data = h.replaceVariables(data)

	if h.GetManager().Config.NBI.Debug {
		log.Infof("invoke genieacs api => %s, params = %s", url, data)
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	common.Must(err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		errstr := fmt.Sprintf("Tr069 invoke error %s", err.Error())
		log.Error(errstr)
		return c.JSON(200, h.RestError(errstr))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errstr := fmt.Sprintf("Tr069 invoke resp read error %s", err.Error())
		log.Error(errstr)
		return c.JSON(200, h.RestError(errstr))
	}
	bodystr := string(body)
	if h.GetManager().Config.Genieacs.Debug {
		log.Info(bodystr)
	}
	return c.JSON(200, h.RestResult(bodystr))
}
