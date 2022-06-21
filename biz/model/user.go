// @Author: wangsheng
// @Description:
// @File:  user
// @Version: 1.0.0
// @Date: 2022/5/19 3:24 下午
// Copyright 2021 QINIU. All rights reserved

package model

import (
	"github.com/qbox/livekit/utils/timestamp"
)

type LiveUserEntity struct {
	ID     uint   `gorm:"primary_key"`
	UserId string `json:"user_id"` //应用内用户ID

	Nick    string  `json:"nick"`                     //昵称
	Avatar  string  `json:"avatar"`                   //头像
	Extends Extends `json:"extends" gorm:"type:json"` //扩展

	ImUserid   int64  `json:"im_userid"`   //IM 用户ID
	ImUsername string `json:"im_username"` //IM 用户名
	ImPassword string `json:"im_password"` //IM 密码

	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt *timestamp.Timestamp
}

func (LiveUserEntity) TableName() string {
	return "live_users"
}
