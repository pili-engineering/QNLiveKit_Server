// @Author: wangsheng
// @Description:
// @File:  user_controller
// @Version: 1.0.0
// @Date: 2022/6/28 9:14 上午
// Copyright 2021 QINIU. All rights reserved

package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterUserRoutes(group *gin.RouterGroup) {
	userGroup := group.Group("/user")
	userGroup.POST("/register", UserController.PostRegister)
	userGroup.POST("/register/batch", UserController.PostRegisterBatch)
	userGroup.PUT("/:id", UserController.PutUserInfo)
	userGroup.GET("/info/:id", UserController.GetUserInfo)
	userGroup.GET("/infos", UserController.GetUsersInfo)
	userGroup.GET("/profile/:id", UserController.GetUserProfile)
	userGroup.GET("/profiles", UserController.GetUsersProfile)
}

var UserController = &userController{}

type userController struct {
}

func (c *userController) PostRegister(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	req := &dto.UserDto{}
	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(req.UserId) == 0 {
		log.Errorf("invalid request, empty user_id")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	userService := user.GetService()
	userModel := model.LiveUserEntity{
		UserId:  req.UserId,
		Nick:    req.Nick,
		Avatar:  req.Avatar,
		Extends: req.Extends,
	}
	err := userService.CreateUser(ctx, &userModel)
	if err != nil {
		log.Errorf("create user  error: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type RegisterBatchResponse struct {
	api.Response
	Data map[string]string `json:"data"`
}

func (c *userController) PostRegisterBatch(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := make([]*dto.UserDto, 0)

	if err := ctx.BindJSON(&req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(req) == 0 {
		log.Errorf("empty request")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	if len(req) > 200 {
		log.Errorf("user count exceed 200, %d", len(req))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	for i, dto := range req {
		if len(dto.UserId) == 0 {
			log.Errorf("invalid request, empty user_id, index %d", i)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
			return
		}
	}

	userService := user.GetService()
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

	resp := &RegisterBatchResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
	}
	if len(failMap) > 0 {
		resp.Data = failMap
	}

	ctx.JSON(http.StatusOK, resp)
}

type GetUserInfoResponse struct {
	api.Response
	Data *dto.UserDto `json:"data"`
}

// GET /server/user/:id
func (*userController) GetUserInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")

	userService := user.GetService()
	userEntity, err := userService.FindUser(ctx, userId)
	if err != nil {
		log.Errorf("find user %s error %v", userId, err)
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
// GET /server/users
func (*userController) GetUsersInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

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

type GetProfileResponse struct {
	api.Response
	Data *dto.UserProfileDto `json:"data"`
}

//获取用户Profile信息
// GET /server/user/profile/:id
func (*userController) GetUserProfile(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")

	userService := user.GetService()
	userEntity, err := userService.FindUser(ctx, userId)
	if err != nil {
		log.Errorf("find user %s error %v", userId, err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrInternal))
		return
	}

	dto := dto.User2ProfileDto(userEntity)
	resp := GetProfileResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dto,
	}

	ctx.JSON(http.StatusOK, resp)
}

type GetUsersProfileResponse struct {
	api.Response
	Data []*dto.UserProfileDto `json:"data"`
}

//批量获取其他用户信息
// GET /server/user/profiles
func (*userController) GetUsersProfile(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

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

	dtos := make([]*dto.UserProfileDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, dto.User2ProfileDto(u))
	}
	resp := GetUsersProfileResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     dtos,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *userController) PutUserInfo(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	userId := ctx.Param("id")

	updateInfo := dto.UserDto{}
	if err := ctx.BindJSON(&updateInfo); err != nil {
		log.Errorf("bind json error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	updateInfo.UserId = userId

	userService := user.GetService()
	err := userService.UpdateUserInfo(ctx, dto.UserDto2Entity(&updateInfo))
	if err != nil {
		log.Errorf("update user info error %v", err)
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}
