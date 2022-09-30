package service

import (
	"context"

	"github.com/qbox/livekit/biz/model"
)

type Request struct {
	LiveId  string        `json:"live_id"`
	UserId  string        `json:"user_id"`
	Mic     bool          `json:"mic"`
	Camera  bool          `json:"camera"`
	Extends model.Extends `json:"extends"`
}

type IService interface {
	KickUser(context context.Context, userId, liveId string) (err error)

	UpMic(context context.Context, req *Request, userId string) (rtcToken string, err error)

	DownMic(context context.Context, req *Request, userId string) (err error)

	DownMicManual(context context.Context, liveId, userId string) (err error)

	ForbidMic(context context.Context, liveId, userId string) (err error)

	UnForbidMic(context context.Context, liveId, userId string) (err error)

	UserStatus(context context.Context, liveId, userId string) (status int, err error)

	LiveMicList(context context.Context, liveId string) (mics []model.LiveMicEntity, err error)

	UpdateMicExtends(context context.Context, liveId, userId string, extends model.Extends) (err error)

	SwitchUserMic(context context.Context, liveId, userId, tp string, flag bool) (err error)
}
