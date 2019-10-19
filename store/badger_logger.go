/* Copyright 2019 zigma authors
 * This file is part of the zigma library.
 *
 * The zigma library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The zigma library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with the zigma library. If not, see <http://www.gnu.org/licenses/>.
 */

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
