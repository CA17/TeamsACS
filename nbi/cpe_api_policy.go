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
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/mikrotik_api"
)

// RunCpeApiPolicy
// params >>
// sn string : cpe sn
// pid string : policy id
func (h *HttpHandler) RunMikrotikCpeApiPolicy(c echo.Context) error {
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(c.QueryParam("sn"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetCpeBySn error %s", err.Error())))
	}
	policy, err := h.GetManager().GetPolicyManager().GetMikrotikApiPolicyByPid(c.QueryParam("pid"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetMikrotikApiPolicyByPid error %s", err.Error())))
	}
	// api params
	apiAddr := common.Must2(cpe.GetApiAddr()).(string)
	apiUser := common.Must2(cpe.GetApiUser()).(string)
	pwdencrypt := common.Must2(cpe.GetApiPwd()).(string)
	apiPwd, err := aes.DecryptFromB64(pwdencrypt, h.GetManager().Config.System.Aeskey)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Api Password Decrypt error %s", err.Error())))
	}
	apiCommand := common.Must2(policy.GetApiCommand()).(string)
	apiParams,_ := policy.GetApiParams()
	apiProps,_ := policy.GetApiProps()

	// connect to cpe
	api := mikrotik_api.NewMikrotikApi(apiUser, apiPwd, apiAddr, false)
	err = api.Connect()
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Connect Vpe error %s", err.Error())))
	}
	defer api.Client.Close()

	reply, err := api.ExecuteCommand(apiCommand, apiParams, apiProps)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("RunMikrotikCpeApiPolicy error %s", err.Error())))
	}
	return c.JSON(http.StatusOK, h.RestResult(reply))
}

