package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/controller/server"
	"github.com/qbox/livekit/app/live/internal/dto"
	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/biz/live"
	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/biz/notify"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/biz/user"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/common/auth/liveauth"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
)

func RegisterCensorRoutes(group *gin.RouterGroup) {
	censorGroup := group.Group("/censor")
	censorGroup.POST("/config", censorController.UpdateCensorConfig)
	censorGroup.GET("/config", censorController.GetCensorConfig)
	censorGroup.POST("/stoplive/:liveId", censorController.PostStopLive)
}

var censorController = &CensorController{}

type CensorController struct {
}

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (c *CensorController) LoginManager(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &LoginRequest{}
	if err := ctx.BindQuery(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	manService := admin.GetManagerService()
	admin, err := manService.FindAdminByUserName(ctx, req.UserName)
	if err != nil {
		log.Errorf("userName:%s, login error:%v", req.UserName, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	} else if admin.Password != req.Password {
		log.Errorf("userName:%s, login error:%v", req.UserName, api.ErrorLoginWrong)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), api.ErrorLoginWrong))
		return
	}

	authToken := token.AuthToken{
		UserId: admin.UserId,
		Role:   "admin",
	}

	tokenService := token.GetService()
	if token, err := tokenService.GenAuthToken(&authToken); err != nil {
		log.Errorf("")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	} else {
		resp := &server.GetAuthTokenResponse{
			Response: api.Response{
				RequestId: log.ReqID(),
				Code:      0,
				Message:   "success",
			},
		}
		resp.Data.AccessToken = token
		resp.Data.ExpiresAt = authToken.ExpiresAt
		ctx.JSON(http.StatusOK, resp)
	}
}

type CensorConfigResponse struct {
	api.Response
	Data *dto.CensorConfigDto `json:"data"`
}

func (c *CensorController) UpdateCensorConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	req := &dto.CensorConfigDto{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}
	censorService := admin.GetCensorService()
	err := censorService.UpdateCensorConfig(ctx, dto.CConfigDtoToEntity(req))
	if err != nil {
		log.Errorf(" UpdateCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &CensorConfigResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: req,
	})
}

func (c *CensorController) GetCensorConfig(ctx *gin.Context) {
	log := logger.ReqLogger(ctx)
	censorService := admin.GetCensorService()
	censorConfig, err := censorService.GetCensorConfig(ctx)
	if err != nil {
		log.Errorf("GetCensorConfig error:%v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}
	ctx.JSON(http.StatusOK, &CensorConfigResponse{
		Response: api.Response{
			RequestId: log.ReqID(),
			Code:      0,
			Message:   "success",
		},
		Data: dto.CConfigEntityToDto(censorConfig),
	})
}

func (c *CensorController) PostStopLive(ctx *gin.Context) {
	adminInfo := liveauth.GetAdminInfo(ctx)

	log := logger.ReqLogger(ctx)
	liveId := ctx.Param("liveId")
	if liveId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorWithRequestId(log.ReqID(), api.ErrInvalidArgument))
		return
	}

	liveEntity, err := live.GetService().LiveInfo(ctx, liveId)
	if err != nil {
		log.Errorf("find live error %s", err.Error())
		ctx.AbortWithStatusJSON(http.StatusOK, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	anchorInfo, err := user.GetService().FindUser(ctx, liveEntity.AnchorId)
	if err != nil {
		log.Errorf("get anchor info for %s error %s", liveEntity.AnchorId, err.Error())
	}
	notifyItem := LiveNotifyItem{
		LiveId:  liveEntity.LiveId,
		Message: "直播涉嫌违规，\n管理员已关闭直播间。",
	}
	err = notify.SendNotifyToLive(ctx, anchorInfo, liveEntity, notify.ActionTypeCensorStop, &notifyItem)
	if err != nil {
		log.Errorf("send notify to live %s error %s", liveEntity.LiveId, err.Error())
	}

	err = live.GetService().AdminStopLive(ctx, liveId, model.LiveStopReasonCensor, adminInfo.UserId)
	if err != nil {
		log.Errorf("stop live failed, err: %v", err)
		ctx.JSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	ctx.JSON(http.StatusOK, api.Response{
		Code:      200,
		Message:   "success",
		RequestId: log.ReqID(),
	})
}

type LiveNotifyItem struct {
	LiveId  string `json:"live_id"`
	Message string `json:"message"`
}
