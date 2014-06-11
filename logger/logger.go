package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	Level int
}

const (
	LogLevelQuiet = iota
	LogLevelUnknown
	LogLevelFatal
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

func (l *Logger) output(level int, tag string, v ...interface{}) {
	if l.Level >= level {
		w := append([]interface{}{tag}, v...)
		log.Println(w...)
	}
}

func (l *Logger) outputf(level int, tag string, format string, v ...interface{}) {
	if l.Level >= level {
		log.Printf(fmt.Sprintf("%v %v\n", tag, format), v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.output(LogLevelDebug, "[DEBUG]", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.outputf(LogLevelDebug, "[DEBUG]", format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.output(LogLevelInfo, "[INFO]", v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.outputf(LogLevelInfo, "[INFO]", format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.output(LogLevelWarn, "[WARN]", v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.outputf(LogLevelWarn, "[WARN]", format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.output(LogLevelError, "[ERROR]", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.outputf(LogLevelError, "[ERROR]", format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.output(LogLevelFatal, "[FATAL]", v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.outputf(LogLevelFatal, "[FATAL]", format, v...)
}

func (l *Logger) Unknown(v ...interface{}) {
	l.output(LogLevelUnknown, "[UNKNOWN]", v...)
}

func (l *Logger) Unknownf(format string, v ...interface{}) {
	l.outputf(LogLevelUnknown, "[UNKNOWN]", format, v...)
}

func (l *Logger) FatalExit(v ...interface{}) {
	l.Fatal(v...)
	os.Exit(1)
}

func (l *Logger) FatalfExit(format string, v ...interface{}) {
	l.Fatalf(format, v...)
	os.Exit(1)
}

func (l *Logger) TimeInfo(start time.Time, v ...interface{}) {
	elapsed := time.Since(start)
	l.Info(append(v, elapsed)...)
}
