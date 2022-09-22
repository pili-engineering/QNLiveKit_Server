package httpq

import (
	"fmt"
	"net"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/qbox/livekit/core/application"
)

var instance *server

type server struct {
	c      *Config
	engine *gin.Engine
}

func newServer(c *Config) *server {
	s := &server{
		c: c,
	}

	s.engine = gin.New()
	s.engine.Use(Cors())
	return s
}

func (s *server) Start() error {
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

func (s *server) Stop(err error) {
}

func Cors() gin.HandlerFunc {
	c := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:    []string{"Content-Type", "Access-Token", "Authorization"},
		MaxAge:          6 * time.Hour,
	}

	return cors.New(c)
}
