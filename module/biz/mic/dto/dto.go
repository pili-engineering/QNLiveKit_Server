package dto

import (
	"github.com/qbox/livekit/biz/model"
	userDto "github.com/qbox/livekit/module/base/user/dto"
)

type MicItemDto struct {
	User    userDto.UserDto `json:"user"`
	Mic     bool            `json:"mic"`
	Camera  bool            `json:"camera"`
	Status  int             `json:"status"`
	Extends model.Extends   `json:"extends"`
}
