// @Author: wangsheng
// @Description:
// @File:  client
// @Version: 1.0.0
// @Date: 2022/5/25 10:17 上午
// Copyright 2021 QINIU. All rights reserved

package pili

import "time"

type Service interface {
	PiliHub() string

	StreamPubURL(roomId string, expectAt *time.Time) (url string)

	StreamRtmpPlayURL(roomId string) (url string)

	StreamFlvPlayURL(roomId string) (url string)

	StreamHlsPlayURL(roomId string) (url string)
}

var service Service

func GetService() Service {
	return service
}

func InitService(conf Config) {
	service = NewQiniuClient(conf)
}
