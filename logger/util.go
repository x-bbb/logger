package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

// 返回错误等级字符串
func getLogLevelText(logLevel int) string {
	switch logLevel {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	}
	return "UNKNOWN"
}

// 获取报错文件和行号
func getLineInfo() (fileName, funcName string, lineNo int) {
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = file
		funcName = runtime.FuncForPC(pc).Name()
		lineNo = line
	}

	return
}

func getLogLevel(level string) int {
	switch level {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	}
	return LogLevelDebug
}

// 写日志通用函数, 将日志写入到chan中，后台有持久化进程
func writeLog(ch chan<- *logData, level int, format string, args ...interface{}) {

	levelStr := getLogLevelText(level)
	// 获取报错点文件名，函数名，行号
	fileName, funcName, lineNo := getLineInfo()
	//去掉目录，只留文件名
	fileName = path.Base(fileName)
	funcName = path.Base(funcName)
	// 格式化时间
	currentTimeStr := time.Now().Format("2006-01-02 15:04:05")

	msg := fmt.Sprintf(format, args...)

	// 日志结构体
	data := &logData{
		fileName: fileName,
		funcName: funcName,
		message:  msg,
		timeStr:  currentTimeStr,
		levelStr: levelStr,
		lineNo:   lineNo,
	}

	//msg = fmt.Sprintf("%s [%s] %s %s:%d %s\n", currentTimeStr, levelStr, fileName, funcName, lineNo, msg)
	select {
	case ch <- data:
	default:
	}

}
