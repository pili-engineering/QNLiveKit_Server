// @Author: wangsheng
// @Description:
// @File:  relay
// @Version: 1.0.0
// @Date: 2022/5/25 9:01 下午
// Copyright 2021 QINIU. All rights reserved

package model

import (
	"github.com/qbox/livekit/utils/timestamp"
)

const (
	RelaySessionStatusWaitAgree   = 0 //等待接收方同意
	RelaySessionStatusAgreed      = 1 //接收方已同意
	RelaySessionStatusInitSuccess = 2 //发起方已经完成跨房，等待对方完成
	RelaySessionStatusRecvSuccess = 3 //接收方已经完成跨房，等待对方完成
	RelaySessionStatusSuccess     = 4 //两方都完成跨房
	RelaySessionStatusRejected    = 5 //接收方拒绝
	RelaySessionStatusStopped     = 6 //结束
)

type RelaySession struct {
	ID         uint                 `gorm:"primary_key"`                      //DB 主Key，无业务含义
	SID        string               `json:"sid" gorm:"column:sid"`            //PK 会话ID
	InitUserId string               `json:"init_user_id"`                     //发起方主播ID
	InitRoomId string               `json:"init_room_id"`                     //发起方直播间ID
	RecvUserId string               `json:"recv_user_id"`                     //接收方主播ID
	RecvRoomId string               `json:"recv_room_id"`                     //接收方直播间ID
	Extends    Extends              `json:"extends" gorm:"type:varchar(512)"` //扩展数据
	Status     int                  `json:"status"`                           //PK 会话状态
	StartAt    *timestamp.Timestamp `json:"start_at"`                         //开始时间
	StopAt     *timestamp.Timestamp `json:"stop_at"`                          //结束时间
	CreatedAt  timestamp.Timestamp  `json:"created_at"`                       //创建时间
	UpdatedAt  timestamp.Timestamp  `json:"updated_at"`                       //更新时间
}

func (s *RelaySession) IsStopped() bool {
	return s.Status == RelaySessionStatusStopped
}
