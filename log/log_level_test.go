package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/log"
)

func TestLevelHead(t *testing.T) {
	assert.Equal(t, "FATAL", log.LogFatal.Head())
	assert.Equal(t, "ERROR", log.LogError.Head())
	assert.Equal(t, "WARN ", log.LogWarn.Head())
	assert.Equal(t, "INFO ", log.LogInfo.Head())
	assert.Equal(t, "DEBUG", log.LogDebug.Head())
	assert.Equal(t, "TRACE", log.LogTrace.Head())
	assert.Panics(t, func() {
		_ = log.Level(9999).Head()
	})
}

func TestLevelString(t *testing.T) {
	assert.Equal(t, "fatal", log.LogFatal.String())
	assert.Equal(t, "error", log.LogError.String())
	assert.Equal(t, "warn", log.LogWarn.String())
	assert.Equal(t, "info", log.LogInfo.String())
	assert.Equal(t, "debug", log.LogDebug.String())
	assert.Equal(t, "trace", log.LogTrace.String())
	assert.Panics(t, func() {
		_ = log.Level(9999).String()
	})
}
