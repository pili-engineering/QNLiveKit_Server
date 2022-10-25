package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	path := "config.yaml"
	c, err := LoadConfig(path)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestConfig_Sub(t *testing.T) {
	path := "config.yaml"
	c, _ := LoadConfig(path)

	sub := c.Sub("cron_config")
	assert.NotNil(t, sub)
}

func TestConfig_Unmarshal(t *testing.T) {
	path := "config.yaml"
	c, _ := LoadConfig(path)

	sub := c.Sub("impl")
	assert.NotNil(t, sub)

	server := Server{}
	err := sub.Unmarshal(&server)
	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", server.Host)
	assert.Equal(t, 8099, server.Port)
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
