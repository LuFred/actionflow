package logutil

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type WriteType int8

type EncoderType int8

const (
	Stdout WriteType = iota - 1
	File
	StdoutAndFile
)

const (
	Console EncoderType = iota - 1
	Json
)

type LoggerConfig struct {
	LogPath     string
	Service     string
	Loglevel    string
	WriteType   WriteType
	EncoderType EncoderType
	FileConfig  LoggerWriteFileConfig
}

type LoggerWriteFileConfig struct {
	Filename string `json:"filename" yaml:"filename"`

	MaxSize int `json:"maxsize" yaml:"maxsize"`

	MaxAge int `json:"maxage" yaml:"maxage"`

	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	LocalTime bool `json:"localtime" yaml:"localtime"`

	Compress bool `json:"compress" yaml:"compress"`
}

// logpath 日志文件路径
// loglevel 日志级别
func NewLogger(cfg LoggerConfig) *zap.Logger {
	// 日志分割
	hook := lumberjack.Logger{
		Filename:   cfg.FileConfig.Filename,   // 日志文件路径，默认 os.TempDir()
		MaxSize:    cfg.FileConfig.MaxAge,     // 每个日志文件保存10M，默认 100M
		MaxBackups: cfg.FileConfig.MaxBackups, // 保留30个备份，默认不限
		MaxAge:     cfg.FileConfig.MaxAge,     // 保留7天，默认不限
		Compress:   cfg.FileConfig.Compress,   // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)

	level := ConvertToZapLevel(cfg.Loglevel)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "lineNum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	var encoder zapcore.Encoder
	switch cfg.EncoderType {
	case Console:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case Json:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var core zapcore.Core
	switch cfg.WriteType {
	case Stdout:
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
	case File:
		core = zapcore.NewCore(
			encoder,
			write,
			level,
		)
	case StdoutAndFile:
		core = zapcore.NewCore(
			encoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(write)),
			level,
		)
	default:
		core = zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
	}

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段,如：添加一个服务器名称
	filed := zap.Fields(zap.String("service", cfg.Service))
	// 构造日志
	return zap.New(core, caller, development, filed)
}
