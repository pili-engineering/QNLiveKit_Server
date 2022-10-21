package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// httpHistogram prometheus 模型
	httpHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "qnlive_api",
		Subsystem:   "",
		Name:        "request",
		Help:        "Histogram of response latency (seconds) of http handlers.",
		ConstLabels: nil,
		Buckets:     nil,
	}, []string{"handler", "method", "code", "path"})
)

// init 初始化prometheus模型
func init() {
	prometheus.MustRegister(httpHistogram)
}

// Prometheus prometheus 监控
func Prometheus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		handler := c.HandlerName()
		path := c.Request.URL.Path
		httpHistogram.WithLabelValues(
			handler,
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
			path,
		).Observe(time.Since(start).Seconds())
	}
}
