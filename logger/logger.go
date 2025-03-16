package logger

import (
	"fmt"
	"github.com/nntaoli/go-tools/logger"
	"io"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = iota + 1
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

type Logger struct {
	*logger.Logger
	level LogLevel
}

var std = &Logger{
	Logger: logger.NewLogger().WithLongFile(),
	level:  WARN}

func init() {
	std.SetPrefix("goex")
}

func SetOut(out io.Writer) {
	std.SetOut(out)
}

func SetLevel(level LogLevel) {
	std.level = level
	std.SetLevel(logger.Level(level))
}

func Debug(args ...interface{}) {
	std.Output(3, logger.DEBUG, "[DEBUG]", fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	std.Output(3, logger.DEBUG, "[DEBUG]", fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) {
	std.Output(3, logger.INFO, "[INFO ]", fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	std.Output(3, logger.INFO, "[INFO ]", fmt.Sprintf(format, args...))
}

func Warn(args ...interface{}) {
	std.Output(3, logger.WARN, "[WARN ]", fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	std.Output(3, logger.WARN, "[WARN ]", fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	std.Output(3, logger.ERROR, "[ERROR]", fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	std.Output(3, logger.ERROR, "[ERROR]", fmt.Sprintf(format, args...))
}

func Fatal(args ...interface{}) {
	if std.level <= FATAL {
		std.Output(3, logger.FATAL, "[FATAL]", fmt.Sprint(args...))
		os.Exit(1)
	}
}

func Fatalf(format string, args ...interface{}) {
	if std.level <= FATAL {
		std.Output(3, logger.FATAL, "[FATAL]", fmt.Sprintf(format, args...))
		os.Exit(1)
	}
}

func Panic(args ...interface{}) {
	if std.level <= PANIC {
		std.Output(3, logger.PANIC, "[PANIC]", fmt.Sprint(args...))
		panic("")
	}
}

func Panicf(format string, args ...interface{}) {
	if std.level <= PANIC {
		std.Output(3, logger.PANIC, "[PANIC]", fmt.Sprintf(format, args...))
		panic("")
	}
}
