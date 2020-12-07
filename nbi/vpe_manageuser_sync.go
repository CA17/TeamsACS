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
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/validutil"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/mikrotik_api"
)

// SyncUser
// sync socks account
func (h *HttpHandler) SyncManageUser(c echo.Context) error {
	frm := web.NewWebForm(c)
	vpeSn, err := frm.GetMustVal("vpe_sn")
	common.Must(err)
	cpeSn, err := frm.GetMustVal("cpe_sn")
	common.Must(err)
	vpe, err := h.GetManager().GetVpeManager().GetVpeBySn(vpeSn)
	common.Must(err)
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(cpeSn)
	// api params
	apiAddr := common.Must2(vpe.GetApiAddr()).(string)
	apiUser := common.Must2(vpe.GetApiUser()).(string)
	pwdencrypt := common.Must2(vpe.GetApiPwd()).(string)
	apiPwd, err := aes.DecryptFromB64(pwdencrypt, h.GetManager().Config.System.Aeskey)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Api Password Decrypt error %s", err.Error())))
	}

	// connect to cpe
	api := mikrotik_api.NewMikrotikApi(apiUser, apiPwd, apiAddr, false)
	err = api.Connect()
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Connect Vpe error %s", err.Error())))
	}
	defer api.Client.Close()

	rdIpaddr, err := maputils.GetStringValueWithErr(*cpe, "rd_ipaddr")
	common.Must(err)
	if !validutil.IsIP(rdIpaddr) {
		return h.GetInternalError("cpe ip format error")
	}
	rdGateway, err := maputils.GetStringValueWithErr(*cpe, "rd_gateway")
	common.Must(err)
	if !validutil.IsIP(rdGateway) {
		return h.GetInternalError("cpe gateway ip format error")
	}

	_ = api.RemovePPPUser(cpeSn)
	err = api.AddPPPUser(cpeSn, cpeSn, rdIpaddr, rdGateway)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("VPE SyncManageUser Execute Api error %s", err.Error())))
	}

	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}
