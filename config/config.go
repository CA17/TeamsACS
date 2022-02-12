package config

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/ca17/teamsacs/common"
	"gopkg.in/yaml.v2"
)

// DBConfig 数据库(PostgreSQL)配置
type DBConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Name   string `yaml:"name"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	Debug  bool   `yaml:"debug"`
}

// SysConfig 系统配置
type SysConfig struct {
	Appid      string `yaml:"appid"`
	Location   string `yaml:"location"`
	Workdir    string `yaml:"workdir"`
	SyslogAddr string `yaml:"syslog_addr"`
	Version    string `yaml:"version"`
	JobEnabled bool   `yaml:"job_enabled"`
	Debug      bool   `yaml:"debug"`
}

// WebConfig WEB 配置
type WebConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Tls    bool   `yaml:"tls"`
	Secret string `yaml:"secret"`
}

// CwmpConfig Cwmp API 配置
type CwmpConfig struct {
	Host  string `yaml:"host" json:"host"`
	Port  int    `yaml:"port" json:"port"`
	Tls   bool   `yaml:"tls" json:"tls"`
	Debug bool   `yaml:"debug" json:"debug"`
}

type AppConfig struct {
	System   SysConfig  `yaml:"system"`
	Web      WebConfig  `yaml:"web"`
	Database DBConfig   `yaml:"database"`
	Cwmp     CwmpConfig `yaml:"cwmp" json:"cwmp"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.System.Workdir, "logs")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.System.Workdir, "data")
}

func (c *AppConfig) GetPublicDir() string {
	return path.Join(c.System.Workdir, "public")
}

func (c *AppConfig) GetPrivateDir() string {
	return path.Join(c.System.Workdir, "private")
}

func (c *AppConfig) GetBackupDir() string {
	return path.Join(c.System.Workdir, "backup")
}

func (c *AppConfig) InitDirs() {
	os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "data"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "cwmpfile"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "public"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "private"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "backup"), 0644)
}

func setEnvValue(name string, val *string) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue
	}
}

func setEnvBoolValue(name string, val *bool) {
	var evalue = os.Getenv(name)
	if evalue != "" {
		*val = evalue == "true" || evalue == "1" || evalue == "on"
	}
}

func setEnvInt64Value(name string, val *int64) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = p
	}
}
func setEnvIntValue(name string, val *int) {
	var evalue = os.Getenv(name)
	if evalue == "" {
		return
	}

	p, err := strconv.ParseInt(evalue, 10, 64)
	if err == nil {
		*val = int(p)
	}
}

var DefaultAppConfig = &AppConfig{
	System: SysConfig{
		Appid:      "MetasLink",
		Location:   "Asia/Shanghai",
		Workdir:    "/var/teamsacs",
		SyslogAddr: "",
		Version:    "latest",
		JobEnabled: true,
		Debug:      true,
	},
	Web: WebConfig{
		Host:   "0.0.0.0",
		Port:   8000,
		Tls:    false,
		Secret: "9b6de5cc-0731-4bf1-xxxx-0f568ac9da37",
	},
	Database: DBConfig{
		Host:   "127.0.0.1",
		Port:   5432,
		Name:   "teamsacs",
		User:   "postgres",
		Passwd: "root",
		Debug:  false,
	},
	Cwmp: CwmpConfig{
		Host:  "0.0.0.0",
		Tls:   false,
		Port:  8106,
		Debug: true,
	},
}

func LoadConfig(cfile string) *AppConfig {
	// 开发环境首先查找当前目录是否存在自定义配置文件
	if cfile == "" {
		cfile = "teamsacs.yml"
	}
	if !common.FileExists(cfile) {
		cfile = "/etc/teamsacs.yml"
	}
	cfg := new(AppConfig)
	if common.FileExists(cfile) {
		data := common.Must2(ioutil.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultAppConfig
	}

	cfg.InitDirs()

	setEnvValue("TEAMSACS_SYSTEM_WORKER_DIR", &cfg.System.Workdir)
	setEnvValue("TEAMSACS_SYSLOG_HOST", &cfg.System.SyslogAddr)
	setEnvBoolValue("TEAMSACS_SYSTEM_DEBUG", &cfg.System.Debug)

	// WEB
	setEnvValue("TEAMSACS_WEB_HOST", &cfg.Web.Host)
	setEnvValue("TEAMSACS_WEB_SECRET", &cfg.Web.Secret)
	setEnvIntValue("TEAMSACS_WEB_PORT", &cfg.Web.Port)
	setEnvBoolValue("TEAMSACS_WEB_TLS", &cfg.Web.Tls)

	// DB
	setEnvValue("TEAMSACS_DB_HOST", &cfg.Database.Host)
	setEnvValue("TEAMSACS_DB_NAME", &cfg.Database.Name)
	setEnvValue("TEAMSACS_DB_USER", &cfg.Database.User)
	setEnvValue("TEAMSACS_DB_PWD", &cfg.Database.Passwd)
	setEnvIntValue("TEAMSACS_DB_PORT", &cfg.Database.Port)
	setEnvBoolValue("TEAMSACS_DB_DEBUG", &cfg.Database.Debug)

	// Cwmp Config
	setEnvValue("TEAMSACS_CWMP_WEB_HOST", &cfg.Cwmp.Host)
	setEnvBoolValue("TEAMSACS_CWMP_WEB_DEBUG", &cfg.Cwmp.Debug)
	setEnvIntValue("TEAMSACS_CWMP_WEB_PORT", &cfg.Cwmp.Port)

	return cfg
}
