package core

import (
	"actionflow/config"
	"actionflow/pkg/logutil"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

// GetLogger returns the logger.
func (cfg Config) GetLogger() *zap.Logger {
	cfg.loggerMu.RLock()
	l := cfg.logger
	cfg.loggerMu.RUnlock()
	return l
}

func (cfg *Config) setupLogging() error {
	switch cfg.YC.Logger {
	case "zap":
		var writeType logutil.WriteType

		switch cfg.YC.LogOutput {
		case config.FileLogOutput:
			writeType = logutil.File
		case config.StdoutAndFileLogOutput:
			writeType = logutil.StdoutAndFile
		case config.StdoutLogOutput:
			fallthrough
		default:
			writeType = logutil.Stdout
		}

		if cfg.ZapLoggerBuilder == nil {
			cfg.ZapLoggerBuilder = func(c *Config) error {
				lc := logutil.LoggerConfig{
					Service:     c.YC.ServiceName,
					Loglevel:    cfg.YC.LogLevel,
					WriteType:   writeType,
					EncoderType: logutil.Json,
				}
				if len(cfg.YC.LogFileConfigJSON) > 1 {
					logfileConfig, err := setupLoggerWriteFileConfig(cfg.YC.LogFileConfigJSON)
					if err != nil {
						return err
					}

					lc.FileConfig = logfileConfig
				}
				c.logger = logutil.NewLogger(lc)
				return nil
			}
		}

		err := cfg.ZapLoggerBuilder(cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown logger option %q", cfg.YC.Logger)
	}

	return nil
}

type logRotationConfig struct {
	*lumberjack.Logger
}

// Sync implements zap.Sink
func (logRotationConfig) Sync() error { return nil }

// setupLoggerWriteFileConfig initializes log file
func setupLoggerWriteFileConfig(logRotateConfigJSON string) (logutil.LoggerWriteFileConfig, error) {
	var loggerWriteFileConfig logutil.LoggerWriteFileConfig

	if err := json.Unmarshal([]byte(logRotateConfigJSON), &loggerWriteFileConfig); err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		var syntaxError *json.SyntaxError
		switch {
		case errors.As(err, &syntaxError):
			return loggerWriteFileConfig, fmt.Errorf("improperly formatted loggerWriteFileConfig: %w", err)
		case errors.As(err, &unmarshalTypeError):
			return loggerWriteFileConfig, fmt.Errorf("invalid loggerWriteFileConfig: %w", err)
		}
	}

	return loggerWriteFileConfig, nil
}
