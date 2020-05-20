package rollingwriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	options := []Option{
		WithTimeTagFormat("200601021504"), WithLogPath("./"), WithFileName("foo"),
		WithAsynchronous(), WithBuffer(), WithBufferThreshold(8), WithCompress(), WithLock(),
		WithMaxRemain(3), WithRollingVolumeSize("100mb"), WithRollingTimePattern("0 0 * * *"),
	}
	cfg := NewDefaultConfig()
	for _, opt := range options {
		opt(&cfg)
	}
	destcfg := Config{
		LogPath:               "./",
		TimeTagFormat:         "200601021504",
		FileName:              "foo",
		MaxRemain:             3,
		RollingPolicy:         TimeRolling,
		RollingTimePattern:    "0 0 * * *",
		RollingVolumeSize:     "100mb",
		WriterMode:            "lock",
		BufferWriterThreshold: 8,
		Compress:              true,
	}
	assert.Equal(t, cfg, destcfg)
}
