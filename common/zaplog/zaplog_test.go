package zaplog

import (
	"testing"

	"github.com/ca17/teamsacs/common/zaplog/log"
)

func TestInfo(t *testing.T) {
	log.Info("test")
}
