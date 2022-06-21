// @Author: wangsheng
// @Description:
// @File:  config
// @Version: 1.0.0
// @Date: 2022/5/19 9:59 上午
// Copyright 2021 QINIU. All rights reserved

package config

import (
	"github.com/qbox/livekit/common/auth/qiniumac"
	"github.com/qbox/livekit/common/im"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/common/rtc"
	"github.com/qbox/livekit/utils/config"
)

// AppConfig global config
var AppConfig Config

type Config struct {
	NodeID int64                    `mapstructure:"node_id"`
	Server Server                   `mapstructure:"service"`
	JwtKey string                   `mapstructure:"jwt_key"`
	Mysqls []*mysql.ConfigStructure `mapstructure:"mysqls"`

	MacConfig qiniumac.Config `mapstructure:"mac_config"`
	RtcConfig rtc.Config      `mapstructure:"rtc_config"`
	ImConfig  im.Config       `mapstructure:"im_config"`
}

// Server service port and host
type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func LoadConfig(confPath string) error {
	return config.LoadConfig(confPath, &AppConfig)
}
