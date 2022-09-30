package client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/mic"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/module/base/auth/internal/middleware"
	"github.com/qbox/livekit/module/base/live/internal/controller/client"
	"github.com/qbox/livekit/module/fun/rtc"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterMicRoutes(group *gin.RouterGroup) {
	liveGroup := group.Group("/mic")
	{
		liveGroup.POST("/", MicController.UpMic)
		liveGroup.DELETE("/", MicController.DownMic)
		liveGroup.GET("/room/list/:live_id", MicController.GetMicList)
		liveGroup.PUT("/extension", MicController.UpdateMicExtends)
		liveGroup.PUT("/switch", MicController.SwitchMic)
	}
}

type micController struct {
}

var MicController = &micController{}

type upMicResponse struct {
	api.Response
	Data struct {
		RtcToken string `json:"rtc_token"`
	} `json:"data"`
}

func (*micController) UpMic(context *gin.Context) {
	userInfo := context.MustGet(middleware.UserCtxKey).(*middleware.UserInfo)
	log := logger.ReqLogger(context)
	request := &mic.Request{}
	if err := context.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		context.JSON(400, gin.H{
			"code":    400,
			"message": "bind json error",
		})
		return
	}
	rtcToken, err := mic.GetService().UpMic(context, request, userInfo.UserId)
	if err != nil {
		log.Errorf("up mic failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "up mic failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, upMicResponse{
		Response: api.Response{
			Code:      http.StatusOK,
			Message:   "success",
			RequestId: log.ReqID(),
		},
		Data: struct {
			RtcToken string `json:"rtc_token"`
		}{
			RtcToken: rtcToken,
		},
	})
}

func (*micController) DownMic(context *gin.Context) {
	userInfo := context.MustGet(middleware.UserCtxKey).(*middleware.UserInfo)
	log := logger.ReqLogger(context)
	request := &mic.Request{}
	if err := context.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		context.JSON(400, gin.H{
			"code":    400,
			"message": "bind json error",
		})
		return
	}
	err := mic.GetService().DownMic(context, request, userInfo.UserId)
	if err != nil {
		log.Errorf("down mic failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "down mic failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

type MicListItem struct {
	User    client.UserInfo `json:"user"`
	Mic     bool            `json:"mic"`
	Camera  bool            `json:"camera"`
	Status  int             `json:"status"`
	Extends model.Extends   `json:"extends"`
}

type MicListResponse struct {
	api.Response
	Data []MicListItem `json:"data"`
}

func (*micController) GetMicList(context *gin.Context) {
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	micList, err := mic.GetService().LiveMicList(context, liveId)
	if err != nil {
		log.Errorf("get mic list failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get mic list failed",
			RequestId: log.ReqID(),
		})
		return
	}
	list := make([]MicListItem, 0, len(micList))
	rtcClient := rtc.GetService()
	for i := range micList {
		// 同步rtc麦位
		if !rtcClient.Online(micList[i].UserId, liveId) {
			continue
		}
		user, err := service.GetService().FindUser(context, micList[i].UserId)
		if err != nil {
			log.Errorf("get user failed, err: %v", err)
			continue
		}
		tmp := MicListItem{
			User: client.UserInfo{
				UserId:     user.UserId,
				ImUserId:   user.ImUserid,
				ImUsername: user.ImUsername,
				Nick:       user.Nick,
				Avatar:     user.Avatar,
				Extends:    user.Extends,
			},
			Status:  micList[i].Status,
			Extends: micList[i].Extends,
			Mic:     micList[i].Mic,
			Camera:  micList[i].Camera,
		}
		list = append(list, tmp)
	}
	context.JSON(http.StatusOK, MicListResponse{
		Response: api.Response{
			Code:      http.StatusOK,
			Message:   "success",
			RequestId: log.ReqID(),
		},
		Data: list,
	})
}

func (*micController) UpdateMicExtends(context *gin.Context) {
	log := logger.ReqLogger(context)
	request := &mic.Request{}
	if err := context.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		context.JSON(400, gin.H{
			"code":    400,
			"message": "bind json error",
		})
		return
	}
	err := mic.GetService().UpdateMicExtends(context, request.LiveId, request.UserId, request.Extends)
	if err != nil {
		log.Errorf("update mic extends failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "update mic extends failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

type switchMicRequest struct {
	LiveId string `json:"live_id"`
	UserId string `json:"user_id"`
	Type   string `json:"type"`
	Flag   bool   `json:"status"`
}

func (*micController) SwitchMic(context *gin.Context) {
	log := logger.ReqLogger(context)
	request := &switchMicRequest{}
	if err := context.BindJSON(request); err != nil {
		log.Error("bind json error", err)
		context.JSON(400, gin.H{
			"code":    400,
			"message": "bind json error",
		})
		return
	}
	err := mic.GetService().SwitchUserMic(context, request.LiveId, request.UserId, request.Type, request.Flag)
	if err != nil {
		log.Errorf("switch mic failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "switch mic failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      http.StatusOK,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}
