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
	"strings"
	"time"

	"github.com/go-routeros/routeros"
	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/aes"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/models"
)

func (h *HttpHandler) QueryCpes(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetParamMap("sortmap")["update_time"] = "desc"
	ispager := params.GetQueryMap().GetString("ispager") == "true"
	if ispager {
		data, err := h.GetManager().GetCpeManager().QueryCpes(params)
		common.Must(err)
		return c.JSON(http.StatusOK, data)
	} else {
		params.Set("limit",100)
		data, err := h.GetManager().GetCpeManager().QueryCpeList(params)
		common.Must(err)
		return c.JSON(http.StatusOK, data)
	}
}

func (h *HttpHandler) AddCpeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsCpe
	common.Must(h.GetManager().GetCpeManager().AddCpeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

func (h *HttpHandler) UpdateCpeData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = models.TeamsacsCpe
	common.Must(h.GetManager().GetCpeManager().UpdateCpeData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

func (h *HttpHandler) replaceVariables(data string) string {
	vars, err := h.GetManager().GetPolicyManager().QueryAllVars()
	if err != nil {
		log.Errorf("Query Policy variables error %s", err.Error())
		return data
	}
	// Replace Variables
	for _, v := range vars {
		name, err := v.GetName()
		if err != nil {
			continue
		}
		value, err := v.GetValue()
		if err != nil {
			continue
		}
		data = strings.ReplaceAll(data, fmt.Sprintf("{{%s}}", name), value)
	}
	return data
}

// RunCpeApiPolicy
// params >>
// sn string : cpe sn
// pid string : policy id
func (h *HttpHandler) RunMikrotikCpeApiPolicy(c echo.Context) error {
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(c.QueryParam("sn"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetCpeBySn error %s", err.Error())))
	}
	policy, err := h.GetManager().GetPolicyManager().GetMikrotikApiPolicyByPid(c.QueryParam("pid"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetMikrotikApiPolicyByPid error %s", err.Error())))
	}
	// api params
	apiAddr := common.Must2(cpe.GetApiAddr()).(string)
	user := common.Must2(cpe.GetApiUser()).(string)
	pwdencrypt := common.Must2(cpe.GetApiPwd()).(string)
	pwd, err := aes.DecryptFromB64(pwdencrypt, h.GetManager().Config.System.Aeskey)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Api Password Decrypt error %s", err.Error())))
	}
	apiCommand := common.Must2(policy.GetApiCommand()).(string)
	apiParams,_ := policy.GetApiParams()
	apiProps,_ := policy.GetApiProps()

	// connect to cpe
	conn, err := routeros.Dial(apiAddr, user, pwd)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Connect Cpe error %s", err.Error())))
	}
	args := make([]string, 0)
	args = append(args, apiCommand)
	for _, p := range strings.Split(apiParams, ",") {
		if p == "" {
			continue
		}
		args = append(args, "?"+p)
	}
	if apiProps != "" {
		args = append(args, "=.proplist="+apiProps)
	}
	if h.GetManager().Config.NBI.Debug {
		log.Infof("%v", args)
	}
	reply, err := conn.Run(args...)
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Execute Api error %s", err.Error())))
	}
	if h.GetManager().Config.NBI.Debug {
		log.Info(reply.String())
	}

	return c.JSON(http.StatusOK, h.RestResult(reply))
}

// RunCpeTr069Policy
// params >>
// sn string : cpe sn
// pid string : policy id
func (h *HttpHandler) RunCpeTr069Policy(c echo.Context) error {
	cpe, err := h.GetManager().GetCpeManager().GetCpeBySn(c.QueryParam("sn"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetCpeBySn error %s", err.Error())))
	}
	policy, err := h.GetManager().GetPolicyManager().GetTr069PolicyByPid(c.QueryParam("pid"))
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("GetTr069PolicyByPid error %s", err.Error())))
	}

	connParamStr := fmt.Sprintf("?timeout=%d&connection_request", 5000)
	if c.QueryParam("async") == "true" {
		connParamStr = ""
	}

	client := &http.Client{}
	client.Timeout = time.Second * 5
	url := common.UrlJoin2(
		h.GetManager().Config.Genieacs.NbiUrl,
		fmt.Sprintf("/devices/%s/tasks%s", cpe.GetDeviceId(), connParamStr))

	data, err := policy.GetTr069ParamData()
	if err != nil {
		return c.JSON(200, h.RestError(fmt.Sprintf("Tr069 Param data error %s", err.Error())))
	}
	data = h.replaceVariables(data)

	if h.GetManager().Config.NBI.Debug {
		log.Infof("invoke genieacs api => %s, params = %s", url, data)
	}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	common.Must(err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		errstr := fmt.Sprintf("Tr069 invoke error %s", err.Error())
		log.Error(errstr)
		return c.JSON(200, h.RestError(errstr))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errstr := fmt.Sprintf("Tr069 invoke resp read error %s", err.Error())
		log.Error(errstr)
		return c.JSON(200, h.RestError(errstr))
	}
	bodystr := string(body)
	if h.GetManager().Config.Genieacs.Debug {
		log.Info(bodystr)
	}
	return c.JSON(200, h.RestResult(bodystr))
}

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
	taskurl := common.UrlJoin2(h.GetManager().Config.Genieacs.NbiUrl, fmt.Sprintf("devices/%s/tasks%s", cpe.GetDeviceId(), connParamStr))
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
