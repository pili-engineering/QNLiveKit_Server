package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/gift"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
	"strconv"
)

func RegisterGiftRoute(group *gin.RouterGroup) {
	giftGroup := group.Group("/gift")
	giftGroup.GET("/config/:type", GiftConfigController.GetGiftConfig)
	giftGroup.POST("/config", GiftConfigController.AddGiftConfig)
	giftGroup.DELETE("/config/:gift_id", GiftConfigController.DeleteGiftConfig)
}

type GiftCController struct {
}

var GiftConfigController = &GiftCController{}

func (*GiftCController) AddGiftConfig(context *gin.Context) {
	log := logger.ReqLogger(context)
	req := &dto.GiftConfigDto{}
	if err := context.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		context.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if req.GiftId == 0 || req.Name == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, &api.Response{
			Code:      http.StatusBadRequest,
			RequestId: log.ReqID(),
			Message:   "Name 不能为空 且 GiftId 需要 >0 ",
		})
		return
	}

	err := gift.GetService().SaveGiftEntity(context, dto.GiftDtoToEntity(req))
	if err != nil {
		log.Errorf("add gift config  failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "add gift config failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

func (*GiftCController) DeleteGiftConfig(context *gin.Context) {
	log := logger.ReqLogger(context)
	typeId := context.Param("gift_id")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("gift_id is not int, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "gift_id is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	err = gift.GetService().DeleteGiftEntity(context, typeIdInt)
	if err != nil {
		log.Errorf("delete gift config  failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "delete gift config failed",
			RequestId: log.ReqID(),
		})
		return
	}
	context.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

type ListGiftConfigResponse struct {
	api.Response
	Data []*dto.GiftConfigDto `json:"data"`
}

func (*GiftCController) GetGiftConfig(context *gin.Context) {
	log := logger.ReqLogger(context)
	typeId := context.Param("type")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("type is not int, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "type is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	giftEntities, err := gift.GetService().GetListGiftEntity(context, typeIdInt)
	if err != nil {
		log.Errorf("get all gift config  failed, err: %v", err)
		context.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "get all gift config failed",
			RequestId: log.ReqID(),
		})
		return
	}
	var giftDtos []*dto.GiftConfigDto
	for _, v := range giftEntities {
		giftDtos = append(giftDtos, dto.GiftEntityToDto(v))
	}
	response := &ListGiftConfigResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data = giftDtos
	context.JSON(http.StatusOK, response)
}
