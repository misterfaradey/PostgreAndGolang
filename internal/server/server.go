package server

import (
	"context"
	"github.com/misterfaradey/PostgreAndGolang/internal/server/controllers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerConf struct {
	GinMode        string
	Address        string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes uint
}

type Server interface {
	Engine() *gin.Engine
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	srv    *http.Server
	engine *gin.Engine
}

func NewServer(
	controller controllers.Controller,
	config *ServerConf,
) *server {

	s := &server{}
	s.setup(controller, config)

	return s
}

func (s *server) setup(controller controllers.Controller, config *ServerConf) {

	gin.SetMode(config.GinMode)

	s.engine = gin.New()
	s.engine.Use(gin.Recovery())

	if config.GinMode != gin.ReleaseMode {
		s.engine.Use(gin.Logger())
	}

	for _, action := range controller.Actions() {
		a := action

		s.engine.Handle(a.HttpMethod, a.RelativePath, controllers.MiddleWare, a.ActionExec)
	}

	s.srv = &http.Server{
		Addr:           config.Address,
		Handler:        s.engine,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: int(config.MaxHeaderBytes),
	}
}

func (s *server) Engine() *gin.Engine {
	return s.engine
}

func (s *server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
