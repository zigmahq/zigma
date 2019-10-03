package log

import (
	"os"
	"time"

	"github.com/go-stack/stack"
)

// Logger implements a simple logging interface
type Logger interface {
	ID() string
	NL()
	Trace(string, ...Field)
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Fatal(string, ...Field)
}

type logger struct {
	id string
}

func (l *logger) print(r *Record) {
	if r.Level > lvl {
		return
	}
	os.Stdout.WriteString(r.String())
}

func (l *logger) log(lv Level, msg string, fields ...Field) {
	l.print(&Record{
		Time:   time.Now().UTC(),
		Level:  lv,
		Msg:    msg,
		Fields: fields,
		Call:   stack.Caller(skipLevel),
	})
}

func (l *logger) ID() string {
	return l.id
}

func (l *logger) NL() {
	os.Stdout.WriteString("\n")
}

func (l *logger) Trace(msg string, fields ...Field) {
	l.log(LogTrace, msg, fields...)
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.log(LogDebug, msg, fields...)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.log(LogInfo, msg, fields...)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.log(LogWarn, msg, fields...)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.log(LogError, msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.log(LogFatal, msg, fields...)
}

// NewLogger initializes and returns a new logger
func NewLogger() Logger {
	return new(logger)
}
