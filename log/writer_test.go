package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/log"
)

func TestWrite(t *testing.T) {
	writer := log.NewWriter()
	writer.Write("test")
	assert.Equal(t, "test\033[0m\n", writer.String())
}

func TestBold(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\033[1mtest", writer.Bold("test"))
}

func TestUnderline(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\033[4mtest", writer.Underline("test"))
}

func TestBlink(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\033[5mtest", writer.Blink("test"))
}

func TestGrey(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\x1b[30mtest", writer.Grey("test"))
}

func TestRed(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\x1b[91mtest", writer.Red("test"))
}

func TestGreen(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\x1b[92mtest", writer.Green("test"))
}

func TestYellow(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\x1b[93mtest", writer.Yellow("test"))
}

func TestBlue(t *testing.T) {
	writer := log.NewWriter()
	assert.Equal(t, "\x1b[94mtest", writer.Blue("test"))
}

func TestChainStyles(t *testing.T) {
	writer := log.NewWriter()
	styled := writer.Bold(writer.Red("test"))
	assert.Equal(t, "\033[1m\x1b[91mtest", styled)
}

func TestChainWrite(t *testing.T) {
	writer := log.NewWriter()
	writer.Write(writer.Bold(writer.Red("1")))
	writer.Write(writer.Underline(writer.Red("2")))
	assert.Equal(t, "\x1b[1m\x1b[91m1\x1b[0m\x1b[4m\x1b[91m2\x1b[0m\n", writer.String())
}
