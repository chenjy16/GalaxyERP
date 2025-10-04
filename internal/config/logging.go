package config

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggingConfig 日志配置
type LoggingConfig struct {
	// 日志级别 (debug, info, warn, error, fatal)
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// 日志格式 (json, console)
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	// 输出目标 (stdout, stderr, file)
	Output string `mapstructure:"output" json:"output" yaml:"output"`
	// 日志文件路径（当output为file时使用）
	FilePath string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
	// 是否启用调用者信息
	EnableCaller bool `mapstructure:"enable_caller" json:"enable_caller" yaml:"enable_caller"`
	// 是否启用堆栈跟踪
	EnableStacktrace bool `mapstructure:"enable_stacktrace" json:"enable_stacktrace" yaml:"enable_stacktrace"`
	// 时间格式
	TimeFormat string `mapstructure:"time_format" json:"time_format" yaml:"time_format"`
	// 日志轮转配置
	Rotation LogRotationConfig `mapstructure:"rotation" json:"rotation" yaml:"rotation"`
}

// LogRotationConfig 日志轮转配置
type LogRotationConfig struct {
	// 是否启用日志轮转
	Enabled bool `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	// 最大文件大小（MB）
	MaxSize int `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	// 最大保留天数
	MaxAge int `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	// 最大备份文件数
	MaxBackups int `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	// 是否压缩备份文件
	Compress bool `mapstructure:"compress" json:"compress" yaml:"compress"`
}

// DefaultLoggingConfig 默认日志配置
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:            "info",
		Format:           "json",
		Output:           "stdout",
		FilePath:         "logs/app.log",
		EnableCaller:     true,
		EnableStacktrace: false,
		TimeFormat:       time.RFC3339,
		Rotation: LogRotationConfig{
			Enabled:    true,
			MaxSize:    100, // 100MB
			MaxAge:     30,  // 30天
			MaxBackups: 10,  // 10个备份文件
			Compress:   true,
		},
	}
}

// GetLogLevel 获取日志级别
func (c LoggingConfig) GetLogLevel() zapcore.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// GetEncoder 获取编码器
func (c LoggingConfig) GetEncoder() zapcore.Encoder {
	encoderConfig := c.getEncoderConfig()
	
	switch strings.ToLower(c.Format) {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}

// getEncoderConfig 获取编码器配置
func (c LoggingConfig) getEncoderConfig() zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     c.getTimeEncoder(),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 如果是控制台格式，使用彩色编码
	if strings.ToLower(c.Format) == "console" {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return config
}

// getTimeEncoder 获取时间编码器
func (c LoggingConfig) getTimeEncoder() zapcore.TimeEncoder {
	timeFormat := c.TimeFormat
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}
}

// GetWriteSyncer 获取写入同步器
func (c LoggingConfig) GetWriteSyncer() zapcore.WriteSyncer {
	switch strings.ToLower(c.Output) {
	case "stdout":
		return zapcore.AddSync(os.Stdout)
	case "stderr":
		return zapcore.AddSync(os.Stderr)
	case "file":
		return c.getFileWriteSyncer()
	default:
		return zapcore.AddSync(os.Stdout)
	}
}

// getFileWriteSyncer 获取文件写入同步器
func (c LoggingConfig) getFileWriteSyncer() zapcore.WriteSyncer {
	// 这里可以集成 lumberjack 或其他日志轮转库
	// 目前先使用简单的文件输出
	file, err := os.OpenFile(c.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// 如果文件打开失败，回退到标准输出
		return zapcore.AddSync(os.Stdout)
	}
	return zapcore.AddSync(file)
}

// BuildLogger 构建日志器
func (c LoggingConfig) BuildLogger() *zap.Logger {
	core := zapcore.NewCore(
		c.GetEncoder(),
		c.GetWriteSyncer(),
		c.GetLogLevel(),
	)

	options := []zap.Option{}
	
	if c.EnableCaller {
		options = append(options, zap.AddCaller())
	}
	
	if c.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zap.New(core, options...)
}

// LoggerOptions 日志器选项
type LoggerOptions struct {
	// 服务名称
	ServiceName string
	// 服务版本
	ServiceVersion string
	// 环境
	Environment string
	// 实例ID
	InstanceID string
}

// BuildLoggerWithOptions 使用选项构建日志器
func (c LoggingConfig) BuildLoggerWithOptions(opts LoggerOptions) *zap.Logger {
	logger := c.BuildLogger()
	
	// 添加全局字段
	fields := []zap.Field{}
	
	if opts.ServiceName != "" {
		fields = append(fields, zap.String("service", opts.ServiceName))
	}
	
	if opts.ServiceVersion != "" {
		fields = append(fields, zap.String("version", opts.ServiceVersion))
	}
	
	if opts.Environment != "" {
		fields = append(fields, zap.String("environment", opts.Environment))
	}
	
	if opts.InstanceID != "" {
		fields = append(fields, zap.String("instance_id", opts.InstanceID))
	}
	
	if len(fields) > 0 {
		logger = logger.With(fields...)
	}
	
	return logger
}

// InitializeGlobalLogger 初始化全局日志器
func InitializeGlobalLogger(config LoggingConfig, opts LoggerOptions) {
	logger := config.BuildLoggerWithOptions(opts)
	zap.ReplaceGlobals(logger)
}

// GetDevelopmentConfig 获取开发环境配置
func GetDevelopmentConfig() LoggingConfig {
	config := DefaultLoggingConfig()
	config.Level = "debug"
	config.Format = "console"
	config.EnableCaller = true
	config.EnableStacktrace = true
	return config
}

// GetProductionConfig 获取生产环境配置
func GetProductionConfig() LoggingConfig {
	config := DefaultLoggingConfig()
	config.Level = "info"
	config.Format = "json"
	config.Output = "file"
	config.EnableCaller = false
	config.EnableStacktrace = false
	return config
}

// GetTestingConfig 获取测试环境配置
func GetTestingConfig() LoggingConfig {
	config := DefaultLoggingConfig()
	config.Level = "warn"
	config.Format = "console"
	config.Output = "stderr"
	config.EnableCaller = false
	config.EnableStacktrace = false
	return config
}