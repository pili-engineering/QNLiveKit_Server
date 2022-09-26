package httpq

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq/middleware"
	"github.com/qbox/livekit/core/module/httpq/monitor"
)

func (s *Server) createEngin() {
	s.engine = gin.New()
	s.engine.Use(middleware.Cors(), middleware.Logger(), middleware.Prometheus(), monitor.Middleware())

	s.clientGroup = s.engine.Group("/client")
	s.serverGroup = s.engine.Group("/server")
	s.adminGroup = s.engine.Group("/admin")
	s.callbackGroup = s.engine.Group("/callback")
}

// Handle 为了保证相同路径前缀的表现一致，如果是预定
func (s *Server) Handle(httpMethod, relativePath string, handler HandlerFunc) {
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}

	route, path := s.selectGroup(relativePath)
	route.Handle(httpMethod, path, makeHandle(handler))
}

func (s *Server) selectGroup(relativePath string) (gin.IRoutes, string) {
	switch {
	case strings.HasPrefix(relativePath, "/client"):
		return s.clientGroup, strings.TrimPrefix(relativePath, "/client")
	case strings.HasPrefix(relativePath, "/server"):
		return s.serverGroup, strings.TrimPrefix(relativePath, "/server")
	case strings.HasPrefix(relativePath, "/admin"):
		return s.adminGroup, strings.TrimPrefix(relativePath, "/admin")
	case strings.HasPrefix(relativePath, "/callback"):
		return s.callbackGroup, strings.TrimPrefix(relativePath, "/callback")
	default:
		return s.engine, relativePath
	}
}

func (s *Server) ClientHandle(httpMethod, relativePath string, handler HandlerFunc) {
	s.clientGroup.Handle(httpMethod, relativePath, makeHandle(handler))
}

func (s *Server) ServerHandle(httpMethod, relativePath string, handler HandlerFunc) {
	s.serverGroup.Handle(httpMethod, relativePath, makeHandle(handler))
}

func (s *Server) AdminHandle(httpMethod, relativePath string, handler HandlerFunc) {
	s.adminGroup.Handle(httpMethod, relativePath, makeHandle(handler))
}

func (s *Server) CallbackHandle(httpMethod, relativePath string, handler HandlerFunc) {
	s.callbackGroup.Handle(httpMethod, relativePath, makeHandle(handler))
}
