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
	"strings"

	"github.com/go-routeros/routeros"
	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

func (h *HttpHandler) QueryVpes(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetParamMap("sortmap")["update_time"] = "desc"
	data, err := h.GetManager().GetVpeManager().QueryVpes(params)
	if err != nil {
		return h.GetInternalError(err)
	}
	return c.JSON(http.StatusOK, data)
}

func (h *HttpHandler) AddVpeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsVpe
	common.Must(h.GetManager().GetVpeManager().AddVpeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

func (h *HttpHandler) UpdateVpeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsVpe
	common.Must(h.GetManager().GetVpeManager().UpdateVpeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}



// RunMikrotikVpeApiPolicy
// sn string
// pid string
func (h *HttpHandler) RunMikrotikVpeApiPolicy(c echo.Context) error {
	cpe, err := h.GetManager().GetVpeManager().GetVpeBySn(c.QueryParam("sn"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetVpeBySn error %s", err.Error())))
	}
	policy, err := h.GetManager().GetPolicyManager().GetMikrotikApiPolicyByPid(c.QueryParam("pid"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetMikrotikApiPolicyByPid error %s", err.Error())))
	}
	// api params
	apiAddr := common.Must2(cpe.GetApiAddr()).(string)
	user := common.Must2(cpe.GetApiUser()).(string)
	pwdencrypt := common.Must2(cpe.GetApiPwd()).(string)
	pwd, err := aes.DecryptFromB64(pwdencrypt, h.GetManager().Config.System.Aeskey)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Api Password Decrypt error %s", err.Error())))
	}
	apiCommand := common.Must2(policy.GetApiCommand()).(string)
	apiParams := common.Must2(policy.GetApiParams()).(string)
	apiProps := common.Must2(policy.GetApiProps()).(string)

	// connect to cpe
	conn, err := routeros.Dial(apiAddr, user, pwd)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Connect Vpe error %s", err.Error())))
	}
	args := make([]string, 0)
	args = append(args, apiCommand)
	for _, p := range strings.Split(apiParams, ",") {
		if p == "" {
			continue
		}
		args = append(args, "?"+p)
	}
	if apiProps != "" {
		args = append(args, "=.proplist="+apiProps)
	}
	if h.GetManager().Config.NBI.Debug {
		log.Infof("%v", args)
	}
	reply, err := conn.Run(args...)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Execute Api error %s", err.Error())))
	}
	if h.GetManager().Config.NBI.Debug {
		log.Info(reply.String())
	}

	return c.JSON(http.StatusOK, h.RestResult(reply.Done))
}

