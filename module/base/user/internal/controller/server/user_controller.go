// @Author: wangsheng
// @Description:
// @File:  user_controller
// @Version: 1.0.0
// @Date: 2022/6/28 9:14 上午
// Copyright 2021 QINIU. All rights reserved

package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/user/dto"
	"github.com/qbox/livekit/module/base/user/internal/impl"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	httpq.ServerHandle(http.MethodPost, "/user/register", UserController.PostRegister)
	httpq.ServerHandle(http.MethodPost, "/user/register/batch", UserController.PostRegisterBatch)
	httpq.ServerHandle(http.MethodPut, "/user/:id", UserController.PutUserInfo)
	httpq.ServerHandle(http.MethodGet, "/user/info/:id", UserController.GetUserInfo)
	httpq.ServerHandle(http.MethodGet, "/user/infos", UserController.GetUsersInfo)
	httpq.ServerHandle(http.MethodGet, "/user/profile/:id", UserController.GetUserProfile)
	httpq.ServerHandle(http.MethodGet, "/user/profiles", UserController.GetUsersProfile)
}

var UserController = &userController{}

type userController struct {
}

// PostRegister 注册一个用户
// return nil
func (c *userController) PostRegister(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	req := &dto.UserDto{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(req.UserId) == 0 {
		log.Errorf("invalid request, empty user_id")
		return nil, rest.ErrBadRequest.WithMessage("empty user_id")
	}

	userService := impl.GetService()
	userModel := model.LiveUserEntity{
		UserId:  req.UserId,
		Nick:    req.Nick,
		Avatar:  req.Avatar,
		Extends: req.Extends,
	}
	err := userService.CreateUser(ctx, &userModel)
	if err != nil {
		log.Errorf("create user  error: %v", err)
		return nil, err
	}

	return nil, nil
}

// PostRegisterBatch 批量注册用户
// return userId -> message 返回失败的部分信息
func (c *userController) PostRegisterBatch(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := make([]*dto.UserDto, 0)

	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(req) == 0 {
		log.Errorf("empty request")
		return nil, rest.ErrBadRequest.WithMessage("empty request")
	}

	if len(req) > 200 {
		log.Errorf("user count exceed 200, %d", len(req))
		return nil, rest.ErrBadRequest.WithMessage("user count exceed 200")
	}

	for i, dto := range req {
		if len(dto.UserId) == 0 {
			log.Errorf("invalid request, empty user_id, index %d", i)
			return nil, rest.ErrBadRequest.WithMessage(fmt.Sprintf("invalid request, empty user_id, index %d", i))
		}
	}

	userService := impl.GetService()
	failMap := make(map[string]string)
	for _, dto := range req {
		userModel := model.LiveUserEntity{
			UserId:  dto.UserId,
			Nick:    dto.Nick,
			Avatar:  dto.Avatar,
			Extends: dto.Extends,
		}
		err := userService.CreateUser(ctx, &userModel)
		if err != nil {
			failMap[dto.UserId] = err.Error()
		}
	}

	return failMap, nil
}

// GetUserInfo 查询用户信息
// GET /server/user/:id
func (*userController) GetUserInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")
	userService := impl.GetService()
	userEntity, err := userService.FindUser(ctx, userId)
	if err != nil {
		log.Errorf("find user %s error %v", userId, err)
		return nil, err
	}

	dto := dto.User2Dto(userEntity)
	return dto, nil
}

type GetUsersInfoRequest struct {
	UserIds []string `json:"user_ids" form:"user_ids"`
}

// GetUsersInfo 批量获取其他用户信息
// GET /server/users
// return []*dto.UserDto
func (*userController) GetUsersInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	req := GetUsersInfoRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind json error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(req.UserIds) == 0 || len(req.UserIds) > 200 {
		log.Errorf("ids cant be empty, or more than 200, %d", len(req.UserIds))
		return nil, rest.ErrBadRequest.WithMessage("ids cant be empty, or more than 200")
	}

	userService := impl.GetService()
	users, err := userService.ListUser(ctx, req.UserIds)
	if err != nil {
		log.Errorf("find user userIds:%s error %v", req.UserIds, err)
		return nil, err
	}

	dtos := make([]*dto.UserDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2Dto(u))
	}

	return dtos, nil
}

// GetUserProfile 获取用户Profile信息
// GET /server/user/profile/:id
// return *dto.UserProfileDto
func (*userController) GetUserProfile(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")
	userService := impl.GetService()
	userEntity, err := userService.FindUser(ctx, userId)
	if err != nil {
		log.Errorf("find user %s error %v", userId, err)
		return nil, err
	}

	dto := dto.User2ProfileDto(userEntity)
	return dto, nil
}

// GetUsersProfile 批量获取其他用户信息
// GET /server/user/profiles
// return []*dto.UserProfileDto
func (*userController) GetUsersProfile(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	req := GetUsersInfoRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		log.Errorf("bind json error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	if len(req.UserIds) == 0 || len(req.UserIds) > 200 {
		log.Errorf("ids cant be empty, or more than 200, %d", len(req.UserIds))
		return nil, rest.ErrInternal.WithMessage("ids cant be empty, or more than 200")
	}

	userService := impl.GetService()
	users, err := userService.ListUser(ctx, req.UserIds)
	if err != nil {
		log.Errorf("find user userIds:%s error %v", req.UserIds, err)
		return nil, err
	}

	dtos := make([]*dto.UserProfileDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2ProfileDto(u))
	}
	return dtos, nil
}

// PutUserInfo 修改用户信息
// return nil
func (c *userController) PutUserInfo(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")
	updateInfo := dto.UserDto{}
	if err := ctx.BindJSON(&updateInfo); err != nil {
		log.Errorf("bind json error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	updateInfo.UserId = userId

	userService := impl.GetService()
	err := userService.UpdateUserInfo(ctx, dto.UserDto2Entity(&updateInfo))
	if err != nil {
		log.Errorf("update user info error %v", err)
		return nil, err
	}

	return nil, nil
}
