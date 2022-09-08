package apimonitor

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Middleware set gin middleware
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		host := c.Request.Host
		method := c.Request.Method
		status := c.Writer.Status()
		duration := time.Since(start).Milliseconds()

		monitor(method, host, path, status, int(duration))
	}
}
