package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/biz/model"
	"github.com/qbox/livekit/module/base/auth"
	"github.com/qbox/livekit/module/store/mysql"
	"github.com/qbox/livekit/utils/logger"
	"github.com/qbox/livekit/utils/timestamp"
)

// OperatorLogMiddleware 用于记录管理员后台的操作日志
func OperatorLogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.ReqLogger(ctx)
		ol := model.OperationLog{}

		//Args      string              `json:"args" `      // 操作内容详情
		ol.IP = ctx.ClientIP()
		ol.Method = ctx.Request.Method
		ol.URL = ctx.Request.URL.RequestURI()
		ol.CreatedAt = timestamp.Now()

		if ol.Method != http.MethodGet {
			buf, _ := ioutil.ReadAll(ctx.Request.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

			ctx.Request.Body = rdr2
			ol.Args = readBody(rdr1)
		}
		ctx.Next()

		adminInfo := auth.GetAdminInfo(ctx)
		if adminInfo == nil {
			log.Errorf("request [%s, %s] with no admin info", ol.Method, ctx.Request.URL.Path)
			return
		}
		ol.UserId = adminInfo.UserId
		ol.StatusCode = ctx.Writer.Status()

		db := mysql.GetLive(log.ReqID())
		db.Save(&ol)
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
