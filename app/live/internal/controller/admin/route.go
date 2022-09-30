package admin

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoute(engine *gin.Engine) {
	group := engine.Group("/manager")
	group.POST("/login", censorController.LoginManager)

	group.POST("/censor/callback", censorController.CallbackCensorJob)
	//group.POST("/gift/test", client.GiftController.Test)
	//adminGroup := engine.Group("/admin", liveauth.AuthAdminHandleFunc(config.AppConfig.JwtKey), middleware.OperatorLogMiddleware())
	//adminGroup := engine.Group("/admin", middleware2.AuthAdminHandleFunc(config.AppConfig.JwtKey), middleware.OperatorLogMiddleware())
	//RegisterCensorRoutes(adminGroup)
	//RegisterGiftRoute(adminGroup)
}
