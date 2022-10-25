package server

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/gift"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
)

type giftController struct {
}

var GiftController = &giftController{}

func RegisterGiftRoutes(group *gin.RouterGroup) {
	giftGroup := group.Group("/gift")

	giftGroup.GET("/config/:type", GiftController.GetGiftConfig)
	giftGroup.POST("/config", GiftController.AddGiftConfig)
	giftGroup.DELETE("/config/:gift_id", GiftController.DeleteGiftConfig)

	giftGroup.GET("/list/live/:live_id", GiftController.ListGiftByLiveId)
	giftGroup.GET("/list/anchor/:anchor_id", GiftController.ListGiftByAnchorId)
	giftGroup.GET("/list/user/:user_id", GiftController.ListGiftByUserId)
}

func (*giftController) AddGiftConfig(context *gin.Context) {
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

func (*giftController) DeleteGiftConfig(context *gin.Context) {
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
	Data []dto.GiftConfigDto `json:"data"`
}

func (*giftController) GetGiftConfig(context *gin.Context) {
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
	giftDtos := make([]dto.GiftConfigDto, 0)
	for _, v := range giftEntities {
		giftDtos = append(giftDtos, *dto.GiftEntityToDto(v))
	}
	response := &ListGiftConfigResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data = giftDtos
	context.JSON(http.StatusOK, response)
}

func (c *giftController) ListGiftByLiveId(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("live_id")
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}

	gifts, count, err := gift.GetService().SearchGiftByLiveId(ctx, liveId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("Search Gift By LiveId failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "Search Gift ByLiveId failed",
			RequestId: log.ReqID(),
		})
		return
	}

	endPage := false
	if len(gifts) < pageSizeInt {
		endPage = true
	}
	response := &LiveGiftListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = count
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = gifts
	ctx.JSON(http.StatusOK, response)
}

func (c *giftController) ListGiftByAnchorId(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	anchorId := ctx.Param("anchor_id")
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	gifts, count, err := gift.GetService().SearchGiftByAnchorId(ctx, anchorId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("Search Gift by AnchorId  failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "Search Gift by AnchorId  failed",
			RequestId: log.ReqID(),
		})
		return
	}

	endPage := false
	if len(gifts) < pageSizeInt {
		endPage = true
	}
	response := &LiveGiftListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = count
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = gifts
	ctx.JSON(http.StatusOK, response)
}

func (c *giftController) ListGiftByUserId(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	userId := ctx.Param("user_id")
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_num is not int",
			RequestId: log.ReqID(),
		})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "page_size is not int",
			RequestId: log.ReqID(),
		})
		return
	}

	gifts, count, err := gift.GetService().SearchGiftByUserId(ctx, userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("Search Gift ByUserId  failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.Response{
			Code:      http.StatusInternalServerError,
			Message:   "Search Gift ByUserId  failed",
			RequestId: log.ReqID(),
		})
		return
	}

	endPage := false
	if len(gifts) < pageSizeInt {
		endPage = true
	}
	response := &LiveGiftListResponse{}
	response.Response.Code = 200
	response.Response.Message = "success"
	response.Response.RequestId = log.ReqID()
	response.Data.TotalCount = count
	response.Data.PageTotal = int(math.Ceil(float64(response.Data.TotalCount) / float64(pageSizeInt)))
	response.Data.EndPage = endPage
	response.Data.List = gifts
	ctx.JSON(http.StatusOK, response)
}

type LiveGiftListResponse struct {
	api.Response
	Data struct {
		TotalCount int               `json:"total_count"`
		PageTotal  int               `json:"page_total"`
		EndPage    bool              `json:"end_page"`
		List       []*model.LiveGift `json:"list"`
	} `json:"data"`
}
