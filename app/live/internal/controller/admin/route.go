package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/app/live/internal/config"
	"github.com/qbox/livekit/common/auth/liveauth"
)

func RegisterRoute(engine *gin.Engine) {
	group := engine.Group("/manager")
	group.GET("/login", censorController.LoginManager)
	group.POST("/censor/callback", censorController.CallbackCensorJob)
	adminGroup := engine.Group("/admin", liveauth.AuthAdminHandleFunc(config.AppConfig.JwtKey))
	RegisterCensorRoutes(adminGroup)
}
