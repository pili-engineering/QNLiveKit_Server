package client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/stats/dto"
	"github.com/qbox/livekit/module/base/stats/service"
	"github.com/qbox/livekit/utils/logger"
)

func RegisterRoutes() {
	//statsGroup := group.Group("/stats")
	//statsGroup.POST("/singleLive", StatsController.PostStatsSingleLive)
	//statsGroup.GET("/singleLive/:live_id", StatsController.GetStatsSingleLive)

	httpq.ClientHandle(http.MethodPost, "/stats/singleLive", StatsController.PostStatsSingleLive)
	httpq.ClientHandle(http.MethodGet, "/stats/singleLive/:live_id", StatsController.GetStatsSingleLive)
}

type statsController struct {
}

var StatsController = &statsController{}

// PostStatsSingleLive 提交一个统计
// return nil
func (c *statsController) PostStatsSingleLive(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	request := &SingleLiveRequest{}
	if err := ctx.BindJSON(request); err != nil {
		log.Errorf("bind request error %s", err.Error())
		return nil, rest.ErrBadRequest.WithMessage(err.Error())
	}
	if len(request.Data) == 0 {
		log.Errorf("Data is empty, invalid request %+v", request)
		return nil, rest.ErrBadRequest.WithMessage("Empty data")
	}

	entities := make([]*model.StatsSingleLiveEntity, 0, len(request.Data))
	for _, d := range request.Data {
		entities = append(entities, dto.StatsSLDtoToEntity(d))
	}

	err := service.Instance.PostStatsSingleLive(ctx, entities)
	if err != nil {
		log.Errorf("post statistic single live error %s", err.Error())
		return nil, err
	}
	return nil, nil
}

type SingleLiveRequest struct {
	Data []*dto.SingleLiveInfo `json:"data"`
}

// GetStatsSingleLive 查询直播间的统计信息
// return *impl.CommonStats
func (c *statsController) GetStatsSingleLive(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("live_id")
	message, err := service.Instance.GetStatsSingleLive(ctx, liveId)
	if err != nil {
		log.Errorf("get statistic single live error %s", err.Error())
		return nil, err
	}

	return message, nil
}
