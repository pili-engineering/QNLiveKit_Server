package client

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/gift"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/report"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
)

type giftController struct {
}

var GiftController = &giftController{}

func RegisterGiftRoutes(group *gin.RouterGroup) {
	giftGroup := group.Group("/gift")
	giftGroup.GET("/config/:type", GiftController.GetGiftConfig)

	giftGroup.GET("/list/live/:live_id", GiftController.ListGiftByLiveId)
	giftGroup.GET("/list/anchor", GiftController.ListGiftByAnchorId)
	giftGroup.GET("/list/user", GiftController.ListGiftByUserId)
	giftGroup.POST("/send", GiftController.SendGift)

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

type ListGiftConfigResponse struct {
	api.Response
	Data []dto.GiftConfigDto `json:"data"`
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
	userInfo := liveauth.GetUserInfo(ctx)
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
	gifts, count, err := gift.GetService().SearchGiftByAnchorId(ctx, userInfo.UserId, pageNumInt, pageSizeInt)
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
	userInfo := liveauth.GetUserInfo(ctx)
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

	gifts, count, err := gift.GetService().SearchGiftByUserId(ctx, userInfo.UserId, pageNumInt, pageSizeInt)
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

func (c *giftController) SendGift(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &gift.SendGiftRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	uInfo := liveauth.GetUserInfo(ctx)
	if uInfo == nil {
		log.Errorf("user info not exist")
		ctx.AbortWithStatusJSON(http.StatusNotFound, api.ErrorWithRequestId(log.ReqID(), api.ErrNotFound))
		return
	}
	sendGift, err := gift.GetService().SendGift(ctx, request, uInfo.UserId)
	if err != nil {
		log.Errorf("Send Gift failed, err: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			&SendResponse{
				Response: api.ErrorWithRequestId(log.ReqID(), err),
				Data:     sendGift,
			})
		return
	}

	rService := report.GetService()
	statsSingleLiveEntity := &model.StatsSingleLiveEntity{
		LiveId: request.LiveId,
		UserId: uInfo.UserId,
		BizId:  sendGift.AnchorId,
		Type:   model.StatsTypeGift,
		Count:  request.Amount,
	}
	rService.UpdateSingleLive(ctx, statsSingleLiveEntity)
	ctx.JSON(http.StatusOK, &SendResponse{
		Response: &api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: sendGift,
	})
}

type SendResponse struct {
	*api.Response
	Data *gift.SendGiftResponse `json:"data"`
}

// Test 用于测试的礼物支付
func (c *giftController) Test(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &gift.PayGiftRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	data := gift.GiftPayTestResp{Status: model.SendGiftStatusSuccess}
	ctx.JSON(http.StatusOK, &gift.PayGiftResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: data,
	})
}
