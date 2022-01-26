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

package apiserver

import (
	"net/http"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/models"
	"github.com/labstack/echo/v4"
)

func (s *ApiServer) ApiStatus(c echo.Context) error {
	return c.JSON(200, "")
}

func (s *ApiServer) FindSettings(c echo.Context) error {
	ctype := c.QueryParam("type")
	var data []models.SysConfig
	query := app.GormDB
	if ctype != "" {
		query = query.Where("type", ctype)
	}
	if err := query.Order("sort").Find(&data).Error; err != nil {
		return c.JSON(http.StatusOK, data)
	}
	return c.JSON(http.StatusOK, data)
}
