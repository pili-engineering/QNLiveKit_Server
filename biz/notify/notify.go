package notify

import (
	"context"
	"encoding/json"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/fun/im"
	"github.com/qbox/livekit/utils/logger"
)

type ActionType string

const (
	ActionTypeCensorNotify ActionType = "censor_notify"
	ActionTypeCensorStop   ActionType = "censor_stop"
	ActionTypeLikeNotify   ActionType = "like_notify"
	ActionTypeGiftNotify   ActionType = "gift_notify"
)

var actionTypeMap = map[ActionType]bool{
	ActionTypeCensorNotify: true,
	ActionTypeCensorStop:   true,
	ActionTypeGiftNotify:   true,
}

func (t ActionType) IsValid() bool {
	_, ok := actionTypeMap[t]
	return ok
}

type LiveCommand struct {
	Action ActionType  `json:"action"`
	Data   interface{} `json:"data"`
}

// SendNotifyToUser 以系统管理员身份，给指定的用户发送通知消息
func SendNotifyToUser(ctx context.Context, user *model.LiveUserEntity, action ActionType, data interface{}) error {
	log := logger.ReqLogger(ctx)
	if user == nil || user.ImUserid == 0 {
		log.Errorf("no target im user info")
		return rest.ErrBadRequest
	}

	command := &LiveCommand{
		Action: action,
		Data:   data,
	}
	content, _ := json.Marshal(command)

	err := im.GetService().SendCommandMessageToUser(ctx, 0, user.ImUserid, string(content))
	if err != nil {
		log.Errorf("SendCommandMessageToUser error %s", err.Error())
	}
	return err
}

// SendNotifyToLive 以指定用户的身份，给直播间发送通知消息
func SendNotifyToLive(ctx context.Context, user *model.LiveUserEntity, live *model.LiveEntity, action ActionType, data interface{}) error {
	log := logger.ReqLogger(ctx)
	if live == nil || live.ChatId == 0 {
		log.Errorf("no live group info ")
		return rest.ErrBadRequest
	}

	if user == nil || user.ImUserid == 0 {
		log.Errorf("no user im info")
		return rest.ErrBadRequest
	}

	command := &LiveCommand{
		Action: action,
		Data:   data,
	}
	content, _ := json.Marshal(command)

	err := im.GetService().SendCommandMessageToGroup(ctx, user.ImUserid, live.ChatId, string(content))
	if err != nil {
		log.Errorf("SendCommandMessageToGroup error %s", err.Error())
	}
	return err
}
