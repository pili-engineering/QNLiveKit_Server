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
	sql := "insert into gift_config(type,name,amount,img,animation_type,animation_img,`order`,extends, created_at,updated_at) " +
		"VALUES(?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE type = ?,name = ?,amount = ?," +
		"img = ?,animation_type = ?,animation_img=?,`order`=?,extends=?, created_at=?,updated_at=?,deleted_at = NULL"
	err := db.Exec(sql, entity.Type, entity.Name, entity.Amount, entity.Img, entity.AnimationType, entity.AnimationImg, entity.Order, entity.Extends, time.Now(), time.Now(),
		entity.Type, entity.Name, entity.Amount, entity.Img, entity.AnimationType, entity.AnimationImg, entity.Order, entity.Extends, time.Now(), time.Now()).Error
	if err != nil {
		log.Errorf("add gift config error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *Service) DeleteGiftEntity(context context.Context, typeId int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Delete(&model.GiftEntity{}, "type = ?", typeId).Error
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
	err := db.Model(&model.GiftEntity{}).Order("`order` asc").Order("created_at desc").Find(&entities).Error
	if err != nil {
		log.Errorf("list gift config error %s", err.Error())
		return nil, api.ErrDatabase
	}
	return entities, nil
}
