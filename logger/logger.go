package logger

import "fmt"

type logData struct {
	message  string
	timeStr  string
	levelStr string
	fileName string
	funcName string
	lineNo   int
}

// 全局日志对象
var logger Log

// console 初始化一个console日志实例
// file 初始化一个文件的日志实例
func InitLogger(name string, config map[string]string) error {
	var err error
	switch name {
	case "file":
		logger, err = NewFileLogger(config)
	case "console":
		logger, err = NewConsoleLog(config)
	default:
		err = fmt.Errorf("UnSupport %s", name)
	}
	return err
}

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	logger.Trace(format, args...)
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	logger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	logger.Fatal(format, args...)
}
