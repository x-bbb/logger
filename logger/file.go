package logger

import (
	"fmt"
	"log"
	"os"
)

type FileLog struct {
	level    int
	fileName string
	file     *os.File
	logChan  chan *logData
}

// 返回文件日志对象
func NewFileLogger(config map[string]string) (Log, error) {
	level, ok := config["level"]
	if !ok {
		err := fmt.Errorf("Not found Level")
		return nil, err
	}
	fileName, ok := config["filename"]
	if !ok {
		err := fmt.Errorf("Not found filename")
		return nil, err
	}

	LogLevel := getLogLevel(level)
	logger := &FileLog{
		level:    LogLevel,
		fileName: fileName,
	}
	logger.init()
	go logger.persistence()
	return logger, nil
}

func (f *FileLog) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

// 从chan中取出日志，持久化到文件中，防止写文件压力过大
func (f *FileLog) persistence() {
	for {
		select {
		case data := <-f.logChan:
			//fmt.Fprintf(f.file, message)
			fmt.Fprintf(f.file, "%s [%s] %s %s:%d %s\n", data.timeStr, data.levelStr,
				data.fileName, data.funcName, data.lineNo, data.message)
		}
	}
}

// 初始化文件句柄
func (f *FileLog) init() {
	file, err := os.OpenFile(f.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("open file err:", err)
	}

	f.logChan = make(chan *logData, 50000)
	f.file = file
}

func (f *FileLog) Debug(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	writeLog(f.logChan, LogLevelDebug, format, args...)
}

func (f *FileLog) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	writeLog(f.logChan, LogLevelTrace, format, args...)

}
func (f *FileLog) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	writeLog(f.logChan, LogLevelInfo, format, args...)
}

func (f *FileLog) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	writeLog(f.logChan, LogLevelWarn, format, args...)
}

func (f *FileLog) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	writeLog(f.logChan, LogLevelError, format, args...)
}

func (f *FileLog) Fatal(format string, args ...interface{}) {
	writeLog(f.logChan, LogLevelFatal, format, args...)
}

func (f *FileLog) Close() {
	f.file.Close()
}
