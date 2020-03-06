package logger

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type FileLog struct {
	level        int
	fileName     string
	file         *os.File
	logChan      chan *logData
	logSplitType string // hour size
	logSplitSize int64  // 当logSplitType为size的时候使用
	time         int
}

// 返回文件日志对象
func NewFileLogger(config map[string]string) (Log, error) {
	var logSplitSize int64

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

	logSplitType, ok := config["logSplitType"]
	// 没有传切分类型，默认使用大小切分
	if !ok {
		logSplitType = LogSplitType
		logSplitSize = LogSplitSize
	} else {
		if logSplitType == "size" {
			logSplitSizeStr, ok := config["logSplitSize"]
			if !ok {
				logSplitSize = LogSplitSize
			}

			logSplitSize, err := strconv.ParseInt(logSplitSizeStr, 10, 64)
			if err != nil {
				logSplitSize = LogSplitSize
			}
			logSplitSize = logSplitSize

		}
		if logSplitType != "time" && logSplitType != "size" {
			err := fmt.Errorf("Unknown logSplitType")
			return nil, err
		}
	}

	LogLevel := getLogLevel(level)
	logger := &FileLog{
		level:        LogLevel,
		fileName:     fileName,
		logSplitType: logSplitType,
		logSplitSize: logSplitSize,
		time:         time.Now().Hour(),
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

// 切分日志逻辑
func (f *FileLog) checkSplitFile() {

	now := time.Now()
	oldFileName := f.fileName
	BakFileName := fmt.Sprintf("%s.%04d-%02d-%02d-%02d-%02d-%02d", f.fileName, now.Year(), now.Month(),
		now.Day(), f.time, now.Minute(), now.Second())

	// 基于时间备份
	if f.logSplitType == "time" {
		if f.time == now.Hour() {
			return
		}
		f.file.Close()
		file, err := backupFile(oldFileName, BakFileName)
		if err != nil {
			fmt.Errorf("backup file error: %v\n", err)
			return
		}
		f.file = file
		f.time = now.Hour()
	}
	// 基于文件大小备份
	if f.logSplitType == "size" {
		//fmt.Println("xxaxxxaaxasasaaaaaaaaaaaaaaaaaa")
		fileStat, err := os.Stat(f.fileName)
		if err != nil {
			fmt.Errorf("stat file error:%v", err)
			return
		}
		if fileStat.Size() >= f.logSplitSize {
			f.file.Close()
			file, err := backupFile(oldFileName, BakFileName)
			if err != nil {
				fmt.Errorf("backup file err:%v\n", err)
				return
			}
			f.file = file
			f.time = now.Hour()
		}
	}

}

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
	f.checkSplitFile()
	writeLog(f.logChan, LogLevelDebug, format, args...)
}

func (f *FileLog) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	f.checkSplitFile()
	writeLog(f.logChan, LogLevelTrace, format, args...)

}
func (f *FileLog) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	f.checkSplitFile()
	writeLog(f.logChan, LogLevelInfo, format, args...)
}

func (f *FileLog) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}
	f.checkSplitFile()
	writeLog(f.logChan, LogLevelWarn, format, args...)
}

func (f *FileLog) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}
	f.checkSplitFile()
	writeLog(f.logChan, LogLevelError, format, args...)
}

func (f *FileLog) Fatal(format string, args ...interface{}) {
	f.checkSplitFile()
	writeLog(f.logChan, LogLevelFatal, format, args...)
}

func (f *FileLog) Close() {
	f.file.Close()
}
