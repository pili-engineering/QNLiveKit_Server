package admin

import (
	"github.com/dgrijalva/jwt-go"
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
	UserId    string `json:"user_id" form:"user_id"`
	PassWord  string `json:"pass_word" form:"pass_word"`
	ExpiresAt int64  `json:"expires_at"form:"expires_at"`
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
	_, err := manService.LoginManager(ctx, req.UserId, req.PassWord)
	if err != nil {
		log.Errorf("get user userId:%s, login error:%v", req.UserId, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorWithRequestId(log.ReqID(), err))
		return
	}

	authToken := token.AuthToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: req.ExpiresAt,
		},
		UserId: req.UserId,
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
