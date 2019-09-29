package config

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Config encapsulates all configuration details for zigma
type Config struct {
	path string
	P2P  *P2P `yaml:"p2p"`
}

// DefaultConfig generates the default settings for zigma
func DefaultConfig() *Config {
	return &Config{
		P2P: DefaultP2P(),
	}
}

// Read reads the yaml configuration from bytes
func Read(b []byte) (*Config, error) {
	config := DefaultConfig()
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// FromFile reads the yaml configuration from specify path
func FromFile(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	conf, err := Read(data)
	if err != nil {
		return nil, err
	}
	conf.path = path
	return conf, nil
}
