package holmes

import (
	"os"
	"path"
	"path/filepath"
	"time"
)

type options struct {
	UseCGroup bool

	DumpPath        string
	DumpProfileType dumpProfileType
	DumpFullStack   bool

	LogLevel int
	Logger   *os.File

	CollectInterval time.Duration

	CoolDown time.Duration

	GrOpts     *grOptions
	MemOpts    *memOptions
	CPUOpts    *cpuOptions
	ThreadOpts *threadOptions
}

type Option interface {
	apply(*options) error
}

type optionFunc func(*options) (err error)

func (f optionFunc) apply(opts *options) error {
	return f(opts)
}

func newOptions() *options {
	return &options{
		GrOpts:          newGrOptions(),
		MemOpts:         newMemOptions(),
		CPUOpts:         newCPUOptions(),
		ThreadOpts:      newThreadOptions(),
		LogLevel:        LogLevelDebug,
		Logger:          os.Stdout,
		CollectInterval: defaultInterval,
		CoolDown:        defaultCooldown,
		DumpPath:        defaultDumpPath,
		DumpProfileType: defaultDumpProfileType,
		DumpFullStack:   false,
	}
}

func WithCollectInterval(interval string) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.CollectInterval, err = time.ParseDuration(interval)
		return
	})
}

func WithCoolDown(coolDown string) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.CoolDown, err = time.ParseDuration(coolDown)
		return
	})
}

func WithDumpPath(dumpPath string, loginfo ...string) Option {
	return optionFunc(func(opts *options) (err error) {
		f := path.Join(dumpPath, defaultLoggerName)
		if len(loginfo) > 0 {
			f = dumpPath + "/" + path.Join(loginfo...)
		}
		opts.DumpPath = filepath.Dir(f)
		opts.Logger, err = os.OpenFile(f, defaultLoggerFlags, defaultLoggerPerm)
		if err != nil && os.IsNotExist(err) {
			if err = os.MkdirAll(opts.DumpPath, 0755); err != nil {
				return
			}
			opts.Logger, err = os.OpenFile(f, defaultLoggerFlags, defaultLoggerPerm)
			if err != nil {
				return
			}
		}
		return
	})
}

func WithBinaryDump() Option {
	return withDumpProfile(binaryDump)
}

func WithTextDump() Option {
	return withDumpProfile(textDump)
}

func WithFullStack(isFull bool) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.DumpFullStack = isFull
		return
	})
}

func withDumpProfile(profileType dumpProfileType) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.DumpProfileType = profileType
		return
	})
}

type grOptions struct {
	Enable                      bool
	GoroutineTriggerNumMin      int
	GoroutineTriggerPercentDiff int
	GoroutineTriggerNumAbs      int
}

func newGrOptions() *grOptions {
	return &grOptions{
		Enable:                      false,
		GoroutineTriggerNumMin:      defaultGoroutineTriggerMin,
		GoroutineTriggerPercentDiff: defaultGoroutineTriggerDiff,
		GoroutineTriggerNumAbs:      defaultGoroutineTriggerAbs,
	}
}

func WithGoroutineDump(min, diff, abs int) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.GrOpts.GoroutineTriggerNumMin = min
		opts.GrOpts.GoroutineTriggerNumAbs = abs
		opts.GrOpts.GoroutineTriggerPercentDiff = diff
		return
	})
}

type memOptions struct {
	Enable                bool
	MemTriggerPercentMin  int
	MemTriggerPercentDiff int
	MenTriggerPercentAbs  int
}

func newMemOptions() *memOptions {
	return &memOptions{
		Enable:                false,
		MemTriggerPercentMin:  defaultMemTriggerMin,
		MemTriggerPercentDiff: defaultMemTriggerDiff,
		MenTriggerPercentAbs:  defaultMemTriggerAbs,
	}
}

func WithMemDump(min, diff, abs int) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.MemOpts.MemTriggerPercentMin = min
		opts.MemOpts.MemTriggerPercentDiff = diff
		opts.MemOpts.MenTriggerPercentAbs = abs
		return
	})
}

type threadOptions struct {
	Enable                   bool
	ThreadTriggerPercentMin  int
	ThreadTriggerPercentDiff int
	ThreadTriggerPercentAbs  int
}

func newThreadOptions() *threadOptions {
	return &threadOptions{
		Enable:                   false,
		ThreadTriggerPercentMin:  defaultThreadTriggerMin,
		ThreadTriggerPercentDiff: defaultThreadTriggerDiff,
		ThreadTriggerPercentAbs:  defaultThreadTriggerAbs,
	}
}

func WithThreadDump(min, diff, abs int) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.ThreadOpts.ThreadTriggerPercentMin = min
		opts.ThreadOpts.ThreadTriggerPercentDiff = diff
		opts.ThreadOpts.ThreadTriggerPercentAbs = abs
		return
	})
}

type cpuOptions struct {
	Enable                bool
	CPUTriggerPercentMin  int
	CPUTriggerPercentDiff int
	CPUTriggerPercentAbs  int
}

func newCPUOptions() *cpuOptions {
	return &cpuOptions{
		Enable:                false,
		CPUTriggerPercentMin:  defaultCPUTriggerMin,
		CPUTriggerPercentDiff: defaultCPUTriggerDiff,
		CPUTriggerPercentAbs:  defaultCPUTriggerAbs,
	}
}

func WithCPUDump(min, diff, abs int) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.CPUOpts.CPUTriggerPercentAbs = abs
		opts.CPUOpts.CPUTriggerPercentDiff = diff
		opts.CPUOpts.CPUTriggerPercentMin = min
		return
	})
}

func WithCGroup(useCGroup bool) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.UseCGroup = useCGroup
		return
	})
}

func WithLoggerLevel(level int) Option {
	return optionFunc(func(opts *options) (err error) {
		opts.LogLevel = level
		return
	})
}
