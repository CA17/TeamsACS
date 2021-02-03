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

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/models"
)

func (h *HttpHandler) RequestToken(c echo.Context) error {
	params, err := h.ParseJsonBody(c)
	common.Must(err)
	username := params.GetMustString("username")
	apisecret := params.GetMustString("apisecret")
	opr, err := h.GetManager().GetOpsManager().GetOperator(username)
	common.Must(err)
	if opr.ApiSecret != apisecret {
		return h.GetInternalError("apisecret error")
	}
	t, err := h.CreateAuthToken(opr)
	common.Must(err)
	return c.JSON(http.StatusOK, h.RestResult(map[string]string{
		"token": t,
	}))
}

func (h *HttpHandler) CreateAuthToken(operator *models.Operator) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["usr"] = operator.Username
	claims["uid"] = operator.ID
	claims["lvl"] = operator.Level
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(h.GetConfig().NBI.JwtSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func (h *HttpHandler) Status(c echo.Context) error {
	return c.NoContent(200)
}
