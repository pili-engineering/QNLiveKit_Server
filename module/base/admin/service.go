package admin

import (
	"context"
	"errors"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
)

type ManagerService interface {
	FindAdminByUserName(ctx context.Context, userName string) (*model.ManagerEntity, error)
	FindAdminByUserId(ctx context.Context, userId string) (*model.ManagerEntity, error)
	FindOrCreateAdminUser(ctx context.Context, userId string) (*model.ManagerEntity, error)
}

type ManService struct {
}

var aService ManagerService = &ManService{}

func GetManagerService() ManagerService {
	return aService
}

func (s *ManService) FindAdminByUserName(ctx context.Context, userName string) (*model.ManagerEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	me := &model.ManagerEntity{}
	result := db.Model(model.ManagerEntity{}).First(me, "user_name = ?", userName)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, ErrUserPassword
		} else {
			log.Errorf("find admin by user name error %v", result.Error)
			return nil, rest.ErrInternal
		}
	}
	return me, nil
}

func (s *ManService) FindAdminByUserId(ctx context.Context, userId string) (*model.ManagerEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	me := &model.ManagerEntity{}
	result := db.Model(model.ManagerEntity{}).First(me, "user_id = ?", userId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, rest.ErrNotFound
		} else {
			log.Errorf("find user by id error %v", result.Error)
			return nil, rest.ErrInternal
		}
	}
	return me, nil
}

func (s *ManService) FindOrCreateAdminUser(ctx context.Context, userId string) (*model.ManagerEntity, error) {
	log := logger.ReqLogger(ctx)
	entity, err := s.FindAdminByUserId(ctx, userId)
	if err != nil && errors.Is(err, rest.ErrNotFound) {
		db := mysql.GetLive(log.ReqID())
		m := &model.ManagerEntity{
			UserId: userId,
		}
		err = db.Model(model.ManagerEntity{}).Save(m).Error
		if err != nil {
			log.Errorf("create admin error %v", err)
			return nil, rest.ErrInternal
		}
		return m, nil
	} else if err != nil {
		return nil, err
	} else {
		return entity, nil
	}
}
