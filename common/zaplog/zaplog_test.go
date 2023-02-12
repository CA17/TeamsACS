package zaplog

import (
	"fmt"
	"testing"
	"time"

	"github.com/ca17/teamsacs/common"
	"github.com/ca17/teamsacs/common/zaplog/log"
)

func TestInfo(t *testing.T) {
	log.Info("test")
}

func TestLoki(t *testing.T) {
	InitGlobalLogger(LogConfig{
		Mode:       Dev,
		LokiEnable: false,
		FileEnable: true,
		Filename:   "/tmp/test.log",
		LokiApi:    "http://127.0.0.1:3100",
		LokiUser:   "test",
		LokiPwd:    "test",
		LokiJob:    "test",
	})
	for i := 0; i < 32; i++ {
		log.Info(fmt.Sprintf("hello world %s", common.UUID()))
		log.Warn(fmt.Sprintf("hello world %s", common.UUID()))
		log.Error(fmt.Sprintf("hello world %s", common.UUID()))
		time.Sleep(time.Millisecond * 10)
	}
	time.Sleep(time.Second * 10)
}
