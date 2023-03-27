package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qbox/livekit/module/biz/relay"
	"github.com/qbox/livekit/module/store/cache"
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
	httpq.Handle(http.MethodPost, "/manager/gift/test", GiftController.Test)
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
	// 记录送礼的积分
	sid, cacheErr := cache.Client.Get(fmt.Sprintf(model.PkIntegral, sendGift.AnchorId))
	if cacheErr == nil && sid != "" {
		// pk记录
		KEY := fmt.Sprintf(model.PkIntegral, sid)
		cache.Client.HIncrBy(KEY, sendGift.AnchorId, int64(sendGift.Amount))
		// 更新拓展字段，添加积分信息
		result, cacheErr := cache.Client.HGetAll(KEY)
		if cacheErr == nil {
			data, _ := json.Marshal(getPkIntegral(ctx, sid, result))
			extends := model.Extends{"pkIntegral": string(data)}
			relay.GetRelayService().UpdateRelayExtends(ctx, sid, extends)
		} else {
			log.Errorf("failed to query redis，error【%v】", cacheErr.Error())
		}
	}
	if err != nil {
		return nil, err
	}
	return sendGift, nil
}

type SendResponse struct {
	*rest.Response
	Data *impl.SendGiftResponse `json:"data"`
}

// PkIntegral 记录pk过程中积分信息
type PkIntegral struct {
	RecvUserId string
	RecvRoomId string
	RecvScore  int
	InitUserId string
	InitRoomId string
	InitScore  int
}

func getPkIntegral(ctx context.Context, sid string, scoreMap map[string]string) *PkIntegral {
	// 查询pk会话信息
	session, err := relay.GetRelayService().GetRelaySession(ctx, sid)
	if err != nil || session == nil {
		return nil
	}
	recvScore := 0
	initScore := 0
	if stringRecvScore, ok := scoreMap[session.RecvUserId]; ok && stringRecvScore != "" {
		recvScore, err = strconv.Atoi(stringRecvScore)
	}
	if stringInitScore, ok := scoreMap[session.InitUserId]; ok && stringInitScore != "" {
		recvScore, err = strconv.Atoi(stringInitScore)
	}
	return &PkIntegral{
		RecvUserId: session.RecvUserId,
		RecvRoomId: session.RecvRoomId,
		RecvScore:  recvScore,
		InitUserId: session.InitUserId,
		InitRoomId: session.InitRoomId,
		InitScore:  initScore,
	}
}

// Test 用于测试的礼物支付
func (c *giftController) Test(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &impl.PayGiftRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	res := &GiftPayTestResp{
		Status: model.SendGiftStatusSuccess,
	}
	return res, nil
}

type GiftPayTestResp struct {
	Status int `json:"status"`
}
