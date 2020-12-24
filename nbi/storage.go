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
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
)

func (h *HttpHandler) UploadCpeBackup(c echo.Context) error {
	cpesn := c.Param("sn")
	if cpesn == "" {
		return h.GetInternalError("cpe sn is empty")
	}
	file, err := c.FormFile("upload")
	common.Must(err)

	srcfile, err := file.Open()
	common.Must(err)
	defer srcfile.Close()

	cpeTmpfilename := fmt.Sprintf("%s/%s.tmp", os.TempDir(), common.UUID())
	dstfile, err := os.Create(cpeTmpfilename)
	common.Must(err)
	defer dstfile.Close()

	_, err = io.Copy(dstfile, srcfile)
	common.Must(err)

	containerName := "cpebackup"
	targetFilepath := fmt.Sprintf("%s/%s.backup", time.Now().Format("20060102"), cpesn )
	r, err := h.GetManager().AzureBlobC.UploadFileObject(containerName, targetFilepath, dstfile)

	common.Must(err)
	log.Infof("UploadCpeBackup done, status=%d", r.Response().StatusCode)
	return c.JSON(http.StatusOK, h.RestSucc("Operation Success"))
}
