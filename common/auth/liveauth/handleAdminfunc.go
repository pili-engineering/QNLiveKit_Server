package liveauth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/qbox/livekit/biz/admin"
	"github.com/qbox/livekit/biz/token"
	"github.com/qbox/livekit/common/api"
	"github.com/qbox/livekit/utils/logger"
	"net/http"
)

const AdminCtxKey = "AdminCtxKey"

type AdminInfo struct {
	UserId string
	Role   string
}

func GetAdminInfo(ctx context.Context) *AdminInfo {
	log := logger.ReqLogger(ctx)
	i := ctx.Value(AdminCtxKey)
	if i == nil {
		return nil
	}

	if t, ok := i.(*AdminInfo); ok {
		return t
	} else {
		log.Errorf("%+v not user info", i)
		return nil
	}
}

func AuthAdminHandleFunc(jwtKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)

		auth := ctx.GetHeader("Authorization")
		if len(auth) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), api.ErrBadToken))
			return
		}

		tokenService := token.GetService()
		authToken, err := tokenService.ParseAuthToken(auth)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}
		manService := admin.GetManagerService()
		_, err = manService.FindAdminByUserId(ctx, authToken.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}
		if authToken.Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, api.ErrorWithRequestId(log.ReqID(), err))
			return
		}

		uInfo := AdminInfo{
			UserId: authToken.UserId,
			Role:   authToken.Role,
		}
		ctx.Set(AdminCtxKey, &uInfo)

		ctx.Next()
	}
}
