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
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"go.elastic.co/apm/module/apmechov4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/common/tpl"
	"github.com/ca17/teamsacs/models"
)

// 运行管理系统
func ListenNBIServer(manager *models.ModelManager) error {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	if os.Getenv("ELASTIC_APM_SERVER_URL") != "" {
		e.Use(apmechov4.Middleware())
	} else {
		e.Use(ServerRecover(manager.Config.NBI.Debug))
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "nbi ${time_rfc3339} ${remote_ip} ${method} ${uri} ${protocol} ${status} ${id} ${user_agent} ${latency} ${bytes_in} ${bytes_out} ${error}\n",
		Output: os.Stdout,
	}))
	manager.WebJwtConfig = &middleware.JWTConfig{
		SigningMethod: middleware.AlgorithmHS256,
		SigningKey:    []byte(manager.Config.NBI.JwtSecret),
		Skipper: func(c echo.Context) bool {
			if strings.HasPrefix(c.Path(), "/nbi/status") ||
			strings.HasPrefix(c.Path(), "/nbi/cpe/backup/upload") ||
				strings.HasPrefix(c.Path(), "/nbi/token") {
				return true
			}
			return false
		},
		ErrorHandler: func(err error) error {
			return NewHTTPError(http.StatusBadRequest, "Missing tokens, Access denied")
		},
	}
	e.Use(middleware.JWTWithConfig(*manager.WebJwtConfig))

	// Init Handlers
	httphandler := NewHttpHandler(&WebContext{
		Manager: manager,
		Config:  manager.Config,
	})
	httphandler.InitAllRouter(e)

	manager.TplRender = tpl.NewCommonTemplate([]string{"/resources/templates"}, manager.Dev, manager.GetTemplateFuncMap())
	e.Renderer = manager.TplRender
	e.HideBanner = true
	e.Logger.SetLevel(common.If(manager.Config.NBI.Debug, elog.DEBUG, elog.INFO).(elog.Lvl))
	e.Debug = manager.Config.NBI.Debug
	log.Info("try start tls web server")
	err := e.StartTLS(fmt.Sprintf("%s:%d", manager.Config.NBI.Host, manager.Config.NBI.Port),
		path.Join(manager.Config.GetPrivateDir(), "teamsacs-nbi.tls.crt"), path.Join(manager.Config.GetPrivateDir(), "teamsacs-nbi.tls.key"))
	if err != nil {
		log.Warningf("start tls server error %s", err)
		log.Info("start web server")
		err = e.Start(fmt.Sprintf("%s:%d", manager.Config.NBI.Host, manager.Config.NBI.Port))
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
