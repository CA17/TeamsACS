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
	"path"
	"sort"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo/v4"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
)

// A generic data CRUD management API with no predefined schema,

// QueryData
func (h *HttpHandler) QueryData(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetParamMap("sortmap")["update_time"] = "desc"
	params["collname"] = c.Param("collname")
	data, err := h.GetManager().GetDataManager().QueryItems(params,c.Param("collname"))
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

// QueryData
func (h *HttpHandler) QueryPageData(c echo.Context) error {
	params := h.RequestParse(c)
	params.GetParamMap("sortmap")["update_time"] = "desc"
	params["collname"] = c.Param("collname")
	data, err := h.GetManager().GetDataManager().QueryPagerItems(params, c.Param("collname"))
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

// QueryData
func (h *HttpHandler) QueryDataOptions(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	data, err := h.GetManager().GetDataManager().QueryItemOptions(params, c.Param("collname"))
	common.Must(err)
	return c.JSON(http.StatusOK, data)
}

// AddData
func (h *HttpHandler) GetData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	r, err := h.GetManager().GetDataManager().GetData(params)
	common.Must(err)
	return c.JSON(http.StatusOK, h.RestResult(r))
}

// AddData
func (h *HttpHandler) GetDataValues(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	r, err := h.GetManager().GetDataManager().GetDataNameValues(params)
	common.Must(err)
	return c.JSON(http.StatusOK, h.RestResult(r))
}

// AddData
func (h *HttpHandler) AddData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	common.Must(h.GetManager().GetDataManager().AddData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

// UpdateData
func (h *HttpHandler) UpdateData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	common.Must(h.GetManager().GetDataManager().UpdateData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

// DeleteData
func (h *HttpHandler) DeleteData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	common.Must(h.GetManager().GetDataManager().DeleteData(params))
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}


// SaveAction
func (h *HttpHandler) SaveAction(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	r, err := h.GetManager().GetDataManager().SaveData(params)
	common.Must(err)
	return c.JSON(http.StatusOK, r)
}


// ImportData
func (h *HttpHandler) ImportData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	collname := params.GetMustString("collname")
	items, err := h.FetchExcelData(c, collname)
	var datas = make([]interface{}, 0)
	for _, item := range items {
		_id, ok := item["_id"]
		if !ok || common.IsEmptyOrNA(_id) {
			item["_id"] = common.UUID()
		}
		datas = append(datas, item)
	}
	common.Must(h.GetManager().GetDataManager().AddBatchData(collname, datas))
	common.Must(err)
	return c.JSON(http.StatusOK, h.RestSucc("Success"))
}

var COLNAMES = map[int]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M",
	13: "N", 14: "O", 15: "P", 16: "Q", 17: "R", 18: "S", 19: "T", 20: "U", 21: "V", 22: "W", 23: "X", 24: "Y",
	25: "Z", 26: "AA", 27: "AB", 28: "AC", 29: "AD", 30: "AE", 31: "AF", 32: "AG", 33: "AH", 34: "AI", 35: "AJ",
	36: "AK", 37: "AL", 38: "AM", 39: "AN",
}

// ExportData
func (h *HttpHandler) ExportData(c echo.Context) error {
	params := h.RequestParse(c)
	params["collname"] = c.Param("collname")
	collname := params.GetMustString("collname")
	data, err := h.GetManager().GetDataManager().QueryItems(params, collname)
	common.Must(err)
	sheet := collname
	filename := fmt.Sprintf("%s-%d.xlsx", sheet, common.UUIDint64())
	filepath := path.Join(h.GetConfig().GetDataDir(), filename)
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheet)
	names := make([]string, 0)
	for i, item := range *data {
		if i == 0 {
			for k, _ := range item {
				names = append(names, k)
			}
			sort.Slice(names, func(i, j int)bool{
				return names[i] == "_id"
			})
			for j, name := range names {
				xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", COLNAMES[j], 1), name)
			}
		}
		for j, name := range names {
			xlsx.SetCellValue(sheet, fmt.Sprintf("%s%d", COLNAMES[j], i+2), item[name].(string))
		}
	}
	xlsx.SetActiveSheet(index)
	err = xlsx.SaveAs(filepath)
	if err != nil {
		log.Error(err)
		return h.GetInternalError(err)
	}
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", sheet))
	return c.File(filepath)
}
