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
	"reflect"
	"testing"

	"github.com/ca17/teamsacs/common/web"
	"github.com/ca17/teamsacs/models"
)

func TestFindSettings(t *testing.T) {
	api.Debug = true
	type args struct {
		ctype string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.SysConfig
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "t1",
			args:    args{},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindSettings(tt.args.ctype)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindSettings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateSettings(t *testing.T) {
	api.Debug = true
	type args struct {
		cfgs []models.SysConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *web.WebRestResult
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "t2",
			args: args{
				cfgs: []models.SysConfig{
					{
						Type:   "test",
						Name:   "test",
						Value:  "test",
						Remark: "test",
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateSettings(tt.args.cfgs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateSettings() got = %v, want %v", got, tt.want)
			}
		})
	}
}
