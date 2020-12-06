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

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

func (h *HttpHandler) QueryCpes(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetParamMap("sortmap")["update_time"] = "desc"
	ispager := params.GetQueryMap().GetString("ispager") == "true"
	if ispager {
		data, err := h.GetManager().GetCpeManager().QueryCpes(params)
		common.Must(err)
		return c.JSON(http.StatusOK, data)
	} else {
		params.Set("limit",100)
		data, err := h.GetManager().GetCpeManager().QueryCpeList(params)
		common.Must(err)
		return c.JSON(http.StatusOK, data)
	}
}

func (h *HttpHandler) AddCpeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsCpe
	common.Must(h.GetManager().GetCpeManager().AddCpeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

func (h *HttpHandler) UpdateCpeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsCpe
	common.Must(h.GetManager().GetCpeManager().UpdateCpeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

func (h *HttpHandler) replaceVariables(data string) string {
	vars, err := h.GetManager().GetPolicyManager().QueryAllVars()
	if err != nil {
		log.Errorf("Query Policy variables error %s", err.Error())
		return data
	}
	// Replace Variables
	for _, v := range vars {
		name, err := v.GetName()
		if err != nil {
			continue
		}
		value, err := v.GetValue()
		if err != nil {
			continue
		}
		data = strings.ReplaceAll(data, fmt.Sprintf("{{%s}}", name), value)
	}
	return data
}
