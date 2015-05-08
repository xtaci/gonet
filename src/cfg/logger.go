package cfg

import (
	"log"
	"os"
	"strings"
)

//---------------------------------------------------------- 通用系统日志
func GetLogger(path string) *log.Logger {
	if !strings.HasPrefix(path, "/") {
		path = os.Getenv("GOPATH") + "/" + path
	}

	// 打开文件
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("error opening file %v\n", err)
		return nil
	}

	// 日志
	logger := log.New(file, "", log.LstdFlags)
	return logger
}

//---------------------------------------------------------- 同步系统日志
// 用于记录至关重要的日志数据(用 O_SYNC实现)
func GetSyncLogger(path string) *log.Logger {
	if !strings.HasPrefix(path, "/") {
		path = os.Getenv("GOPATH") + "/" + path
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		log.Println("error opening file %v\n", err)
		return nil
	}

	logger := log.New(file, "", log.LstdFlags)
	return logger
}

//---------------------------------------------------------- 默认系统日志
func StartLogger(path string) {
	if !strings.HasPrefix(path, "/") {
		path = os.Getenv("GOPATH") + "/" + path
	}

	// 打开日志文件
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("cannot open logfile %v\n", err)
	}

	// 创建MUX
	var r Repeater
	config := Get()
	switch config["log_output"] {
	case "both":
		r.out1 = os.Stdout
		r.out2 = file
	case "file":
		r.out2 = file
	}
	log.SetOutput(&r)
}
