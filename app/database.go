package app

import (
	"fmt"
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/zaplog/log"
	"github.com/ca17/teamsacs/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	dbConnectAttempts = 30
	dbConnectDelay    = 2 * time.Second
)

// getPgDatabase opens PostgreSQL with retries so Docker Compose starts
// remain smooth while the database container is still initializing.
func getPgDatabase(config config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host,
		config.User,
		config.Passwd,
		config.Name,
		config.Port)

	var (
		pool *gorm.DB
		err  error
	)
	for attempt := 1; attempt <= dbConnectAttempts; attempt++ {
		pool, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			},
			Logger: logger.New(
				zap.NewStdLog(zap.L()), // io writer
				logger.Config{
					SlowThreshold:             time.Millisecond * 200,                                                // Slow SQL threshold
					LogLevel:                  common.If(config.Debug, logger.Info, logger.Silent).(logger.LogLevel), // Log level
					IgnoreRecordNotFoundError: true,                                                                  // Ignore ErrRecordNotFound error for logger
					Colorful:                  false,                                                                 // Disable color
				},
			),
		})
		if err == nil {
			sqlDB, dbErr := pool.DB()
			if dbErr != nil {
				err = dbErr
			} else if pingErr := sqlDB.Ping(); pingErr != nil {
				err = pingErr
			} else {
				sqlDB.SetMaxIdleConns(config.IdleConn)
				sqlDB.SetMaxOpenConns(config.MaxConn)
				if attempt > 1 {
					log.Infof("database ready after %d attempts", attempt)
				}
				return pool
			}
		}
		log.Warnf("waiting for database (%d/%d): %v", attempt, dbConnectAttempts, err)
		time.Sleep(dbConnectDelay)
	}
	common.Must(fmt.Errorf("database not ready after %d attempts: %w", dbConnectAttempts, err))
	return pool
}
