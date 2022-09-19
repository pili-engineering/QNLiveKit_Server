package gift

import (
	"context"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
	"time"
)

type GService interface {
	SaveGiftEntity(context context.Context, entity *model.GiftEntity) error
	DeleteGiftEntity(context context.Context, typeId int) error
	GetListGiftEntity(context context.Context) ([]*model.GiftEntity, error)
}

type Service struct {
}

var service GService = &Service{}

func GetService() GService {
	return service
}

func (s *Service) SaveGiftEntity(context context.Context, entity *model.GiftEntity) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(&model.GiftEntity{}).Save(entity).Error
	if err != nil {
		log.Errorf("add gift config error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *Service) DeleteGiftEntity(context context.Context, typeId int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Model(&model.GiftEntity{}).Where("type = ?", typeId).Update("deleted_at", time.Now()).Error
	if err != nil {
		log.Errorf("delete gift config error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *Service) GetListGiftEntity(context context.Context) ([]*model.GiftEntity, error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	entities := make([]*model.GiftEntity, 0)
	err := db.Find(&entities).Error
	if err != nil {
		log.Errorf("list gift config error %s", err.Error())
		return nil, api.ErrDatabase
	}
	return entities, nil
}
