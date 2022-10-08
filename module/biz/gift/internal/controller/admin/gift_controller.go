package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/biz/gift/dto"
	"github.com/qbox/livekit/module/biz/gift/internal/impl"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoute() {
	//giftGroup := group.Group("/gift")
	//giftGroup.GET("/config/:type", GiftConfigController.GetGiftConfig)
	//giftGroup.POST("/config", GiftConfigController.AddGiftConfig)
	//giftGroup.DELETE("/config/:gift_id", GiftConfigController.DeleteGiftConfig)

	httpq.AdminHandle(http.MethodGet, "/gift/config/:type", GiftConfigController.GetGiftConfig)
	httpq.AdminHandle(http.MethodPost, "/gift/config", GiftConfigController.AddGiftConfig)
	httpq.AdminHandle(http.MethodDelete, "/gift/config/:gift_id", GiftConfigController.DeleteGiftConfig)

}

type GiftCController struct {
}

var GiftConfigController = &GiftCController{}

// AddGiftConfig 添加礼物配置
// return nil
func (*GiftCController) AddGiftConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &dto.GiftConfigDto{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if req.GiftId == 0 || req.Name == "" {
		return nil, rest.ErrBadRequest.WithMessage("Name 不能为空 且 GiftId 需要 >0")
	}

	err := impl.GetInstance().SaveGiftEntity(ctx, dto.GiftDtoToEntity(req))
	if err != nil {
		log.Errorf("add gift config  failed, err: %v", err)
		return nil, rest.ErrInternal
	}

	return nil, nil
}

// DeleteGiftConfig 删除一个礼物配置
func (*GiftCController) DeleteGiftConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	typeId := ctx.Param("gift_id")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("gift_id is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("gift_id is not int")
	}
	err = impl.GetInstance().DeleteGiftEntity(ctx, typeIdInt)
	if err != nil {
		log.Errorf("delete gift config  failed, err: %v", err)
		return nil, rest.ErrInternal
	}

	return nil, nil
}

// GetGiftConfig 获取礼物配置
// return []*dto.GiftConfigDto
func (*GiftCController) GetGiftConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	typeId := ctx.Param("type")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("type is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("type is not int")
	}
	giftEntities, err := impl.GetInstance().GetListGiftEntity(ctx, typeIdInt)
	if err != nil {
		log.Errorf("get all gift config  failed, err: %v", err)
		return nil, rest.ErrInternal
	}
	var giftDtos []*dto.GiftConfigDto
	for _, v := range giftEntities {
		giftDtos = append(giftDtos, dto.GiftEntityToDto(v))
	}

	return giftDtos, nil
}
