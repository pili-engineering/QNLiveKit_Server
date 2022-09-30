// @Author: wangsheng
// @Description:
// @File:  user_dto
// @Version: 1.0.0
// @Date: 2022/5/23 11:17 上午
// Copyright 2021 QINIU. All rights reserved

package dto

import "github.com/qbox/livekit/biz/model"

type UserDto struct {
	UserId string `json:"user_id"` //应用内用户ID
	Nick   string `json:"nick"`    //昵称
	Avatar string `json:"avatar"`  //头像

	ImUserid   int64  `json:"im_userid"`
	ImUsername string `json:"im_username,omitempty"`

	Extends model.Extends `json:"extends"` //扩展属性
}

type UserProfileDto struct {
	UserId string `json:"user_id"` //应用内用户ID

	Nick    string        `json:"nick"`    //昵称
	Avatar  string        `json:"avatar"`  //头像
	Extends model.Extends `json:"extends"` //扩展属性

	ImUserid   int64  `json:"im_userid"`   //IM 用户ID
	ImUsername string `json:"im_username"` //IM 用户名
	ImPassword string `json:"im_password"` //IM 密码
}

func User2Dto(userEntity *model.LiveUserEntity) *UserDto {
	dto := UserDto{
		UserId: userEntity.UserId,

		Nick:    userEntity.Nick,
		Avatar:  userEntity.Avatar,
		Extends: userEntity.Extends,

		ImUserid: userEntity.ImUserid,
	}

	return &dto
}

func User2ProfileDto(userEntity *model.LiveUserEntity) *UserProfileDto {
	dto := UserProfileDto{
		UserId: userEntity.UserId,

		Nick:    userEntity.Nick,
		Avatar:  userEntity.Avatar,
		Extends: userEntity.Extends,

		ImUserid:   userEntity.ImUserid,
		ImUsername: userEntity.ImUsername,
		ImPassword: userEntity.ImPassword,
	}

	return &dto
}

func UserDto2Entity(dto *UserDto) *model.LiveUserEntity {
	entity := model.LiveUserEntity{
		UserId:  dto.UserId,
		Nick:    dto.Nick,
		Avatar:  dto.Avatar,
		Extends: dto.Extends,
	}

	return &entity
}
