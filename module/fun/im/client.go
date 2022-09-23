// @Author: wangsheng
// @Description:
// @File:  service
// @Version: 1.0.0
// @Date: 2022/5/20 3:31 下午
// Copyright 2021 QINIU. All rights reserved

package im

import (
	"context"

	"github.com/qbox/livekit/module/fun/im/maxim"
)

type Service interface {
	RegisterUser(ctx context.Context, username, password string) (int64, error)

	GetUserId(ctx context.Context, username string) (int64, error)

	UpdateUserPassword(ctx context.Context, userId int64, password string) (bool, error)

	// CreateChatroom owner 群主的IM 用户ID
	// name 群名
	CreateChatroom(ctx context.Context, owner int64, name string) (int64, error)

	SendCommandMessageToGroup(ctx context.Context, fromUserId int64, toGroupId int64, content string) error

	SendCommandMessageToUser(ctx context.Context, fromUserId int64, toUserId int64, content string) error
}

var service Service

func InitService(conf Config) {
	service = maxim.NewMaximClient(conf.AppId, conf.Token, conf.Endpoint)
}

func GetService() Service {
	return service
}
