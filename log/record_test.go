package log_test

import (
	"testing"
	"time"

	"github.com/go-stack/stack"
	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/log"
)

func TestRecordOutput(t *testing.T) {
	record := &log.Record{
		Time:  time.Date(2019, time.September, 21, 0, 0, 0, 0, time.UTC),
		Level: log.LogInfo,
		Msg:   "hello world",
		Fields: []log.Field{
			log.String("key", "val"),
		},
		Call: stack.Caller(2),
	}
	assert.Equal(t,
		"\x1b[92mINFO \x1b[0m\x1b[92m | \x1b[0m\x1b[30m\x1b[1m2019-09-21T00:00:00\x1b[0m \x1b[0mhello world\x1b[0m \x1b[0m\x1b[37m[\x1b[0m\x1b[37mkey\x1b[0m\x1b[37m: \x1b[0m\x1b[37mval\x1b[0m\x1b[37m] \x1b[0m\n",
		record.String())
}
