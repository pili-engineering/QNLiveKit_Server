package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/app/live/internal/middleware"
	middleware2 "github.com/qbox/livekit/module/base/auth/internal/middleware"
)

func RegisterRoute(engine *gin.Engine) {
	group := engine.Group("/manager")
	group.GET("/login", censorController.LoginManager)

	group.POST("/censor/callback", censorController.CallbackCensorJob)
	adminGroup := engine.Group("/admin", middleware2.AuthAdminHandleFunc(config.AppConfig.JwtKey), middleware.OperatorLogMiddleware())
	RegisterCensorRoutes(adminGroup)
	RegisterGiftRoute(adminGroup)
}
