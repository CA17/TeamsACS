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

package apiclient

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/guonaihong/gout"
)

var (
	api    *TeamsacsApi
	apiurl = "http://127.0.0.1:8000"
	secret = "9b6de5cc-0731-4bf1-xxxx-0f568ac9da37"
	debug  bool
)

type TeamsacsApi struct {
	Apiurl string
	Secret string
	Debug  bool
}

func NewTeamsacsApi(apiurl string, secret string, debug bool) *TeamsacsApi {
	return &TeamsacsApi{Apiurl: apiurl, Secret: secret, Debug: debug}
}

func init() {
	if os.Getenv("TEAMSACS_API_URL") != "" {
		apiurl = os.Getenv("TEAMSACS_API_URL")
	}
	if os.Getenv("TEAMSACS_API_SECRET") != "" {
		secret = os.Getenv("TEAMSACS_API_SECRET")
	}
	debug = os.Getenv("TEAMSACS_API_DEBUG") == "true"
	api = NewTeamsacsApi(apiurl, secret, debug)
}

func (a *TeamsacsApi) CreateAuthorization() gout.H {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["usr"] = "teamsctl"
	claims["uid"] = "teamsctl"
	claims["lvl"] = "api"
	claims["exp"] = time.Now().Add(time.Second * 300).Unix()
	t, _ := token.SignedString([]byte(a.Secret))
	return gout.H{"authorization": "Bearer " + t}
}

func SetDebug(debug bool) {
	api.Debug = debug
}
