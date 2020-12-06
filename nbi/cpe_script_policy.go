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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
)

// RunCpeTr069Policy
// params >>
// sn string : cpe sn
// pid string : policy id
func (h *HttpHandler) RunMikrotikCpeScriptPolicy(c echo.Context) error {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	// Query CpeData
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(c.QueryParam("sn"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetCpeBySn error %s", err.Error())))
	}

	// Query Script PolicyData
	policy, err := h.GetManager().GetPolicyManager().GetMikrotikScriptPolicyByPid(c.QueryParam("pid"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetMikrotikScriptPolicyByPid error %s", err.Error())))
	}

	busyDeviceIds, err := h.GetManager().GetGenieacsManager().GetAcsTaskDeviceIdList()
	if err != nil {
		log.Error(err)
	}

	// CPE can only perform one task at a time.
	if busyDeviceIds != nil && len(busyDeviceIds) > 0 && common.InSlice(cpe.GetDeviceId(), busyDeviceIds) {
		return fmt.Errorf("cpe %s there are still unexecuted tasks, waiting or deleting old ones", cpe.GetDeviceId())
	}

	filename := fmt.Sprintf("%s-latest-task.alter", cpe.GetDeviceId())
	fileurl := common.UrlJoin(h.GetManager().Config.Genieacs.NbiUrl, "files/"+filename)

	// Delete old script
	if h.GetManager().Config.Genieacs.Debug {
		log.Infof("delete old file => %s", fileurl)
	}
	delreq, err := http.NewRequest(http.MethodDelete, fileurl, nil)
	common.Must(err)
	delresp, err := client.Do(delreq)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Delete old script error %s", err.Error())))
	}
	if h.GetManager().Config.Genieacs.Debug {
		log.Infof("delete old file resp statusCode:", delresp.StatusCode)
	}

	// Assembly script content pushed to GenieACS
	scriptValue, err := policy.GetScriptContent()
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Get Policy Script content error %s", err.Error())))
	}

	scriptValue = h.replaceVariables(scriptValue)

	// put new script
	if h.GetManager().Config.Genieacs.Debug {
		log.Infof("put new file => %s", fileurl)
		log.Info(scriptValue)
	}
	addfilereq, err := http.NewRequest(http.MethodPut, fileurl, bytes.NewReader([]byte(scriptValue)))
	common.Must(err)
	addfilereq.Header.Set("Content-Type", "application/json")
	addfilereq.Header.Set("fileType", "3 Vendor Configuration File")
	addfileresp, err := client.Do(addfilereq)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Put new script error %s", err.Error())))
	}
	if h.GetManager().Config.Genieacs.Debug {
		log.Infof("put new file resp statusCode:", addfileresp.StatusCode)
	}

	// dispatch acs task
	connParamStr := fmt.Sprintf("?timeout=%d&connection_request", 5000)
	if c.QueryParam("async") == "true" {
		connParamStr = ""
	}
	taskurl := common.UrlJoin2(h.GetManager().Config.Genieacs.NbiUrl, fmt.Sprintf("/devices/%s/tasks%s", cpe.GetDeviceId(), connParamStr))
	taskdatamap := map[string]string{"name": "download", "fileName": filename, "fileType": "3 Vendor Configuration File"}
	datajson, err := json.MarshalIndent(&taskdatamap, "", "\t")
	if err != nil {
		return err
	}

	if h.GetManager().Config.Genieacs.Debug {
		log.Infof("Tr069 post script task => %s", taskurl)
		log.Info(string(datajson))
	}

	taskreq, err := http.NewRequest(http.MethodPost, taskurl, bytes.NewReader(datajson))
	common.Must(err)
	taskreq.Header.Set("Content-Type", "application/json")
	taskresp, err := client.Do(taskreq)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Tr069 Post script task error %s", err.Error())))
	}
	if h.GetManager().Config.Genieacs.Debug {
		log.Infof("Tr069 Post script task resp statusCode:", taskresp.StatusCode)
	}

	defer taskresp.Body.Close()

	body, err := ioutil.ReadAll(taskresp.Body)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Tr069 script task resp read error %s", err.Error())))
	}
	bodystr := string(body)
	if h.GetManager().Config.Genieacs.Debug {
		log.Info(bodystr)
	}
	return c.JSON(200, h.RestResult(bodystr))

}

