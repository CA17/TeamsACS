package app

import (
	"fmt"
	slog "log"
	"os"
	"runtime/debug"
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/log"
	"github.com/ca17/teamsacs/config"
	"github.com/ca17/teamsacs/models"
	"github.com/op/go-logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	Config    *config.AppConfig
	GormDB    *gorm.DB
	CwmpTable *CwmpEventTable
)

// Init 全局初始化调用
func Init(cfg *config.AppConfig) {
	Config = cfg
	setupLogger()
	GormDB = getPgDatabase(cfg.Database)
	CwmpTable = NewCwmpEventTable()
	setupEventSubscribe()
	if Config.System.JobEnabled {
		setupJobs()
	}
	go Migrate(false)
}

// 初始化日志
func setupLogger() {
	level := logging.INFO
	if Config.System.Debug {
		level = logging.DEBUG
	}
	log.SetupLog(level, Config.System.SyslogAddr, Config.GetLogDir(), Config.System.Appid)
}

// getPgDatabase 获取数据库连接，执行一次
func getPgDatabase(config config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Passwd,
		config.Name,
		config.Port)
	pool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: logger.New(
			slog.New(os.Stdout, "\r\n", slog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,                                                           // Slow SQL threshold
				LogLevel:                  common.If(config.Debug, logger.Info, logger.Silent).(logger.LogLevel), // Log level
				IgnoreRecordNotFoundError: true,                                                                  // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                                                                  // Disable color
			},
		),
	})
	common.Must(err)
	return pool
}

func Migrate(track bool) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			if os.Getenv("GO_DEGUB_TRACE") != "" {
				debug.PrintStack()
			}
			err2, ok := err1.(error)
			if ok {
				err = err2
			}
		}
	}()
	if track {
		return GormDB.Debug().Migrator().AutoMigrate(models.Tables...)
	}
	return GormDB.Migrator().AutoMigrate(models.Tables...)
}

func DropAll() {
	_ = GormDB.Migrator().DropTable(models.Tables...)
}

func InitDb() {
	_ = GormDB.Migrator().DropTable(models.Tables...)
	_ = GormDB.Migrator().AutoMigrate(models.Tables...)
	initTimescaleTable()
	initSettings()
}

func GetSettingsStringValue(stype string, name string) string {
	var value string
	GormDB.Raw("SELECT value FROM sys_config WHERE type = ? and name = ? limit 1", stype, name).Scan(&value)
	return value
}

func GetCwmpSettingsStringValue(name string) string {
	var value string
	GormDB.Raw("SELECT value FROM sys_config WHERE type = ? and name = ? limit 1", "cwmp", name).Scan(&value)
	return value
}
