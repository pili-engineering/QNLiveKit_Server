package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/controller/server"
	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
)

func RegisterManagerRoutes(group *gin.RouterGroup) {
}

var ManagerController = &managerController{}

type managerController struct {
}

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (c *managerController) LoginManager(ctx *gin.Context) {
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
