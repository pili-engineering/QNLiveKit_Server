package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/mysql"
	"github.com/qbox/livekit/utils/logger"
)

type CCService interface {
	UpdateCensorConfig(ctx *gin.Context, mod *model.CensorConfig) error
	GetCensorConfig(ctx *gin.Context) (*model.CensorConfig, error)
}

type CensorService struct {
}

var cService CCService = &CensorService{}

func GetCensorService() CCService {
	return cService
}

func (c *CensorService) UpdateCensorConfig(ctx *gin.Context, mod *model.CensorConfig) error {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	var num int
	err := db.Model(model.CensorConfig{}).Count(&num).Error
	if err != nil {
		return api.ErrDatabase
	}
	if num > 0 {
		db.Delete(model.CensorConfig{}, "id = 1")
	}
	mod.ID = 1
	err = db.Model(model.CensorConfig{}).Save(mod).Error
	if err != nil {
		return api.ErrDatabase
	}
	return nil
}

func (c *CensorService) GetCensorConfig(ctx *gin.Context) (*model.CensorConfig, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	mod := &model.CensorConfig{}
	err := db.Model(model.CensorConfig{}).First(mod).Error
	if err != nil {
		return nil, api.ErrDatabase
	}
	return mod, nil
}
