package httpq

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/application"
)

var instance *Server

type Server struct {
	c *Config

	engine        *gin.Engine
	clientGroup   *gin.RouterGroup
	serverGroup   *gin.RouterGroup
	adminGroup    *gin.RouterGroup
	callbackGroup *gin.RouterGroup

	clientAuthHandle gin.HandlerFunc
	serverAuthHandle gin.HandlerFunc
	adminAuthHandle  gin.HandlerFunc
}

func newServer(c *Config) *Server {
	s := &Server{
		c: c,
	}
	s.createEngine()

	return s
}

func (s *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", s.c.Addr)
	if err != nil {
		return fmt.Errorf("resolve addr %s error %v", s.c.Addr, err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s error %v", s.c.Addr, err)
	}

	go func() {
		err1 := s.engine.RunListener(listener)
		application.Stop(err1)
	}()
	return nil
}

func (s *Server) Stop(err error) {
}
