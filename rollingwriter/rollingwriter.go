package rollingwriter

import (
	"errors"
	"io"
	"os"
	"path"
)

// 三种滚动模式
const (
	WithoutRolling = iota
	TimeRolling
	VolumeRolling
)

// 一些默认的全局变量
var (
	BufferSize      = 0x100000
	QueueSize       = 1024
	Precision       = 1
	DefualtFileMode = os.FileMode(0644)
	DefualtFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND

	// 自定义错误
	ErrInternal        = errors.New("error internal")
	ErrClosed          = errors.New("error write on close")
	ErrInvalidArgument = errors.New("error argument invalid")
)

type Manager interface {
	Fire() chan string
	Close()
}

type RollingWriter interface {
	io.Writer
	Close() error
}

type Config struct {
	TimeTagFormat string `json:"time_tag_format" yaml:"timeTagFormat"` //时间格式化结构
	LogPath       string `json:"log_path" yaml:"logPath"`              // 日志路径
	FileName      string `json:"file_name" yaml:"fileName"`            // 日志文件名称
	MaxRemain     int    `json:"max_remain" yaml:"maxRemain"`          // 日志文件的最大存留数

	// 日志滚动策略，三个选项
	// 0：WithoutRolling:，不滚动
	// 1：TimeRolling，时间滚动策略，
	// 2：VolumeRolling，大小滚动策略
	RollingPolicy      int    `json:"rolling_policy" yaml:"rollingPolicy"`
	RollingTimePattern string `json:"rolling_time_pattern" yaml:"rollingTimePattern"` // 时间滚动策略时的cron表达式
	RollingVolumeSize  string `json:"rolling_volume_size" yaml:"rollingVolumeSize"`   // 大小滚动策略时的截断大小

	WriterMode            string `json:"writer_mode" yaml:"writerMode"`                 // none, lock, async, buffer
	BufferWriterThreshold int    `json:"buffer_threshold" yaml:"bufferWriterThreshold"` // 一部并发是缓存池的大小
	Compress              bool   `json:"compress" yaml:"compress"`                      // 是否压缩历史日志
}

// 默认配置
func NewDefaultConfig() Config {
	return Config{
		TimeTagFormat:         "200601021504",
		LogPath:               "./log",
		FileName:              "log",
		MaxRemain:             -1,
		RollingPolicy:         1,
		RollingTimePattern:    "0 0 * * *",
		RollingVolumeSize:     "1G",
		WriterMode:            "lock",
		BufferWriterThreshold: 64,
		Compress:              false,
	}
}

// 生成日志文件完整路径
func LogFilePath(c *Config) (filepath string) {
	filepath = path.Join(c.LogPath, c.FileName) + ".log"
	return filepath
}

// 配置构造函数，用于更新配置
type Option func(*Config)

// 更新时间格式
func WithTimeTagFormat(format string) Option {
	return func(c *Config) {
		c.TimeTagFormat = format
	}
}

// 更新日志文件目录
func WithLogPath(path string) Option {
	return func(c *Config) {
		c.LogPath = path
	}
}

// 更新日志文件名称
func WithFileName(name string) Option {
	return func(c *Config) {
		c.FileName = name
	}
}

// 改为async模式
func WithAsynchronous() Option {
	return func(c *Config) {
		c.WriterMode = "async"
	}
}

// 改为lock模式
func WithLock() Option {
	return func(c *Config) {
		c.WriterMode = "lock"
	}
}

// 改为buffer模式
func WithBuffer() Option {
	return func(c *Config) {
		c.WriterMode = "buffer"
	}
}

// 修改异步并发时的缓存池大小
func WithBufferThreshold(n int) Option {
	return func(c *Config) {
		c.BufferWriterThreshold = n
	}
}

// 开启压缩历史日志文件
func WithCompress() Option {
	return func(c *Config) {
		c.Compress = true
	}
}

// 更新历史文件保存数
func WithMaxRemain(max int) Option {
	return func(c *Config) {
		c.MaxRemain = max
	}
}

// 设置为不滚动模式
func WithoutRollingPolicy() Option {
	return func(c *Config) {
		c.RollingPolicy = WithoutRolling
	}
}

// 设置为按时间滚动模式，更新滚动时间表达式
func WithRollingTimePattern(pattern string) Option {
	return func(c *Config) {
		c.RollingPolicy = TimeRolling
		c.RollingTimePattern = pattern
	}
}

// 设置为按大小滚动模式，更新滚动时截断的最大值
func WithRollingVolumeSize(size string) Option {
	return func(c *Config) {
		c.RollingPolicy = VolumeRolling
		c.RollingVolumeSize = size
	}
}
