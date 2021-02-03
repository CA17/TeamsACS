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
)

func (h *HttpHandler) QuerySyslog(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetParamMap("sortmap")["timestamp"] = "desc"
	data, err := h.GetManager().GetOpsManager().QuerySyslog(params)
	if err != nil {
		return h.GetInternalError(err)
	}
	return c.JSON(http.StatusOK, data)
}
