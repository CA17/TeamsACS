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
	_ "embed"
	"fmt"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/models"
	"github.com/ca17/teamsacs/teamsctl/apiclient"
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	{
		Name:     "list-file",
		Usage:    "list files",
		Category: "Files",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "query", Aliases: []string{"q"}, Usage: "query keyword", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			result, err := apiclient.FindCwmpFile(c.String("type"))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}

			return nil
		},
	}, {
		Name:     "list-file-task",
		Usage:    "list file download tasks",
		Category: "Files",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Aliases: []string{"id"}, Usage: "query id", Value: ""},
			&cli.StringFlag{Name: "status", Aliases: []string{"s"}, Usage: "query status  initialize | success | error | timeout", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			result, err := apiclient.FindCwmpFileTask(c.String("id"), c.String("status"))
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}
			return nil
		},
	},
	{
		Name:     "upload-file",
		Usage:    "upload file",
		Category: "Files",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "filetype", Aliases: []string{"t"}, Usage: "file type: 1 Firmware Upgrade Image | 2 Web Content | 3 Vendor Configuration File ", Value: 0},
			&cli.StringFlag{Name: "filepath", Aliases: []string{"f"}, Usage: "filepath", Value: ""},
			&cli.StringFlag{Name: "manufacturer", Aliases: []string{"m"}, Usage: "device manufacturer", Value: ""},
			&cli.StringFlag{Name: "product_class", Aliases: []string{"c"}, Usage: "device product_class", Value: ""},
			&cli.StringFlag{Name: "oui", Aliases: []string{"o"}, Usage: "device oui", Value: ""},
			&cli.StringFlag{Name: "version", Aliases: []string{"v"}, Usage: "device version", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			filetype := c.Int("filetype")
			filepath := c.String("filepath")
			manufacturer := c.String("manufacturer")
			productClass := c.String("product_class")
			oui := c.String("oui")
			version := c.String("version")
			if filetype == 0 || filepath == "" {
				fmt.Println("filetype, filename is required")
				return nil
			}
			var filetypeStr string
			switch filetype {
			case 1:
				filetypeStr = "1 Firmware Upgrade Image"
			case 2:
				filetypeStr = "2 Web Content"
			case 3:
				filetypeStr = "3 Vendor Configuration File"
			default:
				fmt.Println("filetype error")
				return nil
			}

			if !common.FileExists(filepath) {
				fmt.Println("file not exists")
				return nil
			}

			cfile := models.CwmpFile{
				Oui:          oui,
				Manufacturer: manufacturer,
				ProductClass: productClass,
				Version:      version,
				FileType:     filetypeStr,
				FilePath:     filepath,
			}
			result, err := apiclient.UploadCwmpFile(cfile)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}

			return nil
		},
	},
	{
		Name:     "remove-file",
		Usage:    "remove file",
		Category: "Files",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Aliases: []string{"id"}, Usage: "file id", Value: ""},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "debug mode"},
		},
		Action: func(c *cli.Context) error {
			apiclient.SetDebug(c.Bool("debug"))
			id := c.String("id")
			if id == "" {
				fmt.Println("file id is required")
				return nil
			}
			result, err := apiclient.RemoveCwmpFile(id)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(common.ToJson(result))
			}

			return nil
		},
	},
}
