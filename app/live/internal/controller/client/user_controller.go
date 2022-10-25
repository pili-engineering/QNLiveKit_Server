// @Author: wangsheng
// @Description:
// @File:  user_controller
// @Version: 1.0.0
// @Date: 2022/5/19 2:47 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"net/http"

	"github.com/qbox/livekit/app/live/internal/dto"

	"github.com/qbox/livekit/biz/user"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterUserRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/user")
	userGroup.GET("/profile", UserController.GetUserProfile)
	userGroup.GET("/user/:id", UserController.GetUserInfo)
	userGroup.PUT("/user", UserController.PutUserInfo)
	userGroup.GET("/users", UserController.GetUsersInfo)
	userGroup.GET("/imusers", UserController.GetImUsersInfo)
}

var UserController = &userController{}

type userController struct {
}

type GetProfileResponse struct {
	api.Response
	Data *dto.UserProfileDto `json:"data"`
}

//获取用户自己的Profile 信息
// GET  /client/user/profile
func (*userController) GetUserProfile(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrNotFound))
		return
	}

	userService := user.GetService()
	userEntity, err := userService.FindOrCreateUser(ctx, uInfo.UserId)
	if err != nil {
		log.Errorf("find or create user %s error %v", uInfo.UserId, err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	dto := dto.User2ProfileDto(userEntity)
	resp := GetProfileResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: dto,
	}

	ctx.JSON(http.StatusOK, resp)
}

//获取用户自己的Profile 信息
// PUT  /client/user
func (*userController) PutUserInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrNotFound))
		return
	}

	updateInfo := dto.UserDto{}
	if err := ctx.BindJSON(&updateInfo); err != nil {
		log.Errorf("bind json error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	updateInfo.UserId = uInfo.UserId

	userService := user.GetService()
	err := userService.UpdateUserInfo(ctx, dto.UserDto2Entity(&updateInfo))
	if err != nil {
		log.Errorf("update user info error %v", err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type GetUserInfoResponse struct {
	api.Response
	Data *dto.UserDto `json:"data"`
}

//获取其他用户信息
// GET /client/user/:id
func (*userController) GetUserInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrBadToken))
		return
	}

	userId := ctx.Param("id")

	userService := user.GetService()
	userEntity, err := userService.FindUser(ctx, userId)
	if err != nil {
		log.Errorf("find  user %s error %v", userId, err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	dto := dto.User2Dto(userEntity)
	resp := GetUserInfoResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dto,
	}

	ctx.JSON(http.StatusOK, resp)
}

type GetUsersInfoRequest struct {
	UserIds []string `json:"user_ids" form:"user_ids"`
}

type GetUsersInfoResponse struct {
	api.Response
	Data []*dto.UserDto `json:"data"`
}

//批量获取其他用户信息
// GET /client/users
func (*userController) GetUsersInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrBadToken))
		return
	}

	req := GetUsersInfoRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind json error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(req.UserIds) == 0 || len(req.UserIds) > 200 {
		log.Errorf("ids cant be empty, or more than 200, %d", len(req.UserIds))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userService := user.GetService()
	users, err := userService.ListUser(ctx, req.UserIds)
	if err != nil {
		log.Errorf("find user userIds:%s error %v", req.UserIds, err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	dtos := make([]*dto.UserDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2Dto(u))
	}
	resp := GetUsersInfoResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dtos,
	}

	ctx.JSON(http.StatusOK, resp)
}

type GetImUsersInfoRequest struct {
	ImUserIds []int64 `json:"im_user_ids" form:"im_user_ids"`
}

//批量获取其他用户信息
// GET /client/user/imusers
func (*userController) GetImUsersInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrBadToken))
		return
	}

	req := GetImUsersInfoRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind json error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(req.ImUserIds) == 0 || len(req.ImUserIds) > 200 {
		log.Errorf("ids cant be empty, or more than 200, %d", len(req.ImUserIds))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userService := user.GetService()
	users, err := userService.ListImUser(ctx, req.ImUserIds)
	if err != nil {
		log.Errorf("find user imUserIds:%s error %v", req.ImUserIds, err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	dtos := make([]*dto.UserDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2Dto(u))
	}
	resp := GetUsersInfoResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dtos,
	}

	ctx.JSON(http.StatusOK, resp)
}
