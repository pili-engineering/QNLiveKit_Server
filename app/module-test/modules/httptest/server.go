package httptest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/module/httpq"
	"github.com/qbox/livekit/core/rest"
)

var instance *Server

type Server struct {
}

func (s *Server) RegisterApi() {
	s.registerClientApis()
	s.registerAdminApis()
	s.registerServerApis()
	s.registerCallbackApis()
}

func (s *Server) registerAdminApis() {
	httpq.AdminHandle(http.MethodGet, "/httptest/hello", func(ctx *gin.Context) (interface{}, error) {
		return gin.H{"message": "hello, admin"}, nil
	})
}

func (s *Server) registerClientApis() {
	httpq.ClientHandle(http.MethodGet, "/httptest/hello", func(ctx *gin.Context) (interface{}, error) {
		return gin.H{"message": "hello, client"}, nil
	})
}

func (s *Server) registerServerApis() {
	httpq.ServerHandle(http.MethodGet, "/httptest/hello", func(ctx *gin.Context) (interface{}, error) {
		return gin.H{"message": "hello, server"}, nil
	})
}

func (s *Server) registerCallbackApis() {
	httpq.CallbackHandle(http.MethodGet, "/httptest/hello", func(ctx *gin.Context) (interface{}, error) {
		return gin.H{"message": "hello, callback"}, nil
	})

	httpq.CallbackHandle(http.MethodGet, "/httptest/error", func(ctx *gin.Context) (interface{}, error) {
		return nil, rest.ErrBadRequest.WithMessage("Just a error")
	})

	httpq.CallbackHandle(http.MethodGet, "/httptest/error/raw", func(ctx *gin.Context) (interface{}, error) {
		return nil, fmt.Errorf("raw error")
	})

	httpq.CallbackHandle(http.MethodGet, "/httptest/panic", func(ctx *gin.Context) (interface{}, error) {
		panic(rest.ErrInternal.WithMessage("what a panic"))
	})

	httpq.CallbackHandle(http.MethodGet, "/httptest/panic/raw", func(ctx *gin.Context) (interface{}, error) {
		panic("o my god")
	})
}
