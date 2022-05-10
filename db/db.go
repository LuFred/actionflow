package db

import (
	"actionflow/core"
	"actionflow/db/entity"
	"actionflow/pkg/ormutil"
	"actionflow/pkg/osutil"
	"go.uber.org/zap"
)

func InitDb(s *core.Server) {
	goorm.SetLogging(true)
	goorm.RegisterModel(&entity.FlowEntity{})
	goorm.RegisterModel(&entity.FlowActionEdgeEntity{})
	goorm.RegisterModel(&entity.ActionFlowJobEntity{})
	goorm.RegisterModel(&entity.TimerEntity{})
	goorm.RegisterModel(&entity.RetryOptionsEntity{})
	goorm.RegisterModel(&entity.ActionEntity{})
	goorm.RegisterModel(&entity.ActionFlowJobRunInstanceEntity{})

	err := goorm.RegisterDataBase("default",
		"mysql",
		s.Cfg.MySqlDb.Host,
		s.Cfg.MySqlDb.DB,
		s.Cfg.MySqlDb.UserName,
		s.Cfg.MySqlDb.PWD,
		10,
		10,
	)

	lg := s.Logger()

	if err != nil {
		lg.Error("Exception initializing database connection", zap.String("err", err.Error()))
		osutil.Exit(1)
	}

	lg.Info("database connection successful")
}
