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
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/ca17/teamsacs/app"
	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
	"github.com/labstack/echo/v4"
)

func (s *ApiServer) FindCwmpFile(c echo.Context) error {
	var data []models.CwmpFile
	prequery := web.NewPreQuery(c).
		DefaultOrderBy("updated_at desc").
		KeyFields("name", "oui", "manufacturer", "product_class", "version", "file_path")

	var total int64
	common.Must(prequery.Query(app.GormDB.Model(&models.CwmpFile{})).Count(&total).Error)

	query := prequery.Query(app.GormDB.Model(&models.CwmpFile{}))
	if query.Find(&data).Error != nil {
		return c.JSON(http.StatusOK, common.EmptyList)
	}
	return c.JSON(http.StatusOK, data)
}

func (s *ApiServer) RemoveCwmpFile(c echo.Context) error {
	var id string
	err := web.NewParamReader(c).
		ReadRequiedString(&id, "id").
		LastError
	if err != nil {
		return c.JSON(200, web.RestError(err.Error()))
	}
	err = app.GormDB.Delete(&models.CwmpFile{}, "id = ?", id).Error
	if err != nil {
		return c.JSON(200, web.RestError(err.Error()))
	}
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}

// UploadCwmpFile 上传 CwmpFile， 如果指定SN， 则立即下发
func (s *ApiServer) UploadCwmpFile(c echo.Context) error {
	file, err := c.FormFile("upload")
	if err != nil {
		return c.JSON(200, web.RestError(err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(200, web.RestError(err.Error()))
	}
	defer src.Close()
	filepath := path.Join(app.Config.System.Workdir, "cwmpfile/"+file.Filename)
	dest, err := os.Create(filepath)
	if err != nil {
		return c.JSON(200, web.RestError("create file error "+err.Error()))
	}
	defer dest.Close()
	io.Copy(dest, src)
	var header = c.Request().Header
	productClass := header.Get("product_class")
	manufacturer := header.Get("manufacturer")
	oui := header.Get("oui")
	version := header.Get("version")
	fileType := header.Get("file_type")
	name := file.Filename
	format := path.Ext(file.Filename)
	err = app.GormDB.Create(&models.CwmpFile{
		ID:           common.UUID(),
		Name:         name,
		Oui:          oui,
		Manufacturer: manufacturer,
		ProductClass: productClass,
		Version:      version,
		FileType:     fileType,
		FilePath:     filepath,
		Format:       format,
		Timeout:      300,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}).Error
	if err != nil {
		return c.JSON(200, web.RestError("save file error "+err.Error()))
	}
	return c.JSON(http.StatusOK, web.RestSucc("success"))
}
