package main

import (
	"actionflow/core"
	"actionflow/db"
	"actionflow/middleware"
	"actionflow/pkg/logutil"
	"actionflow/pkg/temporalutil"
	"actionflow/routers"
	"os"

	"go.uber.org/zap"
)

func main() {
	cfg := core.NewConfig()
	err := cfg.Parse(os.Args[1:])
	lg := cfg.GetLogger()
	if lg == nil {
		// use this logger
		lg = logutil.NewLogger(logutil.LoggerConfig{
			Loglevel:    "info",
			WriteType:   logutil.Stdout,
			EncoderType: logutil.Json,
		})
	}

	defer func() {
		logger := cfg.GetLogger()
		if logger != nil {
			logger.Sync()
		}
	}()

	if err != nil {
		lg.Panic(err.Error())
	}

	srv, err := core.NewServer(*cfg)
	if err != nil {
		lg.Error("Server initialization failed", zap.Error(err))
		lg.Panic(err.Error())
	}
	router := core.NewRouter(srv)
	router.UseMiddlewares(
		middleware.OauthInfoMiddleware,
		middleware.CORSMiddleware)
	router.UseHandlers(
		routers.RegisterHelloworldHandlers(),
		routers.RegisterFlowHandlers())

	db.InitDb(srv)

	err = temporalutil.InitClient(srv.Cfg.TemporalConf.HostPost)
	if err != nil {
		lg.Error("temporal client initialization failed", zap.Error(err))
		lg.Panic(err.Error())
	}

	srv.Start(router)
}
