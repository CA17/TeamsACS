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

package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/httpc"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/webserver"
	"github.com/labstack/echo/v4"
)

func fetchFileData(filename string) ([]byte, error) {
	if common.FileExists(filename) {
		return os.ReadFile(filename)
	}
	if strings.HasPrefix(filename, "http") {
		return httpc.Get(filename, nil, nil, time.Second*10)
	}
	return nil, fmt.Errorf("file not found: %s", filename)
}

func InitRouter() {

	webserver.GET("/admin/files", func(c echo.Context) error {
		return c.Render(http.StatusOK, "files", map[string]interface{}{})
	})

	webserver.GET("/admin/files/query", func(c echo.Context) error {
		type datafile struct {
			Filename   string `json:"filename"`
			Size       int64  `json:"size"`
			Mode       string `json:"mode"`
			UpdateTime string `json:"update_time"`
		}
		dirpath := path.Join(app.GConfig().System.Workdir, "cwmp")
		// if !common.FileExists(dirpath) {
		// 	os.MkdirAll(dirpath, 0644)
		// }
		flist, err := os.ReadDir(dirpath)
		common.Must(err)
		var datas []datafile
		for _, entry := range flist {
			finfo, err := entry.Info()
			if err != nil || finfo.IsDir() {
				continue
			}
			datas = append(datas, datafile{
				Filename:   entry.Name(),
				Size:       finfo.Size(),
				Mode:       finfo.Mode().Perm().String(),
				UpdateTime: finfo.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
		return c.JSON(http.StatusOK, &datas)
	})

	webserver.POST("/admin/files/upload", func(c echo.Context) error {
		file, err := c.FormFile("upload")
		common.Must(err)
		src, err := file.Open()
		common.Must(err)
		defer src.Close()
		filename := file.Filename
		savefile := path.Join(app.GConfig().System.Workdir, "cwmp", filename)
		dst, err := os.Create(savefile)
		common.Must(err)
		defer dst.Close()
		_, err = io.Copy(dst, src)
		common.Must(err)
		return c.JSON(http.StatusOK, web.RestSucc("success"))
	})

	webserver.GET("/admin/files/download/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
			return c.String(http.StatusBadRequest, "bad request")
		}
		c.Response().Header().Set("Content-Disposition", "attachment;filename="+filename)
		return c.File(path.Join(app.GConfig().System.Workdir, "cwmp", filename))
	})

	webserver.GET("/admin/files/delete/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
			return c.String(http.StatusBadRequest, "bad request")
		}
		err := os.Remove(path.Join(app.GConfig().System.Workdir, "cwmp", filename))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, web.RestSucc("success"))
	})

}
