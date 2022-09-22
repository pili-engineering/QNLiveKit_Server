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
	DeleteGiftEntity(context context.Context, giftId int) error
	GetListGiftEntity(context context.Context, typeId int) ([]*model.GiftEntity, error)
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
	sql := "insert into gift_config(gift_id,type,name," +
		"amount,img,animation_type,animation_img," +
		"`order`,extends, created_at, updated_at)  " +
		"VALUES(?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE gift_id = ?,type = ?,name = ?,amount = ?," +
		"img = ?,animation_type = ?,animation_img=?,`order`=?,extends=?, created_at=?,updated_at=?,deleted_at = NULL"
	err := db.Exec(sql, entity.GiftId, entity.Type, entity.Name,
		entity.Amount, entity.Img, entity.AnimationType, entity.AnimationImg,
		entity.Order, entity.Extends, time.Now(), time.Now(),
		entity.GiftId, entity.Type, entity.Name, entity.Amount, entity.Img, entity.AnimationType, entity.AnimationImg, entity.Order, entity.Extends, time.Now(), time.Now()).Error
	if err != nil {
		log.Errorf("add gift config error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *Service) DeleteGiftEntity(context context.Context, giftId int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Delete(&model.GiftEntity{}, "gift_id = ?", giftId).Error
	if err != nil {
		log.Errorf("delete gift config error %s", err.Error())
		return api.ErrDatabase
	}
	return nil
}

func (s *Service) GetListGiftEntity(context context.Context, typeId int) ([]*model.GiftEntity, error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	entities := make([]*model.GiftEntity, 0)
	err := db.Model(&model.GiftEntity{}).Where("type = ?", typeId).Order("`order` asc").Order("created_at desc").Find(&entities).Error
	if err != nil {
		log.Errorf("list gift config error %s", err.Error())
		return nil, api.ErrDatabase
	}
	return entities, nil
}
