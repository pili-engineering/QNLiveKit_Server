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

	statsGroup.POST("/singleLive", StatsController.StatsSingleLive)
}

type statsController struct {
}

var StatsController = &statsController{}

func (c *statsController) StatsSingleLive(ctx *gin.Context) {
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
	err := rService.StatsSingleLive(ctx, entities)
	if err != nil {
		log.Errorf("statistic single live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse(log.ReqID()))
}

type SingleLiveRequest struct {
	Data []*dto.SingleLiveInfo `json:"data"`
}
