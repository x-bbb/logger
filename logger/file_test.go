package logger

import (
	"testing"
	"time"
)

func TestFileLog(t *testing.T) {
	err := InitLogger("file", map[string]string{"level": "debug", "filename": "1.txt"})
	if err != nil {
		t.Errorf("test failed, err:%v", err)
	}
	Warn("test Debug %s", "file not fount")
	time.Sleep(time.Second)
}
