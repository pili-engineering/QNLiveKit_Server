package client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/user"
	userDto "github.com/qbox/livekit/module/base/user/dto"
	"github.com/qbox/livekit/module/biz/mic/dto"
	"github.com/qbox/livekit/module/biz/mic/internal/controller/impl"
	"github.com/qbox/livekit/module/biz/mic/service"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	//liveGroup := group.Group("/mic")
	//{
	//	liveGroup.POST("/", MicController.UpMic)
	//	liveGroup.DELETE("/", MicController.DownMic)
	//	liveGroup.GET("/room/list/:live_id", MicController.GetMicList)
	//	liveGroup.PUT("/extension", MicController.UpdateMicExtends)
	//	liveGroup.PUT("/switch", MicController.SwitchMic)
	//}

	httpq.ClientHandle(http.MethodPost, "/mic/", MicController.UpMic)
	httpq.ClientHandle(http.MethodDelete, "/mic/", MicController.DownMic)
	httpq.ClientHandle(http.MethodGet, "/mic/room/list/:live_id", MicController.GetMicList)
	httpq.ClientHandle(http.MethodPut, "/mic/extension", MicController.UpdateMicExtends)
	httpq.ClientHandle(http.MethodPut, "/mic/switch", MicController.SwitchMic)
}

type micController struct {
}

var MicController = &micController{}

type UpMicResult struct {
	RtcToken string `json:"rtc_token"`
}

// UpMic 用户连麦请求
// return *UpMicResult
func (*micController) UpMic(ctx *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(ctx)
	log := logger.ReqLogger(ctx)
	request := &service.Request{}
	if err := ctx.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	rtcToken, err := impl.GetInstance().UpMic(ctx, request, userInfo.UserId)
	if err != nil {
		log.Errorf("up mic failed, err: %v", err)
		return nil, err
	}

	return &UpMicResult{
		RtcToken: rtcToken,
	}, nil
}

// DownMic 上麦用户请求下麦
// return nil
func (*micController) DownMic(ctx *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(ctx)
	log := logger.ReqLogger(ctx)
	request := &service.Request{}
	if err := ctx.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	err := impl.GetInstance().DownMic(ctx, request, userInfo.UserId)
	if err != nil {
		log.Errorf("down mic failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}

// GetMicList 获取直播间麦位列表
// return []*dto.MicItemDto
func (*micController) GetMicList(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("live_id")
	micList, err := impl.GetInstance().LiveMicList(ctx, liveId)
	if err != nil {
		log.Errorf("get mic list failed, err: %v", err)
		return nil, err
	}
	list := make([]*dto.MicItemDto, 0, len(micList))
	rtcClient := rtc.GetService()
	for i := range micList {
		// 同步rtc麦位
		if !rtcClient.Online(micList[i].UserId, liveId) {
			continue
		}
		userEntity, err := user.GetService().FindUser(ctx, micList[i].UserId)
		if err != nil {
			log.Errorf("get user failed, err: %v", err)
			continue
		}
		tmp := &dto.MicItemDto{
			User: userDto.UserDto{
				UserId:     userEntity.UserId,
				ImUserid:   userEntity.ImUserid,
				ImUsername: userEntity.ImUsername,
				Nick:       userEntity.Nick,
				Avatar:     userEntity.Avatar,
				Extends:    userEntity.Extends,
			},
			Status:  micList[i].Status,
			Extends: micList[i].Extends,
			Mic:     micList[i].Mic,
			Camera:  micList[i].Camera,
		}
		list = append(list, tmp)
	}
	return list, nil
}

// UpdateMicExtends 更新麦位的扩展信息
// return nil
func (*micController) UpdateMicExtends(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	request := &service.Request{}
	if err := context.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	err := impl.GetInstance().UpdateMicExtends(context, request.LiveId, request.UserId, request.Extends)
	if err != nil {
		log.Errorf("update mic extends failed, err: %v", err)
		return nil, err
	}

	return nil, nil
}

type switchMicRequest struct {
	LiveId string `json:"live_id"`
	UserId string `json:"user_id"`
	Type   string `json:"type"`
	Flag   bool   `json:"status"`
}

// SwitchMic 切换连麦状态
// return nil
func (*micController) SwitchMic(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	request := &switchMicRequest{}
	if err := context.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	err := impl.GetInstance().SwitchUserMic(context, request.LiveId, request.UserId, request.Type, request.Flag)
	if err != nil {
		log.Errorf("switch mic failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}
