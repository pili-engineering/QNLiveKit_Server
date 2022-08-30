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
	mod.ID = 1
	err := db.Model(model.CensorConfig{}).Save(mod).Error
	if err != nil {
		return api.ErrDatabase
	}
	return nil
}

func (c *CensorService) GetCensorConfig(ctx *gin.Context) (*model.CensorConfig, error) {
	log := logger.ReqLogger(ctx)
	db := mysql.GetLive(log.ReqID())
	mod := &model.CensorConfig{}
	result := db.Model(model.CensorConfig{}).First(mod)
	if result.Error != nil {
		if result.RecordNotFound() {
			mod.Enable = true
			mod.ID = 1
			mod.Pulp = true
			mod.Interval = 20
			err := db.Model(model.CensorConfig{}).Save(mod).Error
			if err != nil {
				return nil, api.ErrDatabase
			}
		} else {
			return nil, api.ErrDatabase
		}
	}
	return mod, nil
}
