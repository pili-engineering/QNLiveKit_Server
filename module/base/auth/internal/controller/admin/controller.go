package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/admin"
	"github.com/qbox/livekit/module/base/auth/internal/controller/server"
	"github.com/qbox/livekit/module/base/auth/internal/impl"
	token2 "github.com/qbox/livekit/module/base/auth/internal/token"
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

// LoginManager 管理员用户名，密码登录
// return
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
		return nil, rest.ErrForbidden.WithMessage("Invalid username or password")
	}

	authToken := token2.AuthToken{
		UserId: admin.UserId,
		Role:   "admin",
	}

	if token, err := impl.GetInstance().GenAuthToken(&authToken); err != nil {
		log.Errorf("gen token error %v", err)
		return nil, rest.ErrInternal
	} else {
		result := &server.GetAuthTokenResult{
			AccessToken: token,
			ExpiresAt:   authToken.ExpiresAt,
		}

		return result, nil
	}
}
