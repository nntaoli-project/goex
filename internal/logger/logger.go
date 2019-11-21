package logger

import (
	"fmt"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota + 1
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

type Logger struct {
	*log.Logger
	level Level
}

var Log = NewLogger()

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		level:  INFO,
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) Debug(args ...interface{}) {
	if l.level <= DEBUG {
		var msg []interface{}
		msg = append(msg, "[D]")
		msg = append(msg, args...)
		l.Output(2, fmt.Sprint(msg...))
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level <= DEBUG {
		l.Output(2, fmt.Sprintf("[D] %s", fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.level <= INFO {
		var msg []interface{}
		msg = append(msg, "[I]")
		msg = append(msg, args...)
		l.Output(2, fmt.Sprint(msg...))
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.level <= INFO {
		l.Output(2, fmt.Sprintf("[I] %s", fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Warn(args ...interface{}) {
	if l.level <= WARN {
		var msg []interface{}
		msg = append(msg, "[W]")
		msg = append(msg, args...)
		l.Output(2, fmt.Sprint(msg...))
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.level <= WARN {
		l.Output(2, fmt.Sprintf("[W] %s", fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.level <= ERROR {
		var msg []interface{}
		msg = append(msg, "[E]")
		msg = append(msg, args...)
		l.Output(2, fmt.Sprint(msg...))
	}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.level <= ERROR {
		l.Output(2, fmt.Sprintf("[E] %s", fmt.Sprintf(format, args...)))
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.level <= FATAL {
		var msg []interface{}
		msg = append(msg, "[F]")
		msg = append(msg, args...)
		l.Output(2, fmt.Sprint(args...))
		os.Exit(1)
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.level <= FATAL {
		l.Output(2, fmt.Sprintf("[F] %s", fmt.Sprintf(format, args...)))
		os.Exit(1)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	if l.level <= PANIC {
		s := fmt.Sprint(args...)
		s = fmt.Sprintf("[P] %s", s)
		l.Output(2, s)
		panic(s)
	}
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.level <= PANIC {
		s := fmt.Sprintf("[P] %s", fmt.Sprintf(format, args...))
		l.Output(2, s)
		panic(s)
	}
}
