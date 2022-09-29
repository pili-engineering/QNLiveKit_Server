// @Author: wangsheng
// @Description:
// @File:  live_controller
// @Version: 1.0.0
// @Date: 2022/5/19 2:48 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/biz/report"
	user2 "github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterLiveRoutes(group *gin.RouterGroup) {
	liveGroup := group.Group("/live")
	{
		liveGroup.POST("/room/instance", LiveController.CreateLive)
		liveGroup.DELETE("/room/instance/:live_id", LiveController.DeleteLive)
		liveGroup.GET("/room/info/:live_id", LiveController.LiveRoomInfo)
		liveGroup.DELETE("/room/:live_id", LiveController.StopLive)
		liveGroup.PUT("/room/:live_id", LiveController.StartLive)
		liveGroup.GET("/room", LiveController.SearchLive)
		liveGroup.POST("/room/user/:live_id", LiveController.JoinLive)
		liveGroup.GET("/room/list", LiveController.LiveList)
		liveGroup.GET("/room/list/anchor", LiveController.LiveListAnchor)
		liveGroup.DELETE("/room/user/:live_id", LiveController.LeaveLive)
		liveGroup.GET("/room/heartbeat/:live_id", LiveController.Heartbeat)
		liveGroup.PUT("/room/extends", LiveController.UpdateExtends)
		liveGroup.GET("/room/user_list", LiveController.LiveUserList)
		liveGroup.PUT("/room/:live_id/like", LiveController.PutLike)
	}
}

type liveController struct {
}

var LiveController = &liveController{}

type LiveResponse struct {
	api.Response
	Data dto.LiveInfoDto `json:"data"`
}

type LiveListResponse struct {
	api.Response
	Data struct {
		TotalCount int               `json:"total_count"`
		PageTotal  int               `json:"page_total"`
		EndPage    bool              `json:"end_page"`
		List       []dto.LiveInfoDto `json:"list"`
	} `json:"data"`
}

func (*liveController) CreateLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	request := &live.CreateLiveRequest{}
	if context.Bind(request) != nil {
		context.JSON(400, api.Response{
			Code:      400,
			Message:   "invalid request",
			RequestId: log.ReqID(),
		})
		return
	}
	request.AnchorId = userInfo.UserId
	liveEntity, err := live.GetService().CreateLive(context, request)
	if err != nil {
		log.Errorf("create liveEntity failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "create liveEntity failed",
			RequestId: log.ReqID(),
		})
		return
	}
	user, err := user2.GetService().FindUser(context, userInfo.UserId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "find user failed",
			RequestId: log.ReqID(),
		})
		return
	}
	response := &LiveResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.LiveId = liveEntity.LiveId
	response.Data.Title = liveEntity.Title
	response.Data.Notice = liveEntity.Notice
	response.Data.CoverUrl = liveEntity.CoverUrl
	response.Data.Extends = liveEntity.Extends
	response.Data.AnchorInfo.UserId = user.UserId
	response.Data.AnchorInfo.ImUserid = user.ImUserid
	response.Data.AnchorInfo.Nick = user.Nick
	response.Data.AnchorInfo.Avatar = user.Avatar
	response.Data.AnchorInfo.Extends = user.Extends
	response.Data.RoomToken = ""
	response.Data.PkId = liveEntity.PkId
	response.Data.OnlineCount = liveEntity.OnlineCount
	if liveEntity.StartAt != nil {
		response.Data.StartTime = liveEntity.StartAt.Unix()
	}
	response.Data.EndTime = liveEntity.EndAt.Unix()
	response.Data.ChatId = liveEntity.ChatId
	response.Data.PushUrl = liveEntity.PushUrl
	response.Data.HlsUrl = liveEntity.HlsPlayUrl
	response.Data.RtmpUrl = liveEntity.RtmpPlayUrl
	response.Data.FlvUrl = liveEntity.FlvPlayUrl
	response.Data.Pv = 0
	response.Data.Uv = 0
	response.Data.TotalCount = 0
	response.Data.TotalMics = 0
	response.Data.LiveStatus = liveEntity.Status
	context.JSON(http.StatusOK, response)
}

func (*liveController) DeleteLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	err := live.GetService().DeleteLive(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("delete live failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "delete live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (c *liveController) LiveRoomInfo(context *gin.Context) {
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	liveInfo, err := live.GetService().LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get liveInfo info failed",
			RequestId: log.ReqID(),
		})
		return
	}
	user, err := user2.GetService().FindUser(context, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "find user failed",
			RequestId: log.ReqID(),
		})
		return
	}
	_, onlineCount, err := live.GetService().LiveUserList(context, liveId, 1, 10)
	if err != nil {
		log.Errorf("get live user list failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get live user list failed",
			RequestId: log.ReqID(),
		})
		return
	}

	anchorStatus, err := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("get live anchor status error: %v", err)
		context.JSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	response := &LiveResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.LiveId = liveInfo.LiveId
	response.Data.Title = liveInfo.Title
	response.Data.Notice = liveInfo.Notice
	response.Data.CoverUrl = liveInfo.CoverUrl
	response.Data.Extends = liveInfo.Extends
	response.Data.AnchorInfo.UserId = user.UserId
	response.Data.AnchorInfo.ImUserid = user.ImUserid
	response.Data.AnchorInfo.Nick = user.Nick
	response.Data.AnchorInfo.Avatar = user.Avatar
	response.Data.AnchorInfo.Extends = user.Extends
	response.Data.AnchorStatus = anchorStatus
	response.Data.RoomToken = ""
	response.Data.PkId = liveInfo.PkId
	response.Data.OnlineCount = onlineCount
	response.Data.EndTime = liveInfo.EndAt.Unix()
	response.Data.ChatId = liveInfo.ChatId
	response.Data.PushUrl = liveInfo.PushUrl
	response.Data.HlsUrl = liveInfo.HlsPlayUrl
	response.Data.RtmpUrl = liveInfo.RtmpPlayUrl
	response.Data.FlvUrl = liveInfo.FlvPlayUrl
	response.Data.Pv = 0
	response.Data.Uv = 0
	response.Data.TotalCount = 0
	response.Data.TotalMics = 0
	response.Data.LiveStatus = liveInfo.Status
	response.Data.StopReason = liveInfo.StopReason
	response.Data.StopUserId = liveInfo.StopUserId
	if liveInfo.StartAt != nil {
		response.Data.StartTime = liveInfo.StartAt.Unix()
	}
	if liveInfo.StopAt != nil {
		response.Data.StopTime = liveInfo.StopAt.Unix()
	}
	context.JSON(http.StatusOK, response)
}

func (*liveController) getLiveAnchorStatus(ctx context.Context, liveId string, anchorId string) (model.LiveRoomUserStatus, error) {
	log := logger.ReqLogger(ctx)
	liveUser, err := live.GetService().FindLiveRoomUser(ctx, liveId, anchorId)
	if err != nil {
		if !api.IsNotFoundError(err) {
			log.Errorf("find anchor user error: %v", err)
			return model.LiveRoomUserStatusLeave, err
		} else {
			return model.LiveRoomUserStatusLeave, nil
		}
	} else {
		return liveUser.Status, nil
	}
}

func (*liveController) StopLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	err := live.GetService().StopLive(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("stop live failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "stop live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (*liveController) StartLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	roomToken, err := live.GetService().StartLive(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("start live failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "start live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	liveInfo, err := live.GetService().LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get liveInfo info failed",
			RequestId: log.ReqID(),
		})
		return
	}
	user, err := user2.GetService().FindUser(context, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "find user failed",
			RequestId: log.ReqID(),
		})
		return
	}
	response := &LiveResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.LiveId = liveInfo.LiveId
	response.Data.Title = liveInfo.Title
	response.Data.Notice = liveInfo.Notice
	response.Data.CoverUrl = liveInfo.CoverUrl
	response.Data.Extends = liveInfo.Extends
	response.Data.AnchorInfo.UserId = user.UserId
	response.Data.AnchorInfo.ImUserid = user.ImUserid
	response.Data.AnchorInfo.Nick = user.Nick
	response.Data.AnchorInfo.Avatar = user.Avatar
	response.Data.AnchorInfo.Extends = user.Extends
	response.Data.AnchorStatus = model.LiveRoomUserStatusOnline
	response.Data.RoomToken = roomToken
	response.Data.PkId = liveInfo.PkId
	response.Data.OnlineCount = liveInfo.OnlineCount
	if liveInfo.StartAt != nil {
		response.Data.StartTime = liveInfo.StartAt.Unix()
	}
	response.Data.EndTime = liveInfo.EndAt.Unix()
	response.Data.ChatId = liveInfo.ChatId
	response.Data.PushUrl = liveInfo.PushUrl
	response.Data.HlsUrl = liveInfo.HlsPlayUrl
	response.Data.RtmpUrl = liveInfo.RtmpPlayUrl
	response.Data.FlvUrl = liveInfo.FlvPlayUrl
	response.Data.Pv = 0
	response.Data.Uv = 0
	response.Data.TotalCount = 0
	response.Data.TotalMics = 0
	response.Data.LiveStatus = liveInfo.Status
	context.JSON(http.StatusOK, response)
}

func (c *liveController) SearchLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	keyword := context.Query("keyword")
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	liveList1, totalCount1, err := live.GetService().SearchLive(context, keyword, 1, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search live failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "search live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	liveList2, totalCount2, err := live.GetService().SearchLive(context, keyword, 2, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search live failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "search live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	liveList3, totalCount3, err := live.GetService().SearchLive(context, keyword, 3, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search live failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "search live failed",
			RequestId: log.ReqID(),
		})
		return
	}
	liveList := append(liveList1)
	liveList = append(liveList2)
	liveList = append(liveList3)
	endPage := false
	if len(liveList) < pageSizeInt {
		endPage = true
	}
	response := &LiveListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = totalCount1 + totalCount2 + totalCount3
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	list := make([]dto.LiveInfoDto, len(liveList), len(liveList))
	for i := range liveList {
		liveInfo, err := live.GetService().LiveInfo(context, liveList[i].LiveId)
		if err != nil {
			log.Errorf("get liveInfo info failed, err: %v", err)
			continue
		}
		user, err := user2.GetService().FindUser(context, userInfo.UserId)
		if err != nil {
			log.Errorf("find user failed, err: %v", err)
			continue
		}

		anchorStatus, _ := c.getLiveAnchorStatus(context, liveList[i].LiveId, liveList[i].AnchorId)

		list[i].LiveId = liveInfo.LiveId
		list[i].Title = liveInfo.Title
		list[i].Notice = liveInfo.Notice
		list[i].CoverUrl = liveInfo.CoverUrl
		list[i].Extends = liveInfo.Extends
		list[i].AnchorInfo.UserId = user.UserId
		list[i].AnchorInfo.ImUserid = user.ImUserid
		list[i].AnchorInfo.Nick = user.Nick
		list[i].AnchorInfo.Avatar = user.Avatar
		list[i].AnchorInfo.Extends = user.Extends
		list[i].AnchorStatus = anchorStatus
		list[i].RoomToken = ""
		list[i].PkId = liveInfo.PkId
		list[i].OnlineCount = liveInfo.OnlineCount
		list[i].EndTime = liveInfo.EndAt.Unix()
		list[i].ChatId = liveInfo.ChatId
		list[i].PushUrl = liveInfo.PushUrl
		list[i].HlsUrl = liveInfo.HlsPlayUrl
		list[i].RtmpUrl = liveInfo.RtmpPlayUrl
		list[i].FlvUrl = liveInfo.FlvPlayUrl
		list[i].Pv = 0
		list[i].Uv = 0
		list[i].TotalCount = 0
		list[i].TotalMics = 0
		list[i].LiveStatus = liveInfo.Status
		list[i].StopReason = liveInfo.StopReason
		list[i].StopUserId = liveInfo.StopUserId
		if liveInfo.StartAt != nil {
			list[i].StartTime = liveInfo.StartAt.Unix()
		}
		if liveInfo.StopAt != nil {
			list[i].StopTime = liveInfo.StopAt.Unix()
		}
	}
	response.Data.List = list
	context.JSON(http.StatusOK, response)
}

func (c *liveController) JoinLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	err := live.GetService().JoinLiveRoom(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("get live info failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get live info failed",
			RequestId: log.ReqID(),
		})
		return
	}
	liveInfo, err := live.GetService().LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get liveInfo info failed",
			RequestId: log.ReqID(),
		})
		return
	}
	user, err := user2.GetService().FindUser(context, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "find user failed",
			RequestId: log.ReqID(),
		})
		return
	}
	rService := report.GetService()
	statsSingleLiveEntity := &model.StatsSingleLiveEntity{
		LiveId: liveId,
		UserId: userInfo.UserId,
		Type:   model.StatsTypeLive,
		Count:  1,
	}
	rService.UpdateSingleLive(context, statsSingleLiveEntity)
	anchorStatus, _ := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)
	response := &LiveResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.LiveId = liveInfo.LiveId
	response.Data.Title = liveInfo.Title
	response.Data.Notice = liveInfo.Notice
	response.Data.CoverUrl = liveInfo.CoverUrl
	response.Data.Extends = liveInfo.Extends
	response.Data.AnchorInfo.UserId = user.UserId
	response.Data.AnchorInfo.ImUserid = user.ImUserid
	response.Data.AnchorInfo.Nick = user.Nick
	response.Data.AnchorInfo.Avatar = user.Avatar
	response.Data.AnchorInfo.Extends = user.Extends
	response.Data.AnchorStatus = anchorStatus
	response.Data.RoomToken = ""
	response.Data.PkId = liveInfo.PkId
	response.Data.OnlineCount = liveInfo.OnlineCount
	response.Data.EndTime = liveInfo.EndAt.Unix()
	response.Data.ChatId = liveInfo.ChatId
	response.Data.PushUrl = liveInfo.PushUrl
	response.Data.HlsUrl = liveInfo.HlsPlayUrl
	response.Data.RtmpUrl = liveInfo.RtmpPlayUrl
	response.Data.FlvUrl = liveInfo.FlvPlayUrl
	response.Data.Pv = 0
	response.Data.Uv = 0
	response.Data.TotalCount = 0
	response.Data.TotalMics = 0
	response.Data.LiveStatus = liveInfo.Status
	response.Data.StopReason = liveInfo.StopReason
	response.Data.StopUserId = liveInfo.StopUserId
	if liveInfo.StartAt != nil {
		response.Data.StartTime = liveInfo.StartAt.Unix()
	}
	if liveInfo.StopAt != nil {
		response.Data.StopTime = liveInfo.StopAt.Unix()
	}
	context.JSON(http.StatusOK, response)
}

func (c *liveController) LiveList(context *gin.Context) {
	log := logger.ReqLogger(context)
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page num is not int, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page size is not int, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	if pageNumInt <= 0 || pageSizeInt <= 0 {
		log.Errorf("page num or page size is not right, page num: %v, page size: %v", pageNumInt, pageSizeInt)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page num or page size is not right",
			RequestId: log.ReqID(),
		})
		return
	}
	liveList, totalCount, err := live.GetService().LiveList(context, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("get live list failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get live list failed",
			RequestId: log.ReqID(),
		})
		return
	}
	endPage := false
	if len(liveList) < pageSizeInt {
		endPage = true
	}
	response := &LiveListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = totalCount
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	list := make([]dto.LiveInfoDto, len(liveList), len(liveList))
	for i := range liveList {
		liveInfo, err := live.GetService().LiveInfo(context, liveList[i].LiveId)
		if err != nil {
			log.Errorf("get liveInfo info failed, err: %v", err)
			continue
		}
		log.Infof("liveInfo: %v", liveInfo)
		user, err := user2.GetService().FindUser(context, liveInfo.AnchorId)
		log.Infof("user: %v", user)
		if err != nil {
			log.Errorf("find user failed, err: %v", err)
			continue
		}
		anchorStatus, _ := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)

		list[i].LiveId = liveInfo.LiveId
		list[i].Title = liveInfo.Title
		list[i].Notice = liveInfo.Notice
		list[i].CoverUrl = liveInfo.CoverUrl
		list[i].Extends = liveInfo.Extends
		list[i].AnchorInfo.UserId = user.UserId
		list[i].AnchorInfo.ImUserid = user.ImUserid
		list[i].AnchorInfo.Nick = user.Nick
		list[i].AnchorInfo.Avatar = user.Avatar
		list[i].AnchorInfo.Extends = user.Extends
		list[i].AnchorStatus = anchorStatus
		list[i].RoomToken = ""
		list[i].PkId = liveInfo.PkId
		list[i].OnlineCount = liveInfo.OnlineCount
		list[i].EndTime = liveInfo.EndAt.Unix()
		list[i].ChatId = liveInfo.ChatId
		list[i].PushUrl = liveInfo.PushUrl
		list[i].HlsUrl = liveInfo.HlsPlayUrl
		list[i].RtmpUrl = liveInfo.RtmpPlayUrl
		list[i].FlvUrl = liveInfo.FlvPlayUrl
		list[i].Pv = 0
		list[i].Uv = 0
		list[i].TotalCount = 0
		list[i].TotalMics = 0
		list[i].LiveStatus = liveInfo.Status
		list[i].StopReason = liveInfo.StopReason
		list[i].StopUserId = liveInfo.StopUserId

		if liveInfo.StopAt != nil {
			list[i].StopTime = liveInfo.StopAt.Unix()
		}
		log.Infof("liveInfo: %v", liveInfo)
	}
	response.Data.List = list
	context.JSON(http.StatusOK, response)
}

func (c *liveController) LiveListAnchor(context *gin.Context) {
	log := logger.ReqLogger(context)
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	uInfo := liveauth.GetUserInfo(context)
	if uInfo == nil {
		log.Errorf("user info not exist")
		context.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), api.ErrNotFound))
		return
	}
	anchorId := uInfo.UserId
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page num is not int, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page size is not int, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	if pageNumInt <= 0 || pageSizeInt <= 0 {
		log.Errorf("page num or page size is not right, page num: %v, page size: %v", pageNumInt, pageSizeInt)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page num or page size is not right",
			RequestId: log.ReqID(),
		})
		return
	}
	liveList, totalCount, err := live.GetService().LiveListAnchor(context, pageNumInt, pageSizeInt, anchorId)
	if err != nil {
		log.Errorf("get live list anchor failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get live list anchor failed",
			RequestId: log.ReqID(),
		})
		return
	}
	endPage := false
	if len(liveList) < pageSizeInt {
		endPage = true
	}
	response := &LiveListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = totalCount
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	list := make([]dto.LiveInfoDto, len(liveList), len(liveList))
	for i := range liveList {
		liveInfo, err := live.GetService().LiveInfo(context, liveList[i].LiveId)
		if err != nil {
			log.Errorf("get liveInfo info failed, err: %v", err)
			continue
		}
		log.Infof("liveInfo: %v", liveInfo)
		user, err := user2.GetService().FindUser(context, liveInfo.AnchorId)
		log.Infof("user: %v", user)
		if err != nil {
			log.Errorf("find user failed, err: %v", err)
			continue
		}
		anchorStatus, _ := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)

		list[i].LiveId = liveInfo.LiveId
		list[i].Title = liveInfo.Title
		list[i].Notice = liveInfo.Notice
		list[i].CoverUrl = liveInfo.CoverUrl
		list[i].Extends = liveInfo.Extends
		list[i].AnchorInfo.UserId = user.UserId
		list[i].AnchorInfo.ImUserid = user.ImUserid
		list[i].AnchorInfo.Nick = user.Nick
		list[i].AnchorInfo.Avatar = user.Avatar
		list[i].AnchorInfo.Extends = user.Extends
		list[i].AnchorStatus = anchorStatus
		list[i].RoomToken = ""
		list[i].PkId = liveInfo.PkId
		list[i].OnlineCount = liveInfo.OnlineCount
		list[i].EndTime = liveInfo.EndAt.Unix()
		list[i].ChatId = liveInfo.ChatId
		list[i].PushUrl = liveInfo.PushUrl
		list[i].HlsUrl = liveInfo.HlsPlayUrl
		list[i].RtmpUrl = liveInfo.RtmpPlayUrl
		list[i].FlvUrl = liveInfo.FlvPlayUrl
		list[i].Pv = 0
		list[i].Uv = 0
		list[i].TotalCount = 0
		list[i].TotalMics = 0
		list[i].LiveStatus = liveInfo.Status
		list[i].StopReason = liveInfo.StopReason
		list[i].StopUserId = liveInfo.StopUserId
		if liveInfo.StartAt != nil {
			list[i].StartTime = liveInfo.StartAt.Unix()
		}
		if liveInfo.StopAt != nil {
			list[i].StopTime = liveInfo.StopAt.Unix()
		}
		log.Infof("liveInfo: %v", liveInfo)
	}
	response.Data.List = list
	context.JSON(http.StatusOK, response)
}

func (*liveController) LeaveLive(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	if liveId == "" {
		log.Errorf("live_id is empty")
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "live_id is empty",
			RequestId: log.ReqID(),
		})
		return
	}
	err := live.GetService().LeaveLiveRoom(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("leave live room failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "leave live room failed",
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

type HeartBeatResponse struct {
	api.Response
	Data struct {
		LiveId string `json:"live_id"`
		Status int    `json:"live_status"`
	} `json:"data"`
}

func (*liveController) Heartbeat(context *gin.Context) {
	userInfo := context.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	if liveId == "" {
		log.Errorf("live_id is empty")
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "live_id is empty",
			RequestId: log.ReqID(),
		})
		return
	}
	live, err := live.GetService().Heartbeat(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("heartbeat failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "heartbeat failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, &HeartBeatResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data: struct {
			LiveId string `json:"live_id"`
			Status int    `json:"live_status"`
		}{
			LiveId: live.LiveId,
			Status: live.Status,
		},
	})
}

type UpdateExtendsRequest struct {
	LiveId  string        `json:"live_id"`
	Extends model.Extends `json:"extends"`
}

func (*liveController) UpdateExtends(context *gin.Context) {
	log := logger.ReqLogger(context)
	updateExtendsRequest := &UpdateExtendsRequest{}
	err := context.BindJSON(updateExtendsRequest)
	if err != nil {
		log.Errorf("bind json failed, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "bind json failed",
			RequestId: log.ReqID(),
		})
		return
	}
	if updateExtendsRequest.LiveId == "" {
		log.Errorf("live_id is empty")
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "live_id is empty",
			RequestId: log.ReqID(),
		})
		return
	}
	err = live.GetService().UpdateExtends(context, updateExtendsRequest.LiveId, updateExtendsRequest.Extends)
	if err != nil {
		log.Errorf("update extends failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "update extends failed",
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

type UserInfo struct {
	UserId     string        `json:"user_id"`
	ImUserId   int64         `json:"im_userid"`
	ImUsername string        `json:"im_username"`
	Nick       string        `json:"nick"`
	Avatar     string        `json:"avatar"`
	Extends    model.Extends `json:"extends"`
}

type LiveUserListResponse struct {
	api.Response
	Data struct {
		TotalCount int        `json:"total_count"`
		PageTotal  int        `json:"page_total"`
		EndPage    bool       `json:"end_page"`
		List       []UserInfo `json:"list"`
	} `json:"data"`
}

func (*liveController) LiveUserList(context *gin.Context) {
	log := logger.ReqLogger(context)
	liveId := context.Query("live_id")
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	if liveId == "" {
		log.Errorf("live_id is empty")
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "live_id is empty",
			RequestId: log.ReqID(),
		})
		return
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		context.JSON(http.StatusBadRequest, api.Response{
			Code:      http.StatusBadRequest,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	liveUserList, totalCount, err := live.GetService().LiveUserList(context, liveId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("get live user list failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get live user list failed",
			RequestId: log.ReqID(),
		})
		return
	}
	userInfoList := make([]UserInfo, len(liveUserList), len(liveUserList))
	for i := range liveUserList {
		userInfo, err := user2.GetService().FindUser(context, liveUserList[i].UserId)
		if err != nil {
			log.Errorf("get user info failed, err: %v", err)
			continue
		}
		userInfoList[i].UserId = liveUserList[i].UserId
		userInfoList[i].ImUserId = userInfo.ImUserid
		userInfoList[i].ImUsername = userInfo.ImUsername
		userInfoList[i].Nick = userInfo.Nick
		userInfoList[i].Avatar = userInfo.Avatar
		userInfoList[i].Extends = userInfo.Extends
	}
	endPage := false
	if len(userInfoList) < pageSizeInt {
		endPage = true
	}
	response := &LiveUserListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = totalCount
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = userInfoList
	context.JSON(http.StatusOK, response)
}

type PutLikeRequest struct {
	Count int64 `json:"count"`
}

type PutLikeResponse struct {
	api.Response
	Data struct {
		Count int64 `json:"count"` //我在直播间内的点赞总数
		Total int64 `json:"total"` //直播间的点赞总数
	} `json:"data"`
}

func (*liveController) PutLike(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	userInfo := ctx.MustGet(liveauth.UserCtxKey).(*liveauth.UserInfo)
	req := PutLikeRequest{}
	ctx.ShouldBindJSON(&req)
	if req.Count == 0 {
		req.Count = 1
	}

	liveId := ctx.Param("live_id")
	liveInfo, err := live.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
	}

	my, total, err := live.GetService().AddLike(ctx, liveId, userInfo.UserId, req.Count)
	if err != nil {
		log.Errorf("add like error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	u, err := user2.GetService().FindUser(ctx, userInfo.UserId)
	if err == nil {
		item := &notify.LikeNotifyItem{
			LiveId: liveId,
			UserId: userInfo.UserId,
			Count:  req.Count,
		}
		go notify.SendNotifyToLive(ctx, u, liveInfo, notify.ActionTypeLikeNotify, item)
	}

	resp := &PutLikeResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data: struct {
			Count int64 `json:"count"`
			Total int64 `json:"total"`
		}{
			Count: my,
			Total: total,
		},
	}
	ctx.JSON(http.StatusOK, resp)
}
