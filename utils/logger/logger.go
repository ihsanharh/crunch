package logger

import (
	"log"
	"os"
	"time"

	"github.com/XzFrosT/crunch/utils"
)

var (
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
)

func Fatal(v ...interface{}) {
	print(stderr, "\u001b[31m", "FATAL", v...)
	os.Exit(1)
}

func Error(v ...interface{}) {
	print(stderr, "\u001b[31m", "ERROR", v...)
}

func Warn(v ...interface{}) {
	print(stderr, "\u001b[33m", "WARN", v...)
}

func Info(v ...interface{}) {
	print(stdout, "\u001b[32m", "INFO", v...)
}

func Debug(v ...interface{}) {
	print(stdout, "\u001b[35m", "DEBUG", v...)
}

func FatalF(text string, args ...interface{}) {
	Fatal(utils.Fmt(text, args...))
}

func ErrorF(text string, args ...interface{}) {
	Error(utils.Fmt(text, args...))
}

func WarnF(text string, args ...interface{}) {
	Warn(utils.Fmt(text, args...))
}

func InfoF(text string, args ...interface{}) {
	Info(utils.Fmt(text, args...))
}

func DebugF(text string, args ...interface{}) {
	Debug(utils.Fmt(text, args...))
}

func print(std *log.Logger, color string, tag string, v ...interface{}) {
	for _, line := range v {
		std.Println(utils.Fmt("[%s] %s[%s]\u001b[0m %+v", time.Now().Format("02/01/2006 - 15:04:05"), color, tag, line))
	}
}