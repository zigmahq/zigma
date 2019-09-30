package log

import (
	"fmt"
	"strings"
)

// Level defines the importance and urgency of the log message
type Level int

// Defines the importance level of logs
const (
	LogFatal Level = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
	LogTrace
)

var logLevelRefs = map[Level]string{
	LogFatal: "fatal",
	LogError: "error",
	LogWarn:  "warn",
	LogInfo:  "info",
	LogDebug: "debug",
	LogTrace: "trace",
}

var logLevelIds = map[string]Level{
	"fatal": LogFatal,
	"error": LogError,
	"warn":  LogWarn,
	"info":  LogInfo,
	"debug": LogDebug,
	"trace": LogTrace,
}

// Head pads log level reference to 5 characters, this would
// left-justifies the strings and add spaces to fill the empty
// columns on the right.
func (l Level) Head() string {
	if s, ok := logLevelRefs[l]; ok {
		return fmt.Sprintf("%-5s", strings.ToUpper(s))
	}
	panic("unrecognized log level")
}

func (l Level) String() string {
	if s, ok := logLevelRefs[l]; ok {
		return s
	}
	panic("unrecognized log level")
}

// LevelFromString returns the log level enum from a string
func LevelFromString(s string) Level {
	if l, ok := logLevelIds[s]; ok {
		return l
	}
	panic("unrecognized log level")
}
