package client

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/app/live/internal/report"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
)

func RegisterStatsRoutes(group *gin.RouterGroup) {
	statsGroup := group.Group("/stats")

	statsGroup.POST("/singleLive", StatsController.PostStatsSingleLive)
	statsGroup.GET("/singleLive/:live_id", StatsController.GetStatsSingleLive)
}

type statsController struct {
}

var StatsController = &statsController{}

func (c *statsController) PostStatsSingleLive(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	request := &SingleLiveRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	if len(request.Data) == 0 {
		log.Errorf("Data is empty, invalid request %+v", request)
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	entities := make([]*model.StatsSingleLiveEntity, 0, len(request.Data))
	for _, d := range request.Data {
		entities = append(entities, dto.StatsSLDtoToEntity(d))
	}

	rService := report.GetService()
	err := rService.PostStatsSingleLive(ctx, entities)
	if err != nil {
		log.Errorf("post statistic single live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type SingleLiveRequest struct {
	Data []*dto.SingleLiveInfo `json:"data"`
}

func (c *statsController) GetStatsSingleLive(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("live_id")
	rService := report.GetService()
	message, err := rService.GetStatsSingleLive(ctx, liveId)
	if err != nil {
		log.Errorf("get statistic single live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	response := &StatsSingleLiveResponse{
		Response: api.SuccessResponse(log.ReqID()),
		Data:     *message,
	}
	ctx.JSON(http.StatusOK, response)
}

type StatsSingleLiveResponse struct {
	api.Response
	Data report.CommonStats `json:"data"`
}
