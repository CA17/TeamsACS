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
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/maputils"
	"github.com/ca17/teamsacs/common/validutil"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/mikrotik_api"
)

// SyncUser
// sync socks account
func (h *HttpHandler) SyncManageUser(c echo.Context) error {
	var errorResult = func(err error) error {
		return c.JSON(200, h.RestError(err.Error()))
	}
	log.Info("SyncManageUser Start")
	frm := web.NewWebForm(c)
	vpeSn, err := frm.GetMustVal("vpe_sn")
	if err != nil {
		return errorResult(err)
	}
	cpeSn, err := frm.GetMustVal("cpe_sn")
	if err != nil {
		return errorResult(err)
	}
	vpe, err := h.GetManager().GetVpeManager().GetVpeBySn(vpeSn)
	if err != nil {
		return errorResult(err)
	}
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(cpeSn)
	log.Info("SyncManageUser Start -> parse params")
	// api params
	apiAddr,err := vpe.GetApiAddr()
	if err != nil {
		return errorResult(err)
	}
	apiUser,err := vpe.GetApiUser()
	if err != nil {
		return errorResult(err)
	}
	pwdencrypt,err := vpe.GetApiPwd()
	if err != nil {
		return errorResult(err)
	}
	apiPwd, err := aes.DecryptFromB64(pwdencrypt, h.GetManager().Config.System.Aeskey)
	if err != nil {
		return errorResult(fmt.Errorf("api Password Decrypt error %s", err.Error()))
	}
	log.Info("SyncManageUser -> connect to vpe")
	// connect to cpe
	api := mikrotik_api.NewMikrotikApi(apiUser, apiPwd, apiAddr, false)
	err = api.Connect()
	if err != nil {
		return errorResult(fmt.Errorf("connect Vpe error %s", err.Error()))
	}
	defer api.Client.Close()

	rdIpaddr, err := maputils.GetStringValueWithErr(*cpe, "rd_ipaddr")
	if err != nil {
		return errorResult(err)
	}
	if !validutil.IsIP(rdIpaddr) {
		return errorResult(fmt.Errorf("cpe ip format error"))
	}

	rdGateway, err := maputils.GetStringValueWithErr(*cpe, "rd_gateway")
	if err != nil {
		return errorResult(err)
	}
	if !validutil.IsIP(rdGateway) {
		return errorResult(fmt.Errorf("cpe gateway ip format error"))
	}

	remark, _ := maputils.GetStringValueWithErr(*cpe, "remark")
	remark = fmt.Sprintf("%s updated:%s", remark, time.Now().String())

	log.Info("SyncManageUser -> exec api")
	_ = api.RemovePPPUser(cpeSn)
	hexremark, _ := common.ToGbkHexString(remark)
	err = api.AddPPPUser(cpeSn, cpeSn, rdIpaddr, rdGateway, strings.ReplaceAll(hexremark, "\\", "%"))
	if err != nil {
		return errorResult(fmt.Errorf("VPE SyncManageUser Execute Api error %s", err.Error()))
	}

	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}
