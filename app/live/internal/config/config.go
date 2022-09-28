// @Author: wangsheng
// @Description:
// @File:  config
// @Version: 1.0.0
// @Date: 2022/5/19 9:59 上午
// Copyright 2021 QINIU. All rights reserved

package config

import (
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/fun/im"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/module/store/cache"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/config"
)

// AppConfig global config
var AppConfig Config

type Config struct {
	NodeID         int64                    `mapstructure:"node_id"`
	Server         Server                   `mapstructure:"impl"`
	JwtKey         string                   `mapstructure:"jwt_key"`
	CensorCallback string                   `mapstructure:"censor_callback"`
	CensorBucket   string                   `mapstructure:"censor_bucket"`
	Callback       string                   `mapstructure:"callback"`
	ReportHost     string                   `mapstructure:"report_host"`
	Mysqls         []*mysql.ConfigStructure `mapstructure:"mysqls"`

	MacConfig   auth.Config  `mapstructure:"mac_config"`
	RtcConfig   rtc.Config   `mapstructure:"rtc_config"`
	ImConfig    im.Config    `mapstructure:"im_config"`
	CacheConfig cache.Config `mapstructure:"cache_config"`
}

// Server impl port and host
type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func LoadConfig(confPath string) error {
	return config.LoadConfig(confPath, &AppConfig)
}
