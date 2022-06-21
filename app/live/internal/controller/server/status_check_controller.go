package server

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
)

var StatusCheckController = &statusCheckController{}

type statusCheckController struct {
}

func (*statusCheckController) CheckStatus(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)

	resp := &api.Response{
		RequestId: log.ReqID(),
		Code:      0,
		Message:   "success",
	}
	ctx.JSON(http.StatusOK, resp)

}
