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

package settings

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/teamsctl/apiclient"
	"github.com/urfave/cli/v2"
)

//go:embed settings-init.json
var settingsJson []byte

var Commands = []*cli.Command{
	{
		Name:     "list-settings",
		Usage:    "list settings",
		Category: "Settings",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Usage: "settings type: system | cwmp | cwmp ", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			result, err := apiclient.FindSettings(c.String("type"))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}

			return nil
		},
	},
	{
		Name:     "update-settings",
		Usage:    "update settings",
		Category: "Settings",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Usage: "settings type: system | cwmp | cwmp ", Value: ""},
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "settings name", Value: ""},
			&cli.StringFlag{Name: "value", Aliases: []string{"v"}, Usage: "settings value", Value: ""},
			&cli.StringFlag{Name: "remark", Aliases: []string{"r"}, Usage: "settings remark", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			ctype := c.String("type")
			name := c.String("name")
			value := c.String("value")
			remark := c.String("remark")
			if ctype == "" || name == "" || value == "" {
				fmt.Println("type, name, value is required")
				return nil
			}
			result, err := apiclient.UpdateSettings(models.SysConfig{
				Type:   ctype,
				Name:   name,
				Value:  value,
				Remark: remark,
			})
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}

			return nil
		},
	},
	{
		Name:     "remove-settings",
		Usage:    "remove settings",
		Category: "Settings",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "type", Aliases: []string{"t"}, Usage: "settings type: system | cwmp | cwmp ", Value: ""},
			&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "settings name", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			ctype := c.String("type")
			name := c.String("name")
			if ctype == "" || name == "" {
				fmt.Println("type, name is required")
				return nil
			}
			result, err := apiclient.RemoveSettings(ctype, name)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}

			return nil
		},
	},
	// 初始化设置
	{
		Name:     "init-settings",
		Usage:    "init settings from json file",
		Category: "Settings",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Usage: "settings json file", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			var jsondata = settingsJson
			file := c.String("file")
			if common.FileExists(file) {
				bs, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Println("read settings file error " + err.Error())
				} else {
					jsondata = bs
				}
			}
			var cfgs []models.SysConfig
			err := json.Unmarshal(jsondata, &cfgs)
			if err != nil {
				fmt.Println("parse settings file error " + err.Error())
				return nil
			}

			result, err := apiclient.UpdateSettings(cfgs...)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}
			return nil
		},
	},
}
