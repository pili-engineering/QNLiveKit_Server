package admin

import (
	"context"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
)

type ManagerService interface {
	LoginManager(ctx context.Context, userId string, passWord string) (*model.ManagerEntity, error)
	FindManager(ctx context.Context, userId string) (*model.ManagerEntity, error)
}

type ManService struct {
}

var aService ManagerService = &ManService{}

func GetManagerService() ManagerService {
	return aService
}

func (s *ManService) LoginManager(ctx context.Context, userId string, passWord string) (*model.ManagerEntity, error) {
	manager, err := s.FindManager(ctx, userId)
	if err != nil {
		return nil, err
	}
	if manager.PassWord == passWord {
		return manager, nil
	} else {
		return nil, api.ErrorPassWordWrong
	}
}

func (s *ManService) FindManager(ctx context.Context, userId string) (*model.ManagerEntity, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	me := &model.ManagerEntity{}
	result := db.Model(model.ManagerEntity{}).First(me, "user_id = ?", userId)
	if result.Error != nil {
		if result.RecordNotFound() {
			return nil, api.ErrNotFound
		} else {
			return nil, api.ErrDatabase
		}
	}
	return me, nil
}
