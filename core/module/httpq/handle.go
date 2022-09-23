package httpq

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/rest"
	"github.com/qbox/livekit/utils/logger"
)

type HandlerFunc func(ctx *gin.Context) (interface{}, error)

func makeHandle(f HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)

		var err error = nil
		var ret interface{} = nil

		defer func() {
			if err != nil {
				restErr, ok := err.(*rest.Error)
				if ok {
					if restErr.StatusCode == 0 {
						restErr = restErr.WithStatusCode(http.StatusInternalServerError)
					}
				} else {
					restErr = rest.ErrInternal.WithMessage(err.Error())
				}
				ctx.JSON(restErr.StatusCode, restErr.WithRequestId(log.ReqID()))
			} else {
				resp := &rest.Response{
					RequestId: log.ReqID(),
					Code:      0,
					Message:   "",
					Data:      ret,
				}
				ctx.JSON(http.StatusOK, resp)
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				const size = 16 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Error("panic: ", r, fmt.Sprintf("\n%s", buf))

				if _, ok := r.(error); ok {
					err = r.(error)
				} else {
					err = fmt.Errorf("%v", r)
				}
			}
		}()
		ret, err = f(ctx)
	}
}
