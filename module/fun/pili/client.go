// @Author: wangsheng
// @Description:
// @File:  client
// @Version: 1.0.0
// @Date: 2022/5/25 10:17 上午
// Copyright 2021 QINIU. All rights reserved

package pili

import (
	"context"
	"time"
)

type SaveStreamRequest struct {
	Fname                     string `json:"fname"`
	Start                     int64  `json:"start"`
	End                       int64  `json:"end"`
	Pipeline                  string `json:"pipeline"`
	Format                    string `json:"uint"`
	ExpireDays                int    `json:"expireDays"`
	Notify                    string `json:"notify"`
	PersistentDeleteAfterDays int    `json:"persistentDeleteAfterDays"`
	FirstTsType               byte   `json:"firstTsType"`
}

type SaveStreamResponse struct {
	Fname        string `json:"fname"`
	Start        int64  `json:"start"`
	End          int64  `json:"end"`
	PersistentID string `json:"persistentID"`
}

type Service interface {
	PiliHub() string

	StreamPubURL(roomId string, expectAt *time.Time) (url string)

	StreamRtmpPlayURL(roomId string) (url string)

	StreamFlvPlayURL(roomId string) (url string)

	StreamHlsPlayURL(roomId string) (url string)

	PlaybackURL(fname string) string

	SaveStream(ctx context.Context, req *SaveStreamRequest, encodedStreamTitle string) (*SaveStreamResponse, error)
}

var service *QiniuClient

func GetService() Service {
	return service
}

func InitService(conf Config) {
	service = NewQiniuClient(conf)
}
