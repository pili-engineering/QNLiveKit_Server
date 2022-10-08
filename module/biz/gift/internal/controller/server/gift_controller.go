package server

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/biz/gift/dto"
	"github.com/qbox/livekit/module/biz/gift/internal/impl"
	"github.com/qbox/livekit/utils/logger"
)

type giftController struct {
}

var GiftController = &giftController{}

func RegisterRoutes() {
	//giftGroup := group.Group("/gift")
	//
	//giftGroup.GET("/config/:type", GiftController.GetGiftConfig)
	//giftGroup.POST("/config", GiftController.AddGiftConfig)
	//giftGroup.DELETE("/config/:gift_id", GiftController.DeleteGiftConfig)
	//
	//giftGroup.GET("/list/live/:live_id", GiftController.ListGiftByLiveId)
	//giftGroup.GET("/list/anchor/:anchor_id", GiftController.ListGiftByAnchorId)
	//giftGroup.GET("/list/user/:user_id", GiftController.ListGiftByUserId)

	httpq.ServerHandle(http.MethodGet, "/gift/config/:type", GiftController.GetGiftConfig)
	httpq.ServerHandle(http.MethodPost, "/gift/config", GiftController.AddGiftConfig)
	httpq.ServerHandle(http.MethodDelete, "/gift/config/:gift_id", GiftController.DeleteGiftConfig)

	httpq.ServerHandle(http.MethodGet, "/gift/list/live/:live_id", GiftController.ListGiftByLiveId)
	httpq.ServerHandle(http.MethodGet, "/gift/list/anchor/:anchor_id", GiftController.ListGiftByAnchorId)
	httpq.ServerHandle(http.MethodGet, "/gift/list/user/:user_id", GiftController.ListGiftByUserId)
}

// AddGiftConfig 添加一个礼物配置
func (*giftController) AddGiftConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &dto.GiftConfigDto{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if req.GiftId == 0 || req.Name == "" {
		return nil, rest.ErrBadRequest.WithMessage("Name 不能为空 且 GiftId 需要 >0 ")
	}

	err := impl.GetInstance().SaveGiftEntity(ctx, dto.GiftDtoToEntity(req))
	if err != nil {
		log.Errorf("add gift config  failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}

// DeleteGiftConfig 删除一个礼物配置
func (*giftController) DeleteGiftConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	typeId := ctx.Param("gift_id")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("gift_id is not int, err: %v", err)
		return nil, rest.ErrBadRequest
	}

	err = impl.GetInstance().DeleteGiftEntity(ctx, typeIdInt)
	if err != nil {
		log.Errorf("delete gift config  failed, err: %v", err)
		return nil, err
	}
	return nil, nil
}

// GetGiftConfig 查询礼物配置
func (*giftController) GetGiftConfig(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	typeId := ctx.Param("type")
	typeIdInt, err := strconv.Atoi(typeId)
	if err != nil {
		log.Errorf("type is not int, err: %v", err)
		return nil, rest.ErrBadRequest
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

// ListGiftByLiveId 根据直播间ID 查询礼物信息
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
		return nil, err
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

// ListGiftByAnchorId 查询主播收到的礼物列表
func (c *giftController) ListGiftByAnchorId(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	anchorId := ctx.Param("anchor_id")
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
	gifts, count, err := impl.GetInstance().SearchGiftByAnchorId(ctx, anchorId, pageNumInt, pageSizeInt)
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

// ListGiftByUserId 查看用户送出去的礼物列表
func (c *giftController) ListGiftByUserId(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	userId := ctx.Param("user_id")
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

	gifts, count, err := impl.GetInstance().SearchGiftByUserId(ctx, userId, pageNumInt, pageSizeInt)
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
