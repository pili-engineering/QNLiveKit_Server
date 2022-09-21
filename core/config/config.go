package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %s error %v", path, err)
	}
	defer f.Close()

	v := viper.New()
	v.SetConfigType("yaml")
	err = v.ReadConfig(f)
	if err != nil {
		return nil, fmt.Errorf("read config file %s error %v", path, err)
	}
	return &Config{
		v,
	}, nil
}

func (c *Config) Sub(key string) *Config {
	if v := c.v.Sub(key); v == nil {
		return nil
	} else {
		return &Config{
			v,
		}
	}
}

func (c *Config) Unmarshal(val interface{}) error {
	return c.v.Unmarshal(val)
}
