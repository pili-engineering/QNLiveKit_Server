package model

import (
	"github.com/qbox/livekit/utils/timestamp"
)

type LiveEntity struct {
	Id              uint                 `gorm:"primary_key" json:"id"`
	LiveId          string               `json:"live_id"`
	Title           string               `json:"title"`
	Notice          string               `json:"notice"`
	CoverUrl        string               `json:"cover_url"`
	Extends         Extends              `json:"extends" gorm:"type:json"`
	AnchorId        string               `json:"anchor_id"`
	Status          int                  `json:"live_status"`
	PkId            string               `json:"pk_id"`
	OnlineCount     int                  `json:"online_count"`
	StartAt         timestamp.Timestamp  `json:"start_at"`
	EndAt           timestamp.Timestamp  `json:"end_at"`
	RelaySessionId  string               `json:"relay_session_id"`
	ChatId          int64                `json:"chat_id"`
	PushUrl         string               `json:"push_url"`
	RtmpPlayUrl     string               `json:"rtmp_play_url"`
	FlvPlayUrl      string               `json:"flv_play_url"`
	HlsPlayUrl      string               `json:"hls_play_url"`
	LastHeartbeatAt timestamp.Timestamp  `json:"last_heartbeat_at"`
	CreatedAt       timestamp.Timestamp  `json:"created_at"`
	UpdatedAt       timestamp.Timestamp  `json:"updated_at"`
	DeletedAt       *timestamp.Timestamp `json:"deleted_at"`
}

type LiveRoomUserStatus int

const (
	LiveRoomUserStatusLeave  LiveRoomUserStatus = 0 //离开直播间
	LiveRoomUserStatusOnline LiveRoomUserStatus = 1 //在直播间，心跳有效
)

type LiveRoomUserEntity struct {
	Id          uint                 `gorm:"primary_key" json:"id"`
	LiveId      string               `json:"live_id"`
	UserId      string               `json:"user_id"` // userId 应该为唯一索引
	Status      LiveRoomUserStatus   `json:"status"`
	HeartBeatAt *timestamp.Timestamp `json:"heart_beat_at"`
	CreatedAt   timestamp.Timestamp  `json:"created_at"`
	UpdatedAt   timestamp.Timestamp  `json:"updated_at"`
	DeletedAt   *timestamp.Timestamp `json:"deleted_at"`
}

type LiveMicEntity struct {
	Id        uint                 `gorm:"primary_key" json:"id"`
	LiveId    string               `json:"room_id"`
	UserId    string               `json:"user_id"`
	Mic       bool                 `json:"mic"`
	Camera    bool                 `json:"camera"`
	Status    int                  `json:"status"`
	Extends   Extends              `json:"extends" gorm:"type:json"`
	CreatedAt timestamp.Timestamp  `json:"created_at"`
	UpdatedAt timestamp.Timestamp  `json:"updated_at"`
	DeletedAt *timestamp.Timestamp `json:"deleted_at"`
}

const (
	LiveStatusPrepare = iota
	LiveStatusOn
	LiveStatusOff
)

const (
	LiveRoomUserMicStatusJoin = iota
	LiveRoomUserMicStatusLeave
	LiveRoomUserMicForbidden
)
