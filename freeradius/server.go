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
	"os"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"go.elastic.co/apm/module/apmechov4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

// 运行管理系统
func ListenFreeRADIUSServer(manager *models.ModelManager) error {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	if os.Getenv("ELASTIC_APM_SERVER_URL") != "" {
		e.Use(apmechov4.Middleware())
	} else {
		e.Use(ServerRecover(manager.Config.Freeradius.Debug))
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "freeradius ${time_rfc3339} ${remote_ip} ${method} ${uri} ${protocol} ${status} ${id} ${user_agent} ${latency} ${bytes_in} ${bytes_out} ${error}\n",
		Output: os.Stdout,
	}))
	// Init Handlers
	httphandler := NewHttpHandler(&WebContext{
		Manager: manager,
		Config:  manager.Config,
	})
	httphandler.InitAllRouter(e)
	e.HideBanner = true
	e.Logger.SetLevel(common.If(manager.Config.Freeradius.Debug, elog.DEBUG, elog.INFO).(elog.Lvl))
	e.Debug = manager.Config.Freeradius.Debug

	var servaddr = fmt.Sprintf("%s:%d", manager.Config.Freeradius.Host, manager.Config.Freeradius.Port)
	log.Info("try start tls web server")
	err := e.StartTLS(servaddr, path.Join(manager.Config.GetPrivateDir(), "freeradius-api.tls.crt"), path.Join(manager.Config.GetPrivateDir(), "freeradius-api.tls.key"))
	if err != nil {
		log.Warningf("start tls server error %s", err)
		log.Infof("start web server %s", servaddr)
		err = e.Start(servaddr)
	}
	return err
}

func ServerRecover(debug bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					if debug {
						log.Errorf("%+v", r)
					}
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, err.Error()))
				}
			}()
			return next(c)
		}
	}
}
