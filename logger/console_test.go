package logger

import (
	"testing"
)

func TestConsoleLog(t *testing.T) {
	err := InitLogger("console", map[string]string{"level": "debug"})
	if err != nil {
		t.Errorf("init logger failed, err:%v", err)
	}

	Debug("xxxx")
}
