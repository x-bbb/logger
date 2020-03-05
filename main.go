package main

import (
	"goProject/logger"
	"log"
	"time"
)

// console 打印日志到终端，file打印日志到文件
func init() {
	err := logger.InitLogger("console", map[string]string{"level": "debug", "filename": "1.txt"})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for {
		logger.Warn("server is running...")
		time.Sleep(time.Second)
	}
}
