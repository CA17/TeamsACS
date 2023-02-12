package config

import (
	"os"
	"path"
	"strconv"

	"github.com/ca17/teamsacs/common"
	"gopkg.in/yaml.v3"
)

// DBConfig Database (PostgreSQL) configuration
type DBConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
	MaxConn  int    `yaml:"max_conn"`
	IdleConn int    `yaml:"idle_conn"`
	Debug    bool   `yaml:"debug"`
}

// SysConfig System Configuration
type SysConfig struct {
	Appid    string `yaml:"appid"`
	Location string `yaml:"location"`
	Workdir  string `yaml:"workdir"`
	Debug    bool   `yaml:"debug"`
}

// WebConfig WEB Configuration
type WebConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	TlsPort int    `yaml:"tls_port"`
	Secret  string `yaml:"secret"`
}

// Tr069Config tr069 API Configuration
type Tr069Config struct {
	Host   string `yaml:"host" json:"host"`
	Port   int    `yaml:"port" json:"port"`
	Tls    bool   `yaml:"tls" json:"tls"`
	Secret string `yaml:"secret" json:"secret"`
	Debug  bool   `yaml:"debug" json:"debug"`
}

type MqttConfig struct {
	Server   string `yaml:"server" json:"server"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Debug    bool   `yaml:"debug" json:"debug"`
}

type AppConfig struct {
	System   SysConfig   `yaml:"system" json:"system"`
	Web      WebConfig   `yaml:"web" json:"web"`
	Database DBConfig    `yaml:"database" json:"database"`
	Tr069    Tr069Config `yaml:"tr069" json:"tr069"`
	Mqtt     MqttConfig  `yaml:"mqtt" json:"mqtt"`
}

func (c *AppConfig) GetLogDir() string {
	return path.Join(c.System.Workdir, "logs")
}

func (c *AppConfig) GetPublicDir() string {
	return path.Join(c.System.Workdir, "public")
}

func (c *AppConfig) GetPrivateDir() string {
	return path.Join(c.System.Workdir, "private")
}

func (c *AppConfig) GetDataDir() string {
	return path.Join(c.System.Workdir, "data")
}
func (c *AppConfig) GetBackupDir() string {
	return path.Join(c.System.Workdir, "backup")
}

func (c *AppConfig) InitDirs() {
	os.MkdirAll(path.Join(c.System.Workdir, "logs"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "public"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "data"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "cwmp"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "data/metrics"), 0755)
	os.MkdirAll(path.Join(c.System.Workdir, "private"), 0644)
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
		Appid:    "TeamsACS",
		Location: "Asia/Shanghai",
		Workdir:  "/var/teamsacs",
		Debug:    true,
	},
	Web: WebConfig{
		Host:    "0.0.0.0",
		Port:    2979,
		TlsPort: 2989,
		Secret:  "9b6de5cc-0001-1203-1100-0f568ac7da37",
	},
	Database: DBConfig{
		Type:     "postgres",
		Host:     "127.0.0.1",
		Port:     5432,
		Name:     "teamsacs_v1",
		User:     "postgres",
		Passwd:   "myroot",
		MaxConn:  100,
		IdleConn: 10,
		Debug:    false,
	},
	Tr069: Tr069Config{
		Host:   "0.0.0.0",
		Tls:    true,
		Port:   2999,
		Secret: "9b6de5cc-1q21-1203-xxtt-0f568ac9d237",
		Debug:  true,
	},
	Mqtt: MqttConfig{
		Server:   "",
		Username: "",
		Password: "",
		Debug:    false,
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
		data := common.Must2(os.ReadFile(cfile))
		common.Must(yaml.Unmarshal(data.([]byte), cfg))
	} else {
		cfg = DefaultAppConfig
	}

	cfg.InitDirs()

	setEnvValue("TEAMSACS_SYSTEM_WORKER_DIR", &cfg.System.Workdir)
	setEnvBoolValue("TEAMSACS_SYSTEM_DEBUG", &cfg.System.Debug)

	// WEB
	setEnvValue("TEAMSACS_WEB_HOST", &cfg.Web.Host)
	setEnvValue("TEAMSACS_WEB_SECRET", &cfg.Web.Secret)
	setEnvIntValue("TEAMSACS_WEB_PORT", &cfg.Web.Port)
	setEnvIntValue("TEAMSACS_WEB_TLS_PORT", &cfg.Web.TlsPort)

	// DB
	setEnvValue("TEAMSACS_DB_HOST", &cfg.Database.Host)
	setEnvValue("TEAMSACS_DB_NAME", &cfg.Database.Name)
	setEnvValue("TEAMSACS_DB_USER", &cfg.Database.User)
	setEnvValue("TEAMSACS_DB_PWD", &cfg.Database.Passwd)
	setEnvIntValue("TEAMSACS_DB_PORT", &cfg.Database.Port)
	setEnvBoolValue("TEAMSACS_DB_DEBUG", &cfg.Database.Debug)

	// TR069 Config
	setEnvValue("TEAMSACS_TR069_WEB_HOST", &cfg.Tr069.Host)
	setEnvValue("TEAMSACS_TR069_WEB_SECRET", &cfg.Tr069.Secret)
	setEnvBoolValue("TEAMSACS_TR069_WEB_TLS", &cfg.Tr069.Tls)
	setEnvBoolValue("TEAMSACS_TR069_WEB_DEBUG", &cfg.Tr069.Debug)
	setEnvIntValue("TEAMSACS_TR069_WEB_PORT", &cfg.Tr069.Port)

	setEnvValue("TEAMSACS_MQTT_SERVER", &cfg.Mqtt.Server)
	setEnvValue("TEAMSACS_MQTT_USERNAME", &cfg.Mqtt.Username)
	setEnvValue("TEAMSACS_MQTT_PASSWORD", &cfg.Mqtt.Password)
	setEnvBoolValue("TEAMSACS_MQTT_DEBUG", &cfg.Mqtt.Debug)

	return cfg
}
