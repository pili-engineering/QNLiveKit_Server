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

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/live/dto"
	"github.com/qbox/livekit/module/base/live/internal/impl"
	"github.com/qbox/livekit/module/base/live/service"
	"github.com/qbox/livekit/module/base/stats"
	"github.com/qbox/livekit/module/base/user"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	httpq.ClientHandle(http.MethodPost, "/live/room/instance", LiveController.CreateLive)
	httpq.ClientHandle(http.MethodDelete, "/live/room/instance/:live_id", LiveController.DeleteLive)
	httpq.ClientHandle(http.MethodGet, "/live/room/info/:live_id", LiveController.LiveRoomInfo)
	httpq.ClientHandle(http.MethodDelete, "/live/room/:live_id", LiveController.StopLive)
	httpq.ClientHandle(http.MethodPut, "/live/room/:live_id", LiveController.StartLive)
	httpq.ClientHandle(http.MethodGet, "/live/room", LiveController.SearchLive)
	httpq.ClientHandle(http.MethodPost, "/live/room/user/:live_id", LiveController.JoinLive)
	httpq.ClientHandle(http.MethodGet, "/live/room/list", LiveController.LiveList)
	httpq.ClientHandle(http.MethodGet, "/live/room/list/anchor", LiveController.LiveListAnchor)
	httpq.ClientHandle(http.MethodDelete, "/live/room/user/:live_id", LiveController.LeaveLive)
	httpq.ClientHandle(http.MethodGet, "/live/room/heartbeat/:live_id", LiveController.Heartbeat)
	httpq.ClientHandle(http.MethodPut, "/live/room/extends", LiveController.UpdateExtends)
	httpq.ClientHandle(http.MethodGet, "/live/room/user_list", LiveController.LiveUserList)
	httpq.ClientHandle(http.MethodPut, "/live/room/:live_id/like", LiveController.PutLike)
}

type liveController struct {
}

var LiveController = &liveController{}

// CreateLive 创建一个直播
// return dto.LiveInfoDto
func (*liveController) CreateLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	request := &service.CreateLiveRequest{}
	if err := context.Bind(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	request.AnchorId = userInfo.UserId
	liveEntity, err := impl.GetInstance().CreateLive(context, request)
	if err != nil {
		log.Errorf("create liveEntity failed, err: %v", err)
		return nil, err
	}

	userEntity, err := user.GetService().FindUser(context, userInfo.UserId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		return nil, err
	}

	return dto.BuildLiveDto(liveEntity, userEntity), nil
}

// DeleteLive 删除一个直播间
// return nil
func (*liveController) DeleteLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	err := impl.GetInstance().DeleteLive(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("delete live failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}

// LiveRoomInfo 查询直播间信息
// return dto.LiveInfoDto
func (c *liveController) LiveRoomInfo(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	liveInfo, err := impl.GetInstance().LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		return nil, err
	}
	userEntity, err := user.GetService().FindUser(context, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		return nil, err
	}
	_, onlineCount, err := impl.GetInstance().LiveUserList(context, liveId, 1, 10)
	if err != nil {
		log.Errorf("get live user list failed, err: %v", err)
		return nil, err
	}

	anchorStatus, err := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("get live anchor status error: %v", err)
		return nil, err
	}

	ret := dto.BuildLiveDto(liveInfo, userEntity)
	ret.AnchorStatus = anchorStatus
	ret.OnlineCount = onlineCount

	return ret, nil
}

func (*liveController) getLiveAnchorStatus(ctx context.Context, liveId string, anchorId string) (model.LiveRoomUserStatus, error) {
	log := logger.ReqLogger(ctx)
	liveUser, err := impl.GetInstance().FindLiveRoomUser(ctx, liveId, anchorId)
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

// StopLive 停止一个直播间
// return nil
func (*liveController) StopLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	err := impl.GetInstance().StopLive(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("stop live failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}

// StartLive 开始直播
// return dto.LiveInfoDto
func (*liveController) StartLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	roomToken, err := impl.GetInstance().StartLive(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("start live failed, err: %v", err)
		return nil, err
	}
	liveInfo, err := impl.GetInstance().LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		return nil, err
	}
	userEntity, err := user.GetService().FindUser(context, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		return nil, err
	}

	liveDto := dto.BuildLiveDto(liveInfo, userEntity)
	liveDto.RoomToken = roomToken

	return liveDto, nil
}

// SearchLive 查询直播间
// return rest.PageResult<*dto.LiveInfoDto>
func (c *liveController) SearchLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	keyword := context.Query("keyword")
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page_num is not int")
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page_size is not int")
	}
	liveList1, totalCount1, err := impl.GetInstance().SearchLive(context, keyword, 1, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search live failed, err: %v", err)
		return nil, err
	}
	liveList2, totalCount2, err := impl.GetInstance().SearchLive(context, keyword, 2, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search live failed, err: %v", err)
		return nil, err
	}
	liveList3, totalCount3, err := impl.GetInstance().SearchLive(context, keyword, 3, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("search live failed, err: %v", err)
		return nil, err
	}
	liveList := append(liveList1)
	liveList = append(liveList2)
	liveList = append(liveList3)
	endPage := false
	if len(liveList) < pageSizeInt {
		endPage = true
	}

	list := make([]*dto.LiveInfoDto, 0, len(liveList))
	for i := range liveList {
		liveInfo, err := impl.GetInstance().LiveInfo(context, liveList[i].LiveId)
		if err != nil {
			log.Errorf("get liveInfo info failed, err: %v", err)
			continue
		}
		userEntity, err := user.GetService().FindUser(context, userInfo.UserId)
		if err != nil {
			log.Errorf("find user failed, err: %v", err)
			continue
		}
		anchorStatus, _ := c.getLiveAnchorStatus(context, liveList[i].LiveId, liveList[i].AnchorId)
		liveDto := dto.BuildLiveDto(liveInfo, userEntity)
		liveDto.AnchorStatus = anchorStatus
		list = append(list, &liveDto)
	}

	totalCount := totalCount1 + totalCount2 + totalCount3
	pageTotal := int(math.Ceil(float64(totalCount) / float64(pageSizeInt)))
	result := rest.PageResult{
		TotalCount: totalCount,
		PageTotal:  pageTotal,
		EndPage:    endPage,
		List:       list,
	}
	return result, nil
}

// JoinLive 用户加入直播间
// return dto.LiveInfoDto
func (c *liveController) JoinLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	err := impl.GetInstance().JoinLiveRoom(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("get live info failed, err: %v", err)
		return nil, err
	}
	liveInfo, err := impl.GetInstance().LiveInfo(context, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		return nil, err
	}
	userEntity, err := user.GetService().FindUser(context, liveInfo.AnchorId)
	if err != nil {
		log.Errorf("find user failed, err: %v", err)
		return nil, err
	}
	rService := stats.GetService()
	statsSingleLiveEntity := &model.StatsSingleLiveEntity{
		LiveId: liveId,
		UserId: userInfo.UserId,
		Type:   1,
		Count:  1,
	}
	rService.UpdateSingleLive(context, statsSingleLiveEntity)
	anchorStatus, _ := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)

	result := dto.BuildLiveDto(liveInfo, userEntity)
	result.AnchorStatus = anchorStatus
	return result, nil
}

// LiveList 列出
// return rest.PageResult<*dto.LiveInfoDto>
func (c *liveController) LiveList(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page num is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page num is not int")
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page size is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page size is not int")
	}
	if pageNumInt <= 0 || pageSizeInt <= 0 {
		log.Errorf("page num or page size is not right, page num: %v, page size: %v", pageNumInt, pageSizeInt)
		return nil, rest.ErrBadRequest.WithMessage("page num or page size is not right")
	}

	liveList, totalCount, err := impl.GetInstance().LiveList(context, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("get live list failed, err: %v", err)
		return nil, err
	}
	endPage := false
	if len(liveList) < pageSizeInt {
		endPage = true
	}

	result := rest.PageResult{
		TotalCount: totalCount,
		PageTotal:  int(math.Ceil(float64(totalCount) / float64(pageSizeInt))),
		EndPage:    endPage,
	}

	list := make([]*dto.LiveInfoDto, 0, len(liveList))
	for i := range liveList {
		liveInfo, err := impl.GetInstance().LiveInfo(context, liveList[i].LiveId)
		if err != nil {
			log.Errorf("get liveInfo info failed, err: %v", err)
			continue
		}
		log.Infof("liveInfo: %v", liveInfo)
		userEntity, err := user.GetService().FindUser(context, liveInfo.AnchorId)
		if err != nil {
			log.Errorf("find user failed, err: %v", err)
			continue
		}
		anchorStatus, _ := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)

		liveDto := dto.BuildLiveDto(liveInfo, userEntity)
		liveDto.AnchorStatus = anchorStatus

		list = append(list, &liveDto)
	}
	result.List = list

	return &result, nil
}

// LiveListAnchor 查看主播自己的直播间
// return rest.PageResult<*dto.LiveInfoDto>
func (c *liveController) LiveListAnchor(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	uInfo := auth.GetUserInfo(context)

	anchorId := uInfo.UserId
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page num is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page num is not int")
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page size is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page size is not int")
	}
	if pageNumInt <= 0 || pageSizeInt <= 0 {
		log.Errorf("page num or page size is not right, page num: %v, page size: %v", pageNumInt, pageSizeInt)
		return nil, rest.ErrBadRequest.WithMessage("page num or page size is not right")
	}

	liveList, totalCount, err := impl.GetInstance().LiveListAnchor(context, pageNumInt, pageSizeInt, anchorId)
	if err != nil {
		log.Errorf("get live list anchor failed, err: %v", err)
		return nil, err
	}
	endPage := false
	if len(liveList) < pageSizeInt {
		endPage = true
	}

	result := rest.PageResult{
		TotalCount: totalCount,
		PageTotal:  int(math.Ceil(float64(totalCount) / float64(pageSizeInt))),
		EndPage:    endPage,
	}

	list := make([]*dto.LiveInfoDto, 0, len(liveList))
	for i := range liveList {
		liveInfo, err := impl.GetInstance().LiveInfo(context, liveList[i].LiveId)
		if err != nil {
			log.Errorf("get liveInfo info failed, err: %v", err)
			continue
		}
		log.Infof("liveInfo: %v", liveInfo)
		userEntity, err := user.GetService().FindUser(context, liveInfo.AnchorId)
		log.Infof("user: %v", userEntity)
		if err != nil {
			log.Errorf("find user failed, err: %v", err)
			continue
		}
		anchorStatus, _ := c.getLiveAnchorStatus(context, liveInfo.LiveId, liveInfo.AnchorId)

		liveDto := dto.BuildLiveDto(liveInfo, userEntity)
		liveDto.AnchorStatus = anchorStatus

		list = append(list, &liveDto)
	}
	result.List = list
	return &result, nil
}

// LeaveLive 用户离开直播间
// return nil
func (*liveController) LeaveLive(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	if liveId == "" {
		log.Errorf("live_id is empty")
		return nil, rest.ErrBadRequest.WithMessage("live_id is empty")
	}
	err := impl.GetInstance().LeaveLiveRoom(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("leave live room failed, err: %v", err)
		return nil, err
	}

	return nil, nil
}

type HeartBeatResult struct {
	LiveId string `json:"live_id"`
	Status int    `json:"live_status"`
}

// Heartbeat 心跳
// return *HeartBeatResult
func (*liveController) Heartbeat(context *gin.Context) (interface{}, error) {
	userInfo := auth.GetUserInfo(context)
	log := logger.ReqLogger(context)
	liveId := context.Param("live_id")
	if liveId == "" {
		log.Errorf("live_id is empty")
		return nil, rest.ErrBadRequest.WithMessage("live_id is empty")
	}
	liveEntity, err := impl.GetInstance().Heartbeat(context, liveId, userInfo.UserId)
	if err != nil {
		log.Errorf("heartbeat failed, err: %v", err)
		return nil, err
	}

	return &HeartBeatResult{
		LiveId: liveEntity.LiveId,
		Status: liveEntity.Status,
	}, nil
}

type UpdateExtendsRequest struct {
	LiveId  string        `json:"live_id"`
	Extends model.Extends `json:"extends"`
}

// UpdateExtends 更新直播的扩展信息
// return nil
func (*liveController) UpdateExtends(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	updateExtendsRequest := &UpdateExtendsRequest{}
	err := context.BindJSON(updateExtendsRequest)
	if err != nil {
		log.Errorf("bind json failed, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if updateExtendsRequest.LiveId == "" {
		log.Errorf("live_id is empty")
		return nil, rest.ErrBadRequest.WithMessage("live_id is empty")
	}

	err = impl.GetInstance().UpdateExtends(context, updateExtendsRequest.LiveId, updateExtendsRequest.Extends)
	if err != nil {
		log.Errorf("update extends failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}

type UserInfo struct {
	UserId     string        `json:"user_id"`
	ImUserId   int64         `json:"im_userid"`
	ImUsername string        `json:"im_username"`
	Nick       string        `json:"nick"`
	Avatar     string        `json:"avatar"`
	Extends    model.Extends `json:"extends"`
}

// LiveUserList 获取直播间用户信息
// return rest.PageResult<*UserInfo>
func (*liveController) LiveUserList(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	liveId := context.Query("live_id")
	pageNum := context.DefaultQuery("page_num", "1")
	pageSize := context.DefaultQuery("page_size", "10")
	if liveId == "" {
		log.Errorf("live_id is empty")
		return nil, rest.ErrBadRequest.WithMessage("live_id is empty")
	}

	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page_num is not int")
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("page_size is not int")
	}

	liveUserList, totalCount, err := impl.GetInstance().LiveUserList(context, liveId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("get live user list failed, err: %v", err)
		return nil, err
	}

	userInfoList := make([]*UserInfo, 0, len(liveUserList))
	for i := range liveUserList {
		userInfo, err := user.GetService().FindUser(context, liveUserList[i].UserId)
		if err != nil {
			log.Errorf("get user info failed, err: %v", err)
			continue
		}

		userInfoList = append(userInfoList, &UserInfo{
			UserId:     liveUserList[i].UserId,
			ImUserId:   userInfo.ImUserid,
			ImUsername: userInfo.ImUsername,
			Nick:       userInfo.Nick,
			Avatar:     userInfo.Avatar,
			Extends:    userInfo.Extends,
		})
	}
	endPage := false
	if len(userInfoList) < pageSizeInt {
		endPage = true
	}

	result := rest.PageResult{
		TotalCount: totalCount,
		PageTotal:  int(math.Ceil(float64(totalCount) / float64(pageSizeInt))),
		EndPage:    endPage,
		List:       userInfoList,
	}
	return &result, nil
}

type PutLikeRequest struct {
	Count int64 `json:"count"`
}

type PutLikeResult struct {
	Count int64 `json:"count"` //我在直播间内的点赞总数
	Total int64 `json:"total"` //直播间的点赞总数
}

// PutLike 用户在直播间点赞
// return *PutLikeResult
func (*liveController) PutLike(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	userInfo := auth.GetUserInfo(ctx)
	req := PutLikeRequest{}
	ctx.ShouldBindJSON(&req)
	if req.Count == 0 {
		req.Count = 1
	}

	liveId := ctx.Param("live_id")
	liveInfo, err := impl.GetInstance().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("get liveInfo info failed, err: %v", err)
		return nil, err
	}

	my, total, err := impl.GetInstance().AddLike(ctx, liveId, userInfo.UserId, req.Count)
	if err != nil {
		log.Errorf("add like error %s", err.Error())
		return nil, err
	}

	u, err := user.GetService().FindUser(ctx, userInfo.UserId)
	if err == nil {
		item := &notify.LikeNotifyItem{
			LiveId: liveId,
			UserId: userInfo.UserId,
			Count:  req.Count,
		}
		go notify.SendNotifyToLive(ctx, u, liveInfo, notify.ActionTypeLikeNotify, item)
	}

	return &PutLikeResult{
		Count: my,
		Total: total,
	}, nil
}
