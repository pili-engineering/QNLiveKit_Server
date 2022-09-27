// @Author: wangsheng
// @Description:
// @File:  user_controller
// @Version: 1.0.0
// @Date: 2022/5/19 2:47 下午
// Copyright 2021 QINIU. All rights reserved

package client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/user/dto"
	"github.com/qbox/livekit/module/base/user/internal/impl"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	httpq.ClientHandle(http.MethodGet, "/user/profile", UserController.GetUserProfile)
	httpq.ClientHandle(http.MethodGet, "/user/user/:id", UserController.GetUserInfo)
	httpq.ClientHandle(http.MethodPut, "/user/user", UserController.PutUserInfo)
	httpq.ClientHandle(http.MethodGet, "/user/users", UserController.GetUserInfo)
	httpq.ClientHandle(http.MethodGet, "/user/imusers", UserController.GetImUsersInfo)
}

var UserController = &userController{}

type userController struct {
}

// GetUserProfile 获取用户自己的Profile 信息
// GET  /client/user/profile
// return *dto.UserProfileDto
func (*userController) GetUserProfile(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)

	userService := impl.GetService()
	userEntity, err := userService.FindUser(ctx, uInfo.UserId)
	if err != nil {
		log.Errorf("find  user %s error %v", uInfo.UserId, err)
		return nil, rest.ErrInternal
	}

	dto := dto.User2ProfileDto(userEntity)
	return dto, nil
}

// PutUserInfo 更新用户自己的信息
// PUT  /client/user
// 无返回值
func (*userController) PutUserInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	uInfo := liveauth.GetUserInfo(ctx)

	updateInfo := dto.UserDto{}
	if err := ctx.BindJSON(&updateInfo); err != nil {
		log.Errorf("bind json error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	updateInfo.UserId = uInfo.UserId

	userService := impl.GetService()
	err := userService.UpdateUserInfo(ctx, dto.UserDto2Entity(&updateInfo))
	if err != nil {
		log.Errorf("update user info error %v", err)
		return nil, rest.ErrInternal
	}

	return nil, nil
}

// GetUserInfo 获取其他用户信息
// GET /client/user/:id
// return *dto.UserDto
func (*userController) GetUserInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")
	userService := impl.GetService()
	userEntity, err := userService.FindUser(ctx, userId)
	if err != nil {
		log.Errorf("find  user %s error %v", userId, err)
		return nil, err
	}

	dto := dto.User2Dto(userEntity)
	return dto, nil
}

type GetUsersInfoRequest struct {
	UserIds []string `json:"user_ids" form:"user_ids"`
}

// GetUsersInfo 批量获取其他用户信息
// GET /client/users
// return []*dto.UserInfo
func (*userController) GetUsersInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	req := GetUsersInfoRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind json error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(req.UserIds) == 0 || len(req.UserIds) > 200 {
		log.Errorf("ids cant be empty, or more than 200, %d", len(req.UserIds))
		return nil, rest.ErrBadRequest.WithMessage("ids cant be empty, or more than 200,")
	}

	userService := impl.GetService()
	users, err := userService.ListUser(ctx, req.UserIds)
	if err != nil {
		log.Errorf("find user userIds:%s error %v", req.UserIds, err)
		return nil, rest.ErrInternal
	}

	dtos := make([]*dto.UserDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2Dto(u))
	}
	return dtos, nil
}

type GetImUsersInfoRequest struct {
	ImUserIds []int64 `json:"im_user_ids" form:"im_user_ids"`
}

// GetImUsersInfo 根据IM用户批量获取其他用户信息
// GET /client/user/imusers
// return []*dto.UserInfo
func (*userController) GetImUsersInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := GetImUsersInfoRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind json error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(req.ImUserIds) == 0 || len(req.ImUserIds) > 200 {
		log.Errorf("ids cant be empty, or more than 200, %d", len(req.ImUserIds))
		return nil, rest.ErrBadRequest.WithMessage("ids cant be empty, or more than 200")
	}

	userService := impl.GetService()
	users, err := userService.ListImUser(ctx, req.ImUserIds)
	if err != nil {
		log.Errorf("find user imUserIds:%s error %v", req.ImUserIds, err)
		return nil, rest.ErrInternal
	}

	dtos := make([]*dto.UserDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2Dto(u))
	}
	return dtos, nil
}
