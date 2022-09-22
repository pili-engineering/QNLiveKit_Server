package httpq

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(ctx *gin.Context) (interface{}, error)

func makeHandle(f HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ret, err := f(ctx)
		if err != nil {
			return
		}

		ctx.JSON(http.StatusOK, ret)
	}
}
