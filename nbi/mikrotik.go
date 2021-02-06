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

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
)

// QueryMikrotikDeviceInterfaces
func (h *HttpHandler) QueryMikrotikDeviceInterfaces(c echo.Context) error {
	var result = make(map[string]interface{})
	sn := c.QueryParam("sn")
	dev := c.QueryParam("dev")
	data, err := h.GetManager().GetMikrotikDeviceManager().QueryDeviceInterfaceList(sn, dev)
	if err != nil {
		return c.JSON(http.StatusOK, h.RestError(err.Error()))
	}
	if data != nil {
		result["data"] = data
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HttpHandler) QueryMikrotikDeviceRouters(c echo.Context) error {
	var result = make(map[string]interface{})
	sn := c.QueryParam("sn")
	dev := c.QueryParam("dev")
	data, err := h.GetManager().GetMikrotikDeviceManager().QueryDeviceRoutes(sn, dev)
	if err != nil {
		return c.JSON(http.StatusOK, h.RestError(err.Error()))
	}
	if data != nil {
		result["data"] = data
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}

func (h *HttpHandler) QueryMikrotikDeviceDnsInfo(c echo.Context) error {
	var result = make(map[string]interface{})
	sn := c.QueryParam("sn")
	dev := c.QueryParam("dev")
	data, err := h.GetManager().GetMikrotikDeviceManager().QueryDeviceDnsinfo(sn, dev)
	if err != nil {
		return c.JSON(http.StatusOK, h.RestError(err.Error()))
	}
	if data != nil {
		result["data"] = data
	} else {
		result["data"] = common.EmptyList
	}
	return c.JSON(http.StatusOK, result)
}
