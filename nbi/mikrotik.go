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
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"

	"github.com/go-routeros/routeros"
)

// InvokeMikrotikApi
func (h *HttpHandler) InvokeMikrotikVpeApi(c echo.Context) error {
	vpe, err := h.GetManager().GetVpeManager().GetVpeBySn(c.QueryParam("sn"))
	common.Must(err)
	apiAddr, err := vpe.GetApiAddr()
	common.Must(err)
	user, err := vpe.GetApiUser()
	common.Must(err)
	pwd, err := vpe.GetApiPwd()
	common.Must(err)
	conn, err := routeros.Dial(apiAddr, user, pwd)
	common.Must(err)
	args := make([]string, 0)
	args = append(args, c.FormValue("command"))
	params := c.FormValue("params")
	for _, p := range strings.Split(params, ",") {
		if p == "" {
			continue
		}
		args = append(args, "?"+p)
	}
	props := c.FormValue("props")
	if props != "" {
		args = append(args, "=.proplist="+props)
	}
	reply, err := conn.Run(args...)
	common.Must(err)
	log.Info(reply.String())

	return c.JSON(http.StatusOK, h.RestResult(reply.Done))
}

// QueryMikrotikDeviceInterfaces
func (h *HttpHandler) QueryMikrotikDeviceInterfaces(c echo.Context) error {
	var result = make(map[string]interface{})
	data, err := h.GetManager().GetGenieacsManager().QueryMikrotikEthernetInterface(c.QueryParam("sn"))
	common.Must(err)
	if data != nil {
		result["data"] = data.Items
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HttpHandler) QueryMikrotikDevicePPPInterfaces(c echo.Context) error {
	var result = make(map[string]interface{})
	data, err := h.GetManager().GetGenieacsManager().QueryMikrotikPPPInterface(c.QueryParam("sn"))
	common.Must(err)
	if data != nil {
		result["data"] = data.Items
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HttpHandler) QueryMikrotikDeviceRouters(c echo.Context) error {
	var result = make(map[string]interface{})
	data, err := h.GetManager().GetGenieacsManager().QueryMikrotikRouters(c.QueryParam("sn"))
	common.Must(err)
	if data != nil {
		result["data"] = data.Items
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HttpHandler) QueryMikrotikDeviceDnsClientServer(c echo.Context) error {
	var result = make(map[string]interface{})
	data, err := h.GetManager().GetGenieacsManager().QueryMikrotikDnsServer(c.QueryParam("sn"))
	common.Must(err)
	if data != nil {
		result["data"] = data.Items
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HttpHandler) QueryMikrotikDeviceIpInterfaces(c echo.Context) error {
	var result = make(map[string]interface{})
	data, err := h.GetManager().GetGenieacsManager().QueryMikrotikIpInterface(c.QueryParam("sn"))
	common.Must(err)
	if data != nil {
		result["data"] = data.Items
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}
