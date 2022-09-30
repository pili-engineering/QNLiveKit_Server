package impl

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/module/base/admin"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/utils/logger"
)

func (s *ServiceImpl) RegisterAdminAuth() {
	httpq.SetAdminAuth(func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)

		authInfo := ctx.GetHeader("Authorization")
		if len(authInfo) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rest.ErrUnauthorized.WithRequestId(log.ReqID()))
			return
		}

		authToken, err := s.tokenService.ParseAuthToken(authInfo)
		if err != nil {
			switch err1 := err.(type) {
			case *rest.Error:
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, err1.WithRequestId(log.ReqID()))
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, rest.ErrInternal.WithRequestId(log.ReqID()))
			}
			return
		}

		if authToken.Role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, rest.ErrForbidden.WithRequestId(log.ReqID()))
			return
		}

		manService := admin.GetManagerService()
		_, err = manService.FindAdminByUserId(ctx, authToken.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, rest.ErrUnauthorized.WithRequestId(log.ReqID()))
			return
		}

		uInfo := auth.AdminInfo{
			UserId: authToken.UserId,
			Role:   authToken.Role,
		}
		ctx.Set(auth.AdminCtxKey, &uInfo)

		ctx.Next()
	})
}
