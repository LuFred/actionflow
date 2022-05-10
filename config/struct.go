package config

const (
	StdoutLogOutput        = "stdout"
	FileLogOutput          = "file"
	StdoutAndFileLogOutput = "stdoutandfile"
)

type ConfigYAML struct {
	Env                 string     `yaml:"envelopment"`
	ServiceName         string     `yaml:"serviceName"`
	TracingTransportURL string     `yaml:"tracingurl"`
	MySqlDb             *MySqlConf `yaml:"mysql"`
	Listen              string     `yaml:"listen"`
	Port                string     `yaml:"port"`

	TemporalConf *TemporalConf `yaml:"temporal"`

	// Logger is logger options: currently only supports "zap".
	Logger string `yaml:"logger"`

	// LogLevel configures log level. Only supports debug, info, warn, error, panic, or fatal. Default 'info'.
	LogLevel string `yaml:"log-level"`

	// LogOutput configures log output. stdout, file, stdoutandfile
	LogOutput string `yaml:"log-outputs"`

	LogFileConfigJSON string `yaml:"log-file-config-json"`

	LogEncoderType string
}

type MySqlConf struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	PWD      string `yaml:"pwd"`
	DB       string `yaml:"db"`
}

type TemporalConf struct {
	HostPost               string `yaml:"hostPort"`
	JobRunInstanceQueue    string `yaml:"jobRunInstanceQueue"`
	JobRunInstanceWorkflow string `yaml:"jobRunInstanceWorkflow"`
}
