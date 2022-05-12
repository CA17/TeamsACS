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

package freeradius

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/validutil"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
	"github.com/labstack/echo/v4"
)

const (
	RadiusAuthSucces        = "success"
	RadiusAuthFailure       = "failure"
	RadiusMaxSessionTimeout = 864000
	RadiusInterimIntelval   = 120
	RadiusAuthlogLevel      = "all"
)

func (s *FreeradiusServer) initRouter() {
	s.root.Add(http.MethodPost, "/freeradius/authorize", s.FreeradiusAuthorize)
	s.root.Add(http.MethodPost, "/freeradius/authenticate", s.FreeradiusAuthenticate)
	s.root.Add(http.MethodPost, "/freeradius/postauth", s.FreeradiusPostauth)
	s.root.Add(http.MethodPost, "/freeradius/accounting", s.FreeradiusAccounting)
}

// AddAuthlog 添加认证日志
func (s *FreeradiusServer) AddAuthlog(username string, nasip string, result string, reason string, level string, cast int64) {
	if level != "all" || result != level {
		err := addRadiusAuthLog(username, nasip, result, reason, cast)
		if err != nil {
			log.Error(err)
		}
	}
}

// FreeradiusAuthorize
// Authorize processing, if the user exists, the password response is sent back for further verification.
//      #  FreeradiusAuthorize/FreeradiusAuthenticate
//      #
//      #  Code   Meaning       Process body  Module code
//      #  404    not found     no            notfound
//      #  410    gone          no            notfound
//      #  403    forbidden     no            userlock
//      #  401    unauthorized  yes           reject
//      #  204    no content    no            ok
//      #  2xx    successful    yes           ok/updated
//      #  5xx    server error  no            fail
//      #  xxx    -             no            invalid
func (s *FreeradiusServer) FreeradiusAuthorize(c echo.Context) error {
	var start = time.Now()
	username := strings.TrimSpace(c.FormValue("username"))
	nasip := c.FormValue("nasip")

	var user models.RadiusUser
	err := app.GormDB.Where("username=?", username).First(&user).Error
	if err != nil {
		s.AddAuthlog(username, nasip, RadiusAuthFailure, "user query err"+err.Error(), RadiusAuthlogLevel, time.Since(start).Milliseconds())
		return c.JSON(501, echo.Map{"Reply-LatestMessage": "user query error, reject auth, " + err.Error()})
	}

	// Check user status
	if user.Status == common.DISABLED {
		s.AddAuthlog(username, nasip, RadiusAuthFailure, "user disabled", RadiusAuthlogLevel, time.Since(start).Milliseconds())
		return c.JSON(501, echo.Map{"Reply-LatestMessage": "user status disabled, reject auth"})
	}

	var expireTime = time.Time(user.ExpireTime)
	// Check user expiration
	if expireTime.Before(time.Now()) {
		s.AddAuthlog(username, nasip, RadiusAuthFailure, "user expire", RadiusAuthlogLevel, time.Since(start).Milliseconds())
		return c.JSON(501, echo.Map{"Reply-LatestMessage": "user expire, reject auth"})
	}

	// Evaluation of online limit
	// Current number online
	count, err := getOnlineCount(username)
	if err != nil {
		s.AddAuthlog(username, nasip, RadiusAuthFailure, "user query count err"+err.Error(), RadiusAuthlogLevel, time.Since(start).Milliseconds())
		return c.JSON(501, echo.Map{"Reply-LatestMessage": "user online count fetch error, reject auth, " + err.Error()})
	}
	var activeNum = user.ActiveNum
	if count > 0 && activeNum > 0 && count >= int64(activeNum) {
		s.AddAuthlog(username, nasip, RadiusAuthFailure, "user online limit", RadiusAuthlogLevel, time.Since(start).Milliseconds())
		return c.JSON(501, echo.Map{"Reply-LatestMessage": "user online over limit, reject auth"})
	}

	// freeradius response
	var password = user.Password
	resp := map[string]interface{}{}
	resp["control:Cleartext-Password"] = strings.TrimSpace(password)
	resp["reply:Mikrotik-Rate-Limit"] = fmt.Sprintf("%dk/%dk", user.UpRate, user.DownRate)
	sessionTimeout := expireTime.Sub(time.Now()).Seconds()
	resp["reply:Session-Timeout"] = fmt.Sprintf("%d", int64(sessionTimeout))
	resp["reply:Acct-Interim-Interval"] = 120

	// Set address pool or static IP
	var userip = user.IpAddr
	var addrpool = user.AddrPool
	if common.IsNotEmptyAndNA(userip) && validutil.IsIP(userip) {
		resp["Framed-IP-Address"] = userip
	} else if common.IsNotEmptyAndNA(addrpool) {
		resp["Framed-Pool"] = addrpool
	}

	s.AddAuthlog(username, nasip, RadiusAuthSucces, RadiusAuthSucces, RadiusAuthlogLevel, time.Since(start).Milliseconds())

	return c.JSON(http.StatusOK, resp)
}

// FreeradiusAuthenticate
// Authenticate processing
//     #  FreeradiusAuthorize/FreeradiusAuthenticate
//     #
//     #  Code   Meaning       Process body  Module code
//     #  404    not found     no            notfound
//     #  410    gone          no            notfound
//     #  403    forbidden     no            userlock
//     #  401    unauthorized  yes           reject
//     #  204    no content    no            ok
//     #  2xx    successful    yes           ok/updated
//     #  5xx    server error  no            fail
//     #  xxx    -             no            invalid
func (s *FreeradiusServer) FreeradiusAuthenticate(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

// FreeradiusPostauth Postauth processing
func (s *FreeradiusServer) FreeradiusPostauth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

// FreeradiusAccounting Accounting processing
func (s *FreeradiusServer) FreeradiusAccounting(c echo.Context) error {
	webform := web.NewWebForm(c)
	err := updateRadiusOnline(webform)
	if err != nil {
		log.Error(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
