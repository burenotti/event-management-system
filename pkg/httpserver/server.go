package httpserver

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"time"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
)

type Server struct {
	server *fasthttp.Server
	notify chan error
	addr   string
}

func New(addr string, handler fasthttp.RequestHandler, logger *logrus.Logger) Server {
	s := Server{
		server: &fasthttp.Server{
			Handler:      handler,
			ReadTimeout:  _defaultReadTimeout,
			WriteTimeout: _defaultWriteTimeout,
			Logger:       logger,
		},
		addr:   addr,
		notify: make(chan error),
	}
	go s.run()
	return s
}

func (s *Server) run() {
	s.notify <- s.server.ListenAndServe(s.addr)
	close(s.notify)
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}
