package log

import (
	"time"

	"github.com/go-stack/stack"
)

const skipLevel = 2

// Record encapsulates a log event
type Record struct {
	Time   time.Time  `json:"time"`
	Level  Level      `json:"level"`
	Msg    string     `json:"message"`
	Fields []Field    `json:"fields"`
	Call   stack.Call `json:"call"`
}

func (r *Record) String() string {
	c := NewWriter()
	c.Write(c.Green(r.Level.Head()))
	c.Write(c.Green(" | "))
	c.Write(c.Grey(c.Bold(r.Time.Format("2006-01-02T15:04:05"))))
	c.Write(" ")
	c.Write(r.Msg)
	c.Write(" ")
	for _, field := range r.Fields {
		c.Write(c.LightGrey("["))
		c.Write(c.LightGrey(field.Key))
		c.Write(c.LightGrey(": "))
		switch s := field.Value(); {
		case len(s) > 20:
			c.Write(c.LightGrey(s + ".."))
		default:
			c.Write(c.LightGrey(s))
		}
		c.Write(c.LightGrey("] "))
	}
	return c.String()
}
