package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qbox/livekit/utils/logger"
)

func TestLoadConfig(t *testing.T) {
	path := "config.yaml"
	ctx := context.Background()
	log := logger.ReqLogger(ctx)
	c, err := LoadConfig(log, path)
	assert.Nil(t, err)
	assert.NotNil(t, c)
}

func TestConfig_Sub(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		path string
		args args
	}{
		{"", "config.yaml", args{key: "cron_config"}},
		{"", "config.yaml", args{key: "httpq"}},
	}
	for _, tt := range tests {
		ctx := context.Background()
		log := logger.ReqLogger(ctx)
		c, _ := LoadConfig(log, tt.path)
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, c.Sub(tt.args.key), "Sub(%v)", tt.args.key)
		})
	}
}

func TestConfig_Unmarshal(t *testing.T) {
	path := "config.yaml"
	ctx := context.Background()
	log := logger.ReqLogger(ctx)
	c, _ := LoadConfig(log, path)

	sub := c.Sub("httpq")
	assert.NotNil(t, sub)

	server := Server{}
	err := sub.Unmarshal(&server)
	assert.Nil(t, err)

	assert.Equal(t, ":8080", server.Addr)
}

type Server struct {
	Addr string `mapstructure:"addr"`
}
