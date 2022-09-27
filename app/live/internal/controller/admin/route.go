package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/app/live/internal/controller/client"
	"github.com/qbox/livekit/app/live/internal/middleware"
	"github.com/qbox/livekit/common/auth/liveauth"
)

func RegisterRoute(engine *gin.Engine) {
	group := engine.Group("/manager")
	group.GET("/login", censorController.LoginManager)

	group.POST("/censor/callback", censorController.CallbackCensorJob)
	group.POST("/gift/test", client.GiftController.Test)
	adminGroup := engine.Group("/admin", liveauth.AuthAdminHandleFunc(config.AppConfig.JwtKey), middleware.OperatorLogMiddleware())
	RegisterCensorRoutes(adminGroup)
	RegisterGiftRoute(adminGroup)
}
