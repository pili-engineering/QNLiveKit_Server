package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	path := "config.yaml"
	path2 := "live.json"
	c, err := LoadConfig(path, path2)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestConfig_Sub(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		path  string
		path2 string
		args  args
	}{
		{"", "config.yaml", "live.json", args{key: "im"}},
		{"", "config.yaml", "live.json", args{key: "service"}},
	}
	for _, tt := range tests {
		c, _ := LoadConfig(tt.path, tt.path2)
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, c.Sub(tt.args.key), "Sub(%v)", tt.args.key)
		})
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	path := "config.yaml"
	path2 := "live.json"
	c, _ := LoadConfig(path, path2)

	sub := c.Sub("service")
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
