package impl

import (
	"context"
	"time"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/biz/gift/config"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
)

var instance *ServiceImpl

type ServiceImpl struct {
	config.Config
}

func ConfigService(conf config.Config) {
	instance = &ServiceImpl{Config: conf}
}

func GetInstance() *ServiceImpl {
	return instance
}

func (s *ServiceImpl) SaveGiftEntity(context context.Context, entity *model.GiftEntity) error {
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
		return rest.ErrInternal
	}
	return nil
}

func (s *ServiceImpl) DeleteGiftEntity(context context.Context, giftId int) error {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	err := db.Delete(&model.GiftEntity{}, "gift_id = ?", giftId).Error
	if err != nil {
		log.Errorf("delete gift config error %s", err.Error())
		return rest.ErrInternal
	}
	return nil
}

func (s *ServiceImpl) GetListGiftEntity(context context.Context, typeId int) (entities []*model.GiftEntity, err error) {
	log := logger.ReqLogger(context)
	db := mysql.GetLive(log.ReqID())
	entities = make([]*model.GiftEntity, 0)
	if typeId == -1 {
		err = db.Model(&model.GiftEntity{}).Order("`order` asc").Order("created_at desc").Find(&entities).Error
	} else {
		err = db.Model(&model.GiftEntity{}).Where("type = ?", typeId).Order("`order` asc").Order("created_at desc").Find(&entities).Error
	}
	if err != nil {
		log.Errorf("list gift config error %s", err.Error())
		return nil, rest.ErrInternal
	}
	return entities, nil
}
