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

package timeutil

import (
	"testing"
	"time"
)

func TestFormatLenTime(t *testing.T) {
	t.Log(FmtDatetime14String(time.Now()))
	t.Log(FmtDatetime8String(time.Now()))
	t.Log(FmtDatetime6String(time.Now()))
	t.Log(FmtDateString(time.Now()))
	t.Log(FmtDatetimeString(time.Now()))
	t.Log(FmtDatetimeMString(time.Now()))
}

func TestTZ(t *testing.T) {
	tz, err := time.LoadLocation("Etc/UTC")
	t.Log(tz, err)
	t.Log(time.Now().In(tz).Format("2006-01-02 15:04:05"))

}
