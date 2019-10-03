package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zigmahq/zigma/config"
)

var cfg = []byte(`
p2p: {}
`)

func TestRead(t *testing.T) {
	c, err := config.Read(cfg)

	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.P2P)
}

func TestFromFile(t *testing.T) {
	f, err := ioutil.TempFile("", "cfg")

	assert.Nil(t, err)
	assert.NotNil(t, f)
	assert.NotEmpty(t, f.Name())
	defer os.Remove(f.Name())

	f.Write(cfg)

	c, err := config.FromFile(f.Name())
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.P2P)
}
