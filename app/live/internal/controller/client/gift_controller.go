package client

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/biz/gift"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
	"math"
	"net/http"
	"strconv"
)

type giftController struct {
}

var GiftController = &giftController{}

func RegisterGiftRoutes(group *gin.RouterGroup) {
	giftGroup := group.Group("/gift")
	giftGroup.GET("/list/live/:live_id", GiftController.ListGiftByLiveId)
	giftGroup.GET("/list/anchor", GiftController.ListGiftByLiveId)
	giftGroup.GET("/list/user", GiftController.ListGiftByLiveId)

}

func (c *giftController) ListGiftByLiveId(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
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
