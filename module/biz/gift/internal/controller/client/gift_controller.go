package client

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/base/stats"
	"github.com/qbox/livekit/module/biz/gift/dto"
	"github.com/qbox/livekit/module/biz/gift/internal/impl"
	"github.com/qbox/livekit/utils/logger"
)

type giftController struct {
}

var GiftController = &giftController{}

func RegisterRoutes() {
	//giftGroup := group.Group("/gift")
	//giftGroup.GET("/config/:type", GiftController.GetGiftConfig)
	//
	//giftGroup.GET("/list/live/:live_id", GiftController.ListGiftByLiveId)
	//giftGroup.GET("/list/anchor", GiftController.ListGiftByAnchorId)
	//giftGroup.GET("/list/user", GiftController.ListGiftByUserId)
	//giftGroup.POST("/send", GiftController.SendGift)

	httpq.ClientHandle(http.MethodGet, "/gift/config/:type", GiftController.GetGiftConfig)

	httpq.ClientHandle(http.MethodGet, "/gift/list/live/:live_id", GiftController.ListGiftByLiveId)
	httpq.ClientHandle(http.MethodGet, "/gift/list/anchor", GiftController.ListGiftByAnchorId)
	httpq.ClientHandle(http.MethodGet, "/gift/list/user", GiftController.ListGiftByUserId)
	httpq.ClientHandle(http.MethodPost, "/gift/send", GiftController.SendGift)
}

// GetGiftConfig 获取礼物配置
// return
func (*giftController) GetGiftConfig(context *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(context)
	typeId := context.Param("type")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("type is not int, err: %v", err)
		return nil, rest.ErrBadRequest.WithMessage("type is not int")
	}

	giftEntities, err := impl.GetInstance().GetListGiftEntity(context, typeIdInt)
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

// ListGiftByLiveId 根据直播间ID 查找礼物
// return rest.PageResult
func (c *giftController) ListGiftByLiveId(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("live_id")
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}

	gifts, count, err := impl.GetInstance().SearchGiftByLiveId(ctx, liveId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("Search Gift By LiveId failed, err: %v", err)
		return nil, rest.ErrInternal
	}

	endPage := false
	if len(gifts) < pageSizeInt {
		endPage = true
	}

	pageResult := rest.PageResult{
		TotalCount: count,
		PageTotal:  int(math.Ceil(float64(count) / float64(pageSizeInt))),
		EndPage:    endPage,
		List:       gifts,
	}
	return &pageResult, nil
}

// ListGiftByAnchorId 获取主播自己的礼物列表
// return rest.PageResult
func (c *giftController) ListGiftByAnchorId(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	userInfo := auth.GetUserInfo(ctx)
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}
	gifts, count, err := impl.GetInstance().SearchGiftByAnchorId(ctx, userInfo.UserId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("Search Gift by AnchorId  failed, err: %v", err)
		return nil, rest.ErrInternal
	}

	endPage := false
	if len(gifts) < pageSizeInt {
		endPage = true
	}

	pageResult := rest.PageResult{
		TotalCount: count,
		PageTotal:  int(math.Ceil(float64(count) / float64(pageSizeInt))),
		EndPage:    endPage,
		List:       gifts,
	}
	return &pageResult, nil
}

// ListGiftByUserId 查询用户送出的礼物列表
// return rest.PageResult
func (c *giftController) ListGiftByUserId(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	userInfo := auth.GetUserInfo(ctx)
	pageNum := ctx.DefaultQuery("page_num", "1")
	pageSize := ctx.DefaultQuery("page_size", "10")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Errorf("page_num is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Errorf("page_size is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}

	gifts, count, err := impl.GetInstance().SearchGiftByUserId(ctx, userInfo.UserId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Errorf("Search Gift ByUserId  failed, err: %v", err)
		return nil, rest.ErrInternal
	}

	endPage := false
	if len(gifts) < pageSizeInt {
		endPage = true
	}

	pageResult := rest.PageResult{
		TotalCount: count,
		PageTotal:  int(math.Ceil(float64(count) / float64(pageSizeInt))),
		EndPage:    endPage,
		List:       gifts,
	}

	return &pageResult, nil
}

// SendGift 发送礼物接口
func (c *giftController) SendGift(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &impl.SendGiftRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}

	uInfo := auth.GetUserInfo(ctx)
	sendGift, err := impl.GetInstance().SendGift(ctx, request, uInfo.UserId)
	if err != nil {
		log.Errorf("Send Gift failed, err: %v", err)
		return nil, err
	}

	rService := stats.GetService()
	statsSingleLiveEntity := &model.StatsSingleLiveEntity{
		LiveId: request.LiveId,
		UserId: uInfo.UserId,
		BizId:  sendGift.AnchorId,
		Type:   model.StatsTypeGift,
		Count:  request.Amount,
	}
	rService.UpdateSingleLive(ctx, statsSingleLiveEntity)
	return sendGift, nil
}

type SendResponse struct {
	*rest.Response
	Data *impl.SendGiftResponse `json:"data"`
}

//// Test 用于测试的礼物支付
//func (c *giftController) Test(ctx *gin.Context) {
//	log := logger.ReqLogger(ctx)
//	request := &impl.PayGiftRequest{}
//	if err := ctx.BindJSON(request); err != nil {
//		log.Errorf("bind request error %s", err.Error())
//		ctx.AbortWithStatusJSON(http.StatusOK, rest.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, &impl.PayGiftResponse{
//		Response: api.Response{
//			RequestId: log.ReqID(),
//			Code:      0,
//			Message:   "success",
//		},
//		Status: model.SendGiftStatusSuccess,
//	})
//}
