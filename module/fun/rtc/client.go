// @Author: wangsheng
// @Description:
// @File:  client
// @Version: 1.0.0
// @Date: 2022/5/25 10:17 上午
// Copyright 2021 QINIU. All rights reserved

package rtc

type Service interface {
	GetRoomToken(userId, roomId string) string

	GetRelayToken(userId, roomId string) string

	ListUser(roomId string) (res []string, err error)

	Online(userId, roomId string) bool
}

var service Service

func GetService() Service {
	return service
}

func InitService(conf Config) {
	service = NewQiniuClient(conf)
}
