package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/admin"
	"github.com/qbox/livekit/module/base/auth/internal/controller/server"
	"github.com/qbox/livekit/utils/logger"
)

//group := engine.Group("/manager")
//group.GET("/login", censorController.LoginManager)

func RegisterRoutes() {
	httpq.Handle(http.MethodPost, "/manager/login", LoginManager)
}

type LoginRequest struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func LoginManager(ctx *gin.Context) (interface{}, error) {
	log := logger.ReqLogger(ctx)
	req := &LoginRequest{}
	if err := ctx.BindJSON(req); err != nil {
		log.Errorf("bind request error %v", err)
		return nil, rest.ErrBadRequest
	}
	manService := admin.GetManagerService()
	admin, err := manService.FindAdminByUserName(ctx, req.UserName)
	if err != nil {
		log.Errorf("userName:%s, login error:%v", req.UserName, err)
		return nil, err
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
