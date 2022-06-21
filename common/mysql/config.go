package mysql

import (
	"fmt"
	"net/url"
)

// ConfigStructure mysql config
type ConfigStructure struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifeTime int    `mapstructure:"conn_max_life_time"`
	Default         string `mapstructure:"default"`
	ReadOnly        bool   `mapstructure:"read_only"`
}

// GetHost get mysql host
func (conf *ConfigStructure) GetHost() string {

	return conf.Host
}

// GetPort get mysql port
func (conf *ConfigStructure) GetPort() int {

	return conf.Port
}

// GetAddress get mysql address
func (conf *ConfigStructure) GetAddress() string {

	return fmt.Sprintf("%s:%d", conf.Host, conf.Port)
}

// String mysql config to string
func (conf *ConfigStructure) String() string {

	return fmt.Sprintf("[Mysql]\nhost = %s\nport = %d\nuserName = %s\ndatabase = %s\nmaxIdleConns = %d\nmaxOpenConns = %d\nconnMaxLifeTime = %d",
		conf.Host, conf.Port, conf.Username, conf.Database, conf.MaxIdleConns, conf.MaxOpenConns, conf.ConnMaxLifeTime)
}

// GetURL get mysql url
func (conf *ConfigStructure) GetURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=%s&parseTime=true",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		url.QueryEscape("Asia/Shanghai"))
}
