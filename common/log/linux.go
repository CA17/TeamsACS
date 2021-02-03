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

// +build !windows

package log

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/op/go-logging"
)

func SetupSyslog(level logging.Level, syslogaddr string, module string) logging.LeveledBackend {
	var format = logging.MustStringFormatter(
		`%{pid} %{shortfile} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
	)
	backend, err := NewSyslogBackend("", syslogaddr, syslog.LOG_INFO)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	backend2Formatter := logging.NewBackendFormatter(backend, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(level, module)
	return backend1Leveled
}
