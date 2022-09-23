// @Author: wangsheng
// @Description:
// @File:  service
// @Version: 1.0.0
// @Date: 2022/5/20 9:21 上午
// Copyright 2021 QINIU. All rights reserved

package user

import (
	"context"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/im"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/password"
	"github.com/qbox/livekit/utils/timestamp"
)

type IUserService interface {
	FindUser(ctx context.Context, userId string) (*model.LiveUserEntity, error)
	FindOrCreateUser(ctx context.Context, userId string) (*model.LiveUserEntity, error)
	ListUser(ctx context.Context, userIds []string) ([]*model.LiveUserEntity, error)
	ListImUser(ctx context.Context, imUserIds []int64) ([]*model.LiveUserEntity, error)
	CreateUser(ctx context.Context, user *model.LiveUserEntity) error
	UpdateUserInfo(ctx context.Context, user *model.LiveUserEntity) error
}

var userService IUserService = &UserService{}

func GetService() IUserService {
	return userService
}

type UserService struct {
}

func (s *UserService) FindUser(ctx context.Context, userId string) (*model.LiveUserEntity, error) {
	return s.findUser(ctx, userId)
}

func (s *UserService) FindOrCreateUser(ctx context.Context, userId string) (*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)
	ue, err := s.findUser(ctx, userId)
	if err == nil {
		return ue, err
	}

	if !api.IsNotFoundError(err) {
		log.Errorf("find user error %v", err)
		return nil, err
	}

	user := model.LiveUserEntity{
		UserId: userId,
	}

	return s.createUser(ctx, &user)
}

func (s *UserService) findUser(ctx context.Context, userId string) (*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)

	db := mysql.GetLiveReadOnly(log.ReqID())
	ue := &model.LiveUserEntity{}
	result := db.Model(model.LiveUserEntity{}).First(ue, "user_id = ?", userId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, api.ErrNotFound
		} else {
			return nil, api.ErrDatabase
		}
	}

	if ue.ImUserid == 0 {
		ue1, err := s.createImUser(ctx, ue)
		if err != nil {
			log.Errorf("create im user error %+v", err)
		} else {
			ue = ue1
		}
	}

	return ue, nil
}

func (s *UserService) ListUser(ctx context.Context, userIds []string) ([]*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)

	ret := make([]*model.LiveUserEntity, 0)
	db := mysql.GetLiveReadOnly(log.ReqID())

	result := db.Table(model.LiveUserEntity{}.TableName()).
		Where("user_id in (?)", userIds).
		Find(&ret)
	if result.Error != nil {
		log.Errorf("list user error %v", result.Error)
		return nil, api.ErrDatabase
	} else {
		return ret, nil
	}
}

func (s *UserService) ListImUser(ctx context.Context, imUserIds []int64) ([]*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)

	ret := make([]*model.LiveUserEntity, 0)
	db := mysql.GetLiveReadOnly(log.ReqID())

	result := db.Table(model.LiveUserEntity{}.TableName()).
		Where("im_userid in (?)", imUserIds).
		Find(&ret)
	if result.Error != nil {
		log.Errorf("list user error %v", result.Error)
		return nil, api.ErrDatabase
	} else {
		return ret, nil
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.LiveUserEntity) error {
	log := logger.ReqLogger(ctx)

	old, err := s.findUser(ctx, user.UserId)
	if err != nil && !api.IsNotFoundError(err) {
		log.Errorf("find user error %v", err)
		return api.ErrDatabase
	}

	if old != nil {
		return api.ErrAlreadyExist
	}

	_, err = s.createUser(ctx, user)
	if err != nil {
		log.Errorf("create user error %v", err)
	}

	return err
}
func (s *UserService) UpdateUserInfo(ctx context.Context, user *model.LiveUserEntity) error {
	log := logger.ReqLogger(ctx)

	db := mysql.GetLive(log.ReqID())
	db = db.Model(model.LiveUserEntity{}).Where("user_id = ?", user.UserId)

	old := model.LiveUserEntity{}
	result := db.First(&old)
	if result.Error != nil {
		log.Errorf("find old user error %+v", result.Error)
		if result.RecordNotFound() {
			return api.ErrNotFound
		} else {
			return api.ErrDatabase
		}
	}

	updates := user2Updates(&old, user)
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = timestamp.Now()
	result = db.Update(updates)
	if result.Error != nil {
		log.Errorf("update user error %v", result.Error)
		return api.ErrDatabase
	} else {
		return nil
	}
}

func user2Updates(oldUser, updateUser *model.LiveUserEntity) map[string]interface{} {
	updates := map[string]interface{}{}
	if len(updateUser.Nick) > 0 {
		updates["nick"] = updateUser.Nick
	}

	if len(updateUser.Avatar) > 0 {
		updates["avatar"] = updateUser.Avatar
	}

	if len(updateUser.Extends) > 0 {
		updates["extends"] = model.CombineExtends(oldUser.Extends, updateUser.Extends)
	}

	if updateUser.ImUserid > 0 {
		updates["im_userid"] = updateUser.ImUserid
	}

	if len(updateUser.ImUsername) > 0 {
		updates["im_username"] = updateUser.ImUsername
	}

	if len(updateUser.ImPassword) > 0 {
		updates["im_password"] = updateUser.ImPassword
	}

	return updates
}

func (s *UserService) createUser(ctx context.Context, user *model.LiveUserEntity) (*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)

	user.CreatedAt = timestamp.Now()
	user.UpdatedAt = timestamp.Now()

	db := mysql.GetLive(log.ReqID())
	result := db.Create(user)
	if result.Error != nil {
		log.Errorf("create user %+v, error %+v", result.Error)
		return nil, api.ErrDatabase
	}

	user1, err := s.createImUser(ctx, user)
	if err != nil {
		log.Errorf("create im user error %+v", err)
		return user, nil
	}
	user = user1

	return user, nil
}

func (s *UserService) createImUser(ctx context.Context, user *model.LiveUserEntity) (*model.LiveUserEntity, error) {
	log := logger.ReqLogger(ctx)

	imUsername := "ql" + user.UserId
	imPassword := password.RandomPassword(8)

	imService := im.GetService()
	imUserid, err := imService.RegisterUser(ctx, imUsername, imPassword)
	if err != nil {
		log.Errorf("register user error %+v", err)
		return nil, err
	}
	user.ImUserid = imUserid
	user.ImUsername = imUsername
	user.ImPassword = imPassword

	updateInfo := model.LiveUserEntity{
		UserId:     user.UserId,
		ImUserid:   imUserid,
		ImUsername: imUsername,
		ImPassword: imPassword,
	}

	s.UpdateUserInfo(ctx, &updateInfo)

	return user, nil
}
