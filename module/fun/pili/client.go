// @Author: wangsheng
// @Description:
// @File:  client
// @Version: 1.0.0
// @Date: 2022/5/25 10:17 上午
// Copyright 2021 QINIU. All rights reserved

package pili

import "time"

type Service interface {
	StreamPubURL(roomId string, expectAt *time.Time) (url string)

	StreamRtmpPlayURL(roomId string) (url string)

	StreamFlvPlayURL(roomId string) (url string)

	StreamHlsPlayURL(roomId string) (url string)
}

type Config struct {
	AppId     string `mapstructure:"app_id"`     // RTC AppId
	AccessKey string `mapstructure:"access_key"` // AK
	SecretKey string `mapstructure:"secret_key"` // SK

	RoomTokenExpireS int64  `mapstructure:"room_token_expire_s"`
	RtcPlayBackUrl   string `mapstructure:"playback_url"`
	Hub              string `mapstructure:"hub"`
	StreamPattern    string `mapstructure:"stream_pattern"`
	PublishUrl       string `mapstructure:"publish_url"`
	PublishDomain    string `mapstructure:"publish_domain"`
	RtmpPlayUrl      string `mapstructure:"rtmp_play_url"`
	FlvPlayUrl       string `mapstructure:"flv_play_url"`
	HlsPlayUrl       string `mapstructure:"hls_play_url"`
	SecurityType     string `mapstructure:"security_type"`    //expiry, expiry_sk, none
	PublishKey       string `mapstructure:"publish_key"`      //推流key
	PublishExpireS   int64  `mapstructure:"publish_expire_s"` //推流URL 过期时间，单位：秒
}

var service Service

func GetService() Service {
	return service
}

func InitService(conf Config) {
	service = NewQiniuClient(conf)
}
