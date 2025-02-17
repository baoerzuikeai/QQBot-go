package logs

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

func Info(v ...any) {
	_, file, line, _ := runtime.Caller(1) // 获取调用日志函数的文件和行号
	log.SetFlags(0)
	fileName := filepath.Base(file)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	log.SetPrefix(fmt.Sprintf("[Info] %s %s:%d ", currentTime, fileName, line)) // 在日志前缀中加入文件位置
	log.Println(v...)
}

func Debug(v ...any) {
	_, file, line, _ := runtime.Caller(1) // 获取调用日志函数的文件和行号
	log.SetFlags(0)
	fileName := filepath.Base(file)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	log.SetPrefix(fmt.Sprintf("[Debug] %s %s:%d ", currentTime, fileName, line)) // 在日志前缀中加入文件位置
	log.Println(v...)
}
