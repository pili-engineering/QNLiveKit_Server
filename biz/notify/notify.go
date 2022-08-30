package notify

import (
	"context"
	"encoding/json"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/im"
	"github.com/qbox/livekit/utils/logger"
)

type ActionType string

const (
	ActionTypeCensorNotify ActionType = "censor_notify"
	ActionTypeCensorStop   ActionType = "censor_stop"
)

var actionTypeMap = map[ActionType]bool{
	ActionTypeCensorNotify: true,
	ActionTypeCensorStop:   true,
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
		return api.ErrInvalidArgument
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

// SendNotifyToLive 以主播的身份，给直播间发送通知消息
func SendNotifyToLive(ctx context.Context, anchor *model.LiveUserEntity, live *model.LiveEntity, action ActionType, data interface{}) error {
	log := logger.ReqLogger(ctx)
	if live == nil || live.ChatId == 0 {
		log.Errorf("no live group info ")
		return api.ErrInvalidArgument
	}

	if anchor == nil || anchor.ImUserid == 0 {
		log.Errorf("no live anchor im info")
		return api.ErrInvalidArgument
	}

	command := &LiveCommand{
		Action: action,
		Data:   data,
	}
	content, _ := json.Marshal(command)

	err := im.GetService().SendCommandMessageToGroup(ctx, anchor.ImUserid, live.ChatId, string(content))
	if err != nil {
		log.Errorf("SendCommandMessageToGroup error %s", err.Error())
	}
	return err
}
