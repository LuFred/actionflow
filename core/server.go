package core

import (
	"actionflow/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"go.uber.org/zap"
)

type Server struct {
	Cfg config.ConfigYAML

	lgMu *sync.RWMutex
	lg   *zap.Logger

	readMu sync.RWMutex

	// stop signals the run goroutine should shutdown.
	stop chan struct{}
	// done is closed when all goroutines from start() complete.
	done chan struct{}

	errorc chan error

	wg sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	router *Router
}

func NewServer(cfg Config) (srv *Server, err error) {
	srv = &Server{
		Cfg:    *cfg.YC,
		lg:     cfg.GetLogger(),
		lgMu:   new(sync.RWMutex),
		errorc: make(chan error, 1),
	}
	return
}

func (s *Server) Logger() *zap.Logger {
	s.lgMu.RLock()
	l := s.lg
	s.lgMu.RUnlock()
	return l
}

func (s *Server) Start(r *Router) {
	s.router = r
	s.start()
	os.Exit(0)
}

func (s *Server) start() {
	s.done = make(chan struct{})
	s.stop = make(chan struct{})
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.run()

}

func (s *Server) run() {
	lg := s.Logger()
	lg.Info(fmt.Sprintf("listen: %s:%s", s.Cfg.Listen, s.Cfg.Port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", s.Cfg.Listen, s.Cfg.Port), s.router.chi); err != nil {
		lg.Panic(err.Error())
	}
}
