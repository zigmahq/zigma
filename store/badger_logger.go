package store

import (
	"fmt"
	"strings"

	"github.com/zigmahq/zigma/log"
)

var logger = log.DefaultLogger

// BadgerLogger implements logging interface for badger datastore
type BadgerLogger struct{}

func (b *BadgerLogger) format(f string, args ...interface{}) string {
	str := fmt.Sprintf(f, args...)
	str = strings.TrimSpace(str)
	return str
}

// Errorf logs an ERROR log message to the logger specified in opts or to the
// global logger if no logger is specified in opts.
func (b *BadgerLogger) Errorf(f string, args ...interface{}) {
	logger.Error(b.format(f, args...))
}

// Warningf logs a WARNING message to the logger specified in opts.
func (b *BadgerLogger) Warningf(f string, args ...interface{}) {
	logger.Warn(b.format(f, args...))
}

// Infof logs an INFO message to the logger specified in opts.
func (b *BadgerLogger) Infof(f string, args ...interface{}) {
	logger.Info(b.format(f, args...))
}

// Debugf logs a DEBUG message to the logger specified in opts.
func (b *BadgerLogger) Debugf(f string, args ...interface{}) {
	logger.Debug(b.format(f, args...))
}
