package core

import (
	"actionflow/config"
	"actionflow/pkg/logutil"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"go.uber.org/zap"

	"gopkg.in/yaml.v3"
)

type Envelopment int32

const (
	DEVELOPMENT Envelopment = iota - 1
	PRODUCTION
)

const (
	DefaultPort   = "8000"
	DefaultListen = "0.0.0.0"
)

// Config struct defines the config structure
type Config struct {
	YC         *config.ConfigYAML
	configFile string
	cf         *configFlags

	// ZapLoggerBuilder is used to build the zap logger.
	ZapLoggerBuilder func(configYAML *Config) error

	// logger logs server-side operations. The default is nil,
	// and "setupLogging" must be called before starting server.
	// Do not set logger directly.
	loggerMu *sync.RWMutex
	logger   *zap.Logger
}

// configFlags has the set of flags used for command line parsing a Config
type configFlags struct {
	flagSet *flag.FlagSet
}

func NewConfig() *Config {
	cfg := &Config{
		cf:       &configFlags{},
		loggerMu: new(sync.RWMutex),
		logger:   nil,
		YC: &config.ConfigYAML{
			Logger:         "zap",
			LogOutput:      config.StdoutLogOutput,
			LogLevel:       logutil.DefaultLogLevel,
			LogEncoderType: "json",
		},
	}

	cfg.cf = &configFlags{
		flagSet: flag.NewFlagSet(cfg.YC.ServiceName, flag.ContinueOnError),
	}
	fs := cfg.cf.flagSet
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, usageline)
	}

	fs.StringVar(&cfg.configFile, "conf-file", "./config/development.yml", "Path to the server configuration file.")
	//fs.StringVar(&cfg.YC.Env, "env", cfg.YC.Env, "Service runtime environment")
	//fs.StringVar(&cfg.YC.ServiceName, "service-name", cfg.YC.ServiceName, "The service name")
	//fs.StringVar(&cfg.YC.Listen, "listen", cfg.YC.Listen, "Service binding address")
	//fs.StringVar(&cfg.YC.Port, "port", cfg.YC.Port, "The port number for the service binding")

	return cfg
}

func (cfg *Config) Parse(arguments []string) error {
	var err error
	err = cfg.cf.flagSet.Parse(arguments)
	switch err {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		os.Exit(2)
	}

	if len(cfg.cf.flagSet.Args()) != 0 {
		err = fmt.Errorf("'%s' is not a valid flag", cfg.cf.flagSet.Arg(0))
	}

	err = cfg.configFromFile(cfg.configFile)

	if len(cfg.YC.Listen) < 1 {
		cfg.YC.Listen = DefaultListen
	}

	if len(cfg.YC.Port) < 1 {
		cfg.YC.Port = DefaultPort
	}
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) configFromFile(path string) error {
	if err := configFromFile(cfg.YC, path); err != nil {
		return err
	}

	return cfg.validate()
}

func configFromFile(cfg *config.ConfigYAML, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) validate() error {
	if err := cfg.setupLogging(); err != nil {
		return err
	}
	return nil
}
