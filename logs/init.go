package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Muskchen/toolkits/rollingwriter"
	"github.com/Muskchen/toolkits/runner"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	// 时间格式
	Format string `json:"format" yaml:"format"`
	// 日志格式，json和console
	Type string `json:"type" yaml:"type"`
	// 是否开通栈追踪，开启后error及以下级别打印栈信息
	Stacktrace  bool `json:"stacktrace" yaml:"stacktrace"`
	Development bool `json:"development" yaml:"development"`
	// 日志文件及级别配置
	Appenders []appender `json:"appenders" yaml:"appenders"`
}

type appender struct {
	// 日志级别
	Level string `json:"level" yaml:"level"`
	// writer信息
	Rolling *rollingwriter.Config `json:"rolling" yaml:"rolling"`
}

var (
	logger *zap.Logger
	config zapcore.EncoderConfig
	encode zapcore.Encoder
)

func init() {
	runner.Init()
	fmt.Printf("HostName: %s, Workerspace: %s\n", runner.Hostname, runner.Cwd)
	config = newEncoderConfig("2006-01-02 15:04:05")
	encode = encoder("json", config)
	writer := os.Stdout
	level := logLevel("debug")
	core := zapcore.NewCore(encode, zapcore.AddSync(writer), level)
	logger = zap.New(core, zap.AddCaller())
	logger = logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	logger.WithOptions(zap.Development())
}

func Init(cfg *Config) {
	config = newEncoderConfig(cfg.Format)
	encode = encoder(cfg.Type, config)
	var Logs []zapcore.Core
	for _, app := range cfg.Appenders {
		writer, err := rollingwriter.NewWriterFromConfig(app.Rolling)
		if err != nil {
			writer = os.Stdout
		}
		level := logLevel(app.Level)
		core := zapcore.NewCore(encode, zapcore.AddSync(writer), level)
		Logs = append(Logs, core)
	}

	core := zapcore.NewTee(Logs...)
	logger = zap.New(core, zap.AddCaller())
	if cfg.Stacktrace {
		logger = logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	}
	if cfg.Development {
		logger.WithOptions(zap.Development())
	}
}

func Logger() *zap.Logger {
	return logger
}

func SLogger() *zap.SugaredLogger {
	return logger.Sugar()
}

func Close() {
	if err := logger.Sync(); err != nil {
		logger.Error("closed err", zap.Error(err))
	}
}

// 初始化配置
func newEncoderConfig(format string) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:    "msg",                       // 日志消息对应的key
		LevelKey:      "level",                     // 日志级别对应的key
		TimeKey:       "ts",                        // 时间对应的key
		CallerKey:     "file",                      // 调用信息对应的key
		StacktraceKey: "stacktrace",                // 栈追踪对应的key
		EncodeLevel:   zapcore.CapitalLevelEncoder, // 大写的日志级别显示
		LineEnding:    zapcore.DefaultLineEnding,   // 日志的换行符，默认为"\n"
		EncodeTime: func(t time.Time, en zapcore.PrimitiveArrayEncoder) {
			en.AppendString(t.Format(format))
		}, // 时间格式化
		EncodeDuration: zapcore.SecondsDurationEncoder, // 序列化时间的类型
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 采用短格式输出字段
	}
}

// 日志输出格式
func encoder(typ string, config zapcore.EncoderConfig) (encoder zapcore.Encoder) {
	typ = strings.TrimSpace(strings.ToLower(typ))
	switch typ {
	case "json":
		return zapcore.NewJSONEncoder(config)
	case "console":
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewJSONEncoder(config)
	}
}

// 日志级别
func logLevel(level string) zapcore.Level {
	level = strings.TrimSpace(strings.ToLower(level))
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
