package log

import (
	"strings"
)

// Composer encapsulates a text writer
type Composer interface {
	Write(string)
	Bold(string) string
	Underline(string) string
	Blink(string) string
	Grey(string) string
	Red(string) string
	Green(string) string
	Yellow(string) string
	Blue(string) string
	LightGrey(string) string
	String() string
}

type composer struct {
	sb strings.Builder
}

func (c *composer) Write(s string) {
	c.sb.WriteString(s)
	c.sb.WriteString("\033[0m")
}

func (c *composer) styled(a, s string) string {
	return a + s
}

func (c *composer) Bold(s string) string {
	return c.styled("\033[1m", s)
}

func (c *composer) Underline(s string) string {
	return c.styled("\033[4m", s)
}

func (c *composer) Blink(s string) string {
	return c.styled("\033[5m", s)
}

func (c *composer) Grey(s string) string {
	return c.styled("\x1b[30m", s)
}

func (c *composer) Red(s string) string {
	return c.styled("\x1b[91m", s)
}

func (c *composer) Green(s string) string {
	return c.styled("\x1b[92m", s)
}

func (c *composer) Yellow(s string) string {
	return c.styled("\x1b[93m", s)
}

func (c *composer) Blue(s string) string {
	return c.styled("\x1b[94m", s)
}

func (c *composer) LightGrey(s string) string {
	return c.styled("\x1b[37m", s)
}

func (c *composer) String() string {
	return c.sb.String() + "\n"
}

// NewWriter returns a new text writer
func NewWriter() Composer {
	var sb strings.Builder
	return &composer{
		sb: sb,
	}
}
