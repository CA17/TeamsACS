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

package config

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/ca17/teamsacs/common"
)

type MongodbConfig struct {
	Url    string `yaml:"url" json:"url"`
	User   string `yaml:"user" json:"user"`
	Passwd string `yaml:"passwd" json:"passwd"`
}

type AzureStorageConfig struct {
	AccountName string `yaml:"account_name"`
	AccountKey  string `yaml:"account_key"`
}

type SysConfig struct {
	Appid      string `yaml:"appid" json:"appid"`
	Workdir    string `yaml:"workdir" json:"workdir"`
	SyslogAddr string `yaml:"syslog_addr" json:"syslog_addr"`
	Location   string `yaml:"location" json:"location"`
	Aeskey     string `yaml:"aeskey" json:"aeskey"`
	Debug      bool   `yaml:"debug" json:"debug"`
}

type NBIConfig struct {
	Host      string `yaml:"host" json:"host"`
	Port      int    `yaml:"port" json:"port"`
	Debug     bool   `yaml:"debug" json:"debug"`
	JwtSecret string `yaml:"jwt_secret" json:"jwt_secret"`
}

type ElasticConfig struct {
	Urls  string `yaml:"urls" json:"urls"`
	User  string `yaml:"user" json:"user"`
	Pwd   string `yaml:"pwd" json:"pwd"`
	Debug bool   `yaml:"debug" json:"debug"`
}

type FreeradiusConfig struct {
	Host  string `yaml:"host" json:"host"`
	Port  int    `yaml:"port" json:"port"`
	Debug bool   `yaml:"debug" json:"debug"`
}

type GenieacsConfig struct {
	NbiUrl string `yaml:"nbi_url" json:"nbi_url"`
	Debug  bool   `yaml:"debug" json:"debug"`
}

type GrpcConfig struct {
	Host  string `yaml:"host" json:"host"`
	Port  int    `yaml:"port" json:"port"`
	Debug bool   `yaml:"debug" json:"debug"`
}

type RadiusdConfig struct {
	Host     string `yaml:"host" json:"host"`
	AuthPort int    `yaml:"auth_port" json:"auth_port"`
	AcctPort int    `yaml:"acct_port" json:"acct_port"`
	Debug    bool   `yaml:"debug" json:"debug"`
}

type SyslogdConfig struct {
	Host        string `yaml:"host" json:"host"`
	Rfc5424Port int    `yaml:"rfc_5424_port" json:"rfc_5424_port"`
	Rfc3164Port int    `yaml:"rfc_3164_port" json:"rfc_3164_port"`
	TextlogPort int    `yaml:"textlog_port" json:"textlog_port"`
	MaxRecodes  int    `yaml:"max_recodes" json:"max_recodes"`
	Debug       bool   `yaml:"debug" json:"debug"`
}

type AppConfig struct {
	System       SysConfig          `yaml:"system" json:"system"`
	NBI          NBIConfig          `yaml:"nbi" json:"nbi"`
	Freeradius   FreeradiusConfig   `yaml:"freeradius" json:"freeradius"`
	Genieacs     GenieacsConfig     `yaml:"genieacs" json:"genieacs"`
	Mongodb      MongodbConfig      `yaml:"mongodb" json:"mongodb"`
	Grpc         GrpcConfig         `yaml:"grpc" json:"grpc"`
	Radiusd      RadiusdConfig      `yaml:"radiusd" json:"radiusd"`
	Elastic      ElasticConfig      `yaml:"elastic" json:"elastic"`
	Syslogd      SyslogdConfig      `yaml:"syslogd" json:"syslogd"`
	AzureStorage AzureStorageConfig `yaml:"azure_storage"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.System.Workdir, "logs")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.System.Workdir, "data")
}

func (c *AppConfig) GetRadiusDir() string {
	return path.Join(c.System.Workdir, "radius")
}

func (c *AppConfig) GetPrivateDir() string {
	return path.Join(c.System.Workdir, "private")
}

func (c *AppConfig) GetResourceDir() string {
	return path.Join(c.System.Workdir, "resource")
}

func (c *AppConfig) GetBackupDir() string {
	return path.Join(c.System.Workdir, "backup")
}

func (c *AppConfig) InitDirs() {
	os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0700)
	os.MkdirAll(path.Join(c.System.Workdir, "radius"), 0700)
	os.MkdirAll(path.Join(c.System.Workdir, "data"), 0700)
	os.MkdirAll(path.Join(c.System.Workdir, "public"), 0700)
	os.MkdirAll(path.Join(c.System.Workdir, "private"), 0700)
	os.MkdirAll(path.Join(c.System.Workdir, "resource"), 0700)
	os.MkdirAll(path.Join(c.System.Workdir, "backup"), 0644)
}

var DefaultAppConfig = &AppConfig{
	System: SysConfig{
		Appid:      "TeamsACS",
		Workdir:    "/var/teamsacs",
		SyslogAddr: "",
		Location:   "Asia/Shanghai",
		Aeskey:     "5f8923be3da19452d3acdc9e69fa24e6",
	},
	NBI: NBIConfig{
		Host:      "0.0.0.0",
		Port:      1979,
		Debug:     true,
		JwtSecret: "9b6de5cc-0731-4bf1-zpms-0f568ac9da37",
	},
	Freeradius: FreeradiusConfig{
		Host:  "0.0.0.0",
		Port:  1980,
		Debug: true,
	},
	Genieacs: GenieacsConfig{
		NbiUrl: "http://127.0.0.1:7557",
		Debug:  true,
	},
	Grpc: GrpcConfig{
		Host:  "0.0.0.0",
		Port:  1981,
		Debug: true,
	},
	Radiusd: RadiusdConfig{
		Host:     "0.0.0.0",
		AuthPort: 1812,
		AcctPort: 1813,
		Debug:    true,
	},
	Elastic: ElasticConfig{
		Urls:  "http://172.26.1.163:9200",
		User:  "",
		Pwd:   "",
		Debug: true,
	},
	Syslogd: SyslogdConfig{
		Host:        "0.0.0.0",
		Rfc5424Port: 1914,
		Rfc3164Port: 1924,
		TextlogPort: 1934,
		MaxRecodes:  100000,
		Debug:       true,
	},
	Mongodb: MongodbConfig{
		Url:    "mongodb://127.0.0.1:27017",
		User:   "",
		Passwd: "",
	},
	AzureStorage: AzureStorageConfig{
		AccountName: "",
		AccountKey:  "",
	},
}

func setEnvValue(name string, f func(v string)) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		f(evalue)
	}
}
func setEnvInt64Value(name string, f func(v int64)) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}
	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		f(p)
	}
}

func LoadConfig(cfile string) *AppConfig {

	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(ioutil.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultAppConfig
	}

	setEnvValue("TEAMSACS_WORKER_DIR", func(v string) {
		cfg.System.Workdir = v
	})

	// Acs Config
	setEnvValue("TEAMSACS_NBI_HOST", func(v string) {
		cfg.NBI.Host = v
	})
	setEnvValue("TEAMSACS_NBI_DEBUG", func(v string) {
		cfg.NBI.Debug = v == "true"
	})
	setEnvValue("TEAMSACS_NBI_SECRET", func(v string) {
		cfg.NBI.JwtSecret = v
	})
	setEnvInt64Value("TEAMSACS_NBI_PORT", func(v int64) {
		cfg.NBI.Port = int(v)
	})

	// GenieAcs Config
	setEnvValue("TEAMSACS_GENIEACS_NBIURL", func(v string) {
		cfg.Genieacs.NbiUrl = v
	})

	// FreeRADIUS Config
	setEnvValue("TEAMSACS_FREERADIUS_WEB_HOST", func(v string) {
		cfg.Freeradius.Host = v
	})

	setEnvValue("TEAMSACS_FREERADIUS_WEB_DEBUG", func(v string) {
		cfg.Freeradius.Debug = v == "true"
	})
	setEnvInt64Value("TEAMSACS_FREERADIUS_WEB_PORT", func(v int64) {
		cfg.Freeradius.Port = int(v)
	})

	// Mongodb Config
	setEnvValue("TEAMSACS_MONGODB_URL", func(v string) {
		cfg.Mongodb.Url = v
	})
	setEnvValue("TEAMSACS_MONGODB_USER", func(v string) {
		cfg.Mongodb.User = v
	})
	setEnvValue("TEAMSACS_MONGODB_PASSWD", func(v string) {
		cfg.Mongodb.Passwd = v
	})

	// Grpc Config
	setEnvValue("TEAMSACS_GRPC_HOST", func(v string) {
		cfg.Grpc.Host = v
	})
	setEnvInt64Value("TEAMSACS_GRPC_PORT", func(v int64) {
		cfg.Grpc.Port = int(v)
	})

	setEnvValue("TEAMSACS_GRPC_DEBUG", func(v string) {
		cfg.Grpc.Debug = v == "true"
	})

	// Radius config
	setEnvInt64Value("TEAMSACS_RADIUS_AUTH_PORT", func(v int64) {
		cfg.Radiusd.AuthPort = int(v)
	})

	setEnvInt64Value("TEAMSACS_RADIUS_ACCT_PORT", func(v int64) {
		cfg.Radiusd.AcctPort = int(v)
	})

	setEnvValue("TEAMSACS_RADIUS_DEBUG", func(v string) {
		cfg.Radiusd.Debug = v == "true"
	})

	// Elastic config
	setEnvValue("TEAMSACS_ELASTIC_URLS", func(v string) {
		cfg.Elastic.Urls = v
	})
	setEnvValue("TEAMSACS_ELASTIC_USER", func(v string) {
		cfg.Elastic.User = v
	})
	setEnvValue("TEAMSACS_ELASTIC_PWD", func(v string) {
		cfg.Elastic.Pwd = v
	})

	setEnvValue("TEAMSACS_ELASTIC_DEBUG", func(v string) {
		cfg.Elastic.Debug = v == "true"
	})

	// AzureStorage Config
	setEnvValue("TEAMSACS_AZURE_STORAGE_ACCOUNT", func(v string) {
		cfg.AzureStorage.AccountName = v
	})
	setEnvValue("TEAMSACS_AZURE_STORAGE_ACCESS_KEY", func(v string) {
		cfg.AzureStorage.AccountKey = v
	})

	return cfg
}
