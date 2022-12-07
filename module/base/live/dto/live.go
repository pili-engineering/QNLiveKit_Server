// @Author: wangsheng
// @Description:
// @File:  live
// @Version: 1.0.0
// @Date: 2022/6/28 11:52 上午
// Copyright 2021 QINIU. All rights reserved

package dto

import (
	"github.com/qbox/livekit/biz/model"
	userDto "github.com/qbox/livekit/module/base/user/dto"
)

type LiveInfoDto struct {
	LiveId       string                   `json:"live_id"`
	Title        string                   `json:"title"`
	Notice       string                   `json:"notice"`
	CoverUrl     string                   `json:"cover_url"`
	Extends      model.Extends            `json:"extends"`
	AnchorInfo   userDto.UserDto          `json:"anchor_info"`
	AnchorStatus model.LiveRoomUserStatus `json:"anchor_status"`
	RoomToken    string                   `json:"room_token"`
	PkId         string                   `json:"pk_id"`
	OnlineCount  int                      `json:"online_count"`
	StartTime    int64                    `json:"start_time"`
	EndTime      int64                    `json:"end_time"`
	ChatId       int64                    `json:"chat_id"`
	PushUrl      string                   `json:"push_url"`
	HlsUrl       string                   `json:"hls_url"`
	RtmpUrl      string                   `json:"rtmp_url"`
	FlvUrl       string                   `json:"flv_url"`
	Pv           int                      `json:"pv"`
	Uv           int                      `json:"uv"`
	TotalCount   int                      `json:"total_count"`
	TotalMics    int                      `json:"total_mics"`
	LiveStatus   int                      `json:"live_status"`

	StopReason string `json:"stop_reason,omitempty"`  //关闭原因：censor 内容违规
	StopUserId string `json:"stop_user_id,omitempty"` //关闭直播的管理员用户ID
	StopTime   int64  `json:"stop_at,omitempty"`      //关闭时间
}

func BuildLiveDto(liveEntity *model.LiveEntity, userEntity *model.LiveUserEntity) LiveInfoDto {
	ret := LiveInfoDto{}

	ret.LiveId = liveEntity.LiveId
	ret.Title = liveEntity.Title
	ret.Notice = liveEntity.Notice
	ret.CoverUrl = liveEntity.CoverUrl
	ret.Extends = liveEntity.Extends
	ret.RoomToken = ""
	ret.PkId = liveEntity.PkId
	ret.OnlineCount = liveEntity.OnlineCount
	if liveEntity.StartAt != nil {
		ret.StartTime = liveEntity.StartAt.Unix()
	}
	ret.EndTime = liveEntity.EndAt.Unix()
	ret.ChatId = liveEntity.ChatId
	ret.PushUrl = liveEntity.PushUrl
	ret.HlsUrl = liveEntity.HlsPlayUrl
	ret.RtmpUrl = liveEntity.RtmpPlayUrl
	ret.FlvUrl = liveEntity.FlvPlayUrl
	ret.Pv = 0
	ret.Uv = 0
	ret.TotalCount = 0
	ret.TotalMics = 0
	ret.LiveStatus = liveEntity.Status
	ret.StopReason = liveEntity.StopReason
	ret.StopUserId = liveEntity.StopUserId
	if liveEntity.StopAt != nil {
		ret.StopTime = liveEntity.StopAt.Unix()
	}

	ret.AnchorInfo.UserId = userEntity.UserId
	ret.AnchorInfo.ImUserid = userEntity.ImUserid
	ret.AnchorInfo.Nick = userEntity.Nick
	ret.AnchorInfo.Avatar = userEntity.Avatar
	ret.AnchorInfo.Extends = userEntity.Extends
	return ret
}
