// @Author: wangsheng
// @Description:
// @File:  live
// @Version: 1.0.0
// @Date: 2022/6/28 11:52 上午
// Copyright 2021 QINIU. All rights reserved

package dto

import "github.com/qbox/livekit/biz/model"

type LiveInfoDto struct {
	LiveId       string        `json:"live_id"`
	Title        string        `json:"title"`
	Notice       string        `json:"notice"`
	CoverUrl     string        `json:"cover_url"`
	Extends      model.Extends `json:"extends"`
	AnchorInfo   UserDto       `json:"anchor_info"`
	AnchorStatus model.LiveRoomUserStatus
	RoomToken    string `json:"room_token"`
	PkId         string `json:"pk_id"`
	OnlineCount  int    `json:"online_count"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	ChatId       int64  `json:"chat_id"`
	PushUrl      string `json:"push_url"`
	HlsUrl       string `json:"hls_url"`
	RtmpUrl      string `json:"rtmp_url"`
	FlvUrl       string `json:"flv_url"`
	Pv           int    `json:"pv"`
	Uv           int    `json:"uv"`
	TotalCount   int    `json:"total_count"`
	TotalMics    int    `json:"total_mics"`
	LiveStatus   int    `json:"live_status"`
}
