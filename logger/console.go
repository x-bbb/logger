package logger

import (
	"fmt"
	"os"
)

type ConsoleLog struct {
	level   int
	logChan chan *logData
}

func NewConsoleLog(config map[string]string) (Log, error) {
	level, ok := config["level"]
	if !ok {
		err := fmt.Errorf("Not found level")
		return nil, err
	}
	logLevel := getLogLevel(level)

	logger := &ConsoleLog{level: logLevel}
	logger.init()
	go logger.persistence()
	return logger, nil

}

func (c *ConsoleLog) init() {
	c.logChan = make(chan *logData, 50000)
}

// 从chan中取日志，打印到终端
func (c *ConsoleLog) persistence() {
	for {
		select {
		case data := <-c.logChan:
			//fmt.Fprintf(os.Stdout, data.message)
			fmt.Fprintf(os.Stdout, "%s [%s] %s %s:%d %s\n", data.timeStr, data.levelStr, data.fileName,
				data.funcName, data.lineNo, data.message)
		}
	}
}

// 设置日志级别
func (c *ConsoleLog) SetLevel(level int) {
	if level < LogLevelInfo || level > LogLevelFatal {
		level = LogLevelInfo
	}
	c.level = level
}

func (c *ConsoleLog) Debug(format string, args ...interface{}) {
	// 如果日志级别大于传入的level,直接返回，不写日志
	if c.level > LogLevelDebug {
		return
	}
	writeLog(c.logChan, LogLevelDebug, format, args...)
}

func (c *ConsoleLog) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}
	writeLog(c.logChan, LogLevelTrace, format, args...)
}

func (c *ConsoleLog) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}
	writeLog(c.logChan, LogLevelInfo, format, args...)
}

func (c *ConsoleLog) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}
	writeLog(c.logChan, LogLevelWarn, format, args...)
}

func (c *ConsoleLog) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}
	writeLog(c.logChan, LogLevelError, format, args...)
}

func (c *ConsoleLog) Fatal(format string, args ...interface{}) {
	writeLog(c.logChan, LogLevelFatal, format, args...)
}

func (c *ConsoleLog) Close() {
	return
}
