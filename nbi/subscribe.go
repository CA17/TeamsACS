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
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/models"
)

func (h *HttpHandler) QuerySubscribes(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetSortMap()["update_time"] = "desc"
	// Check for expiring accounts
	expireDays := params.GetQueryMap().GetInt64("expire_days")
	if expireDays > 0 {
		trmap := params.GetTimeRangeMap()
		trmap["end"] = "expire_time"
		trmap["end_value"] = time.Now().Add(time.Hour * 24 * time.Duration(expireDays)).Format("2006-01-02 15:04:05")
	}

	keyword :=  params.GetQueryMap().GetString("keyword")
	if keyword != "" {
		filtermap := params.GetFilterMap()
		filtermap["filter[remark]"] = keyword
	}

	data, err := h.GetManager().GetSubscribeManager().QuerySubscribes(params)
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

func (h *HttpHandler) AddSubscribeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsSubscribe
	common.Must(h.GetManager().GetSubscribeManager().AddSubscribeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

