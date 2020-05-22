package logs

import (
	"testing"
	"time"

	"github.com/Muskchen/toolkits/rollingwriter"
)

func TestInit(t *testing.T) {
	cfg := &Config{
		Format: "2006-01-02 15:04:05",
		Type:   "json",
		Appenders: []appender{{
			Level: "info",
			Rolling: &rollingwriter.Config{
				TimeTagFormat:      "20060102150405",
				LogPath:            "./",
				FileName:           "info",
				MaxRemain:          3,
				RollingPolicy:      rollingwriter.WithoutRolling,
				RollingTimePattern: "* * * * *",
				WriterMode:         "lock",
				Compress:           true,
			},
		}, {
			Level: "error",
			Rolling: &rollingwriter.Config{
				TimeTagFormat:      "20060102150405",
				LogPath:            "./",
				FileName:           "error",
				MaxRemain:          30,
				RollingPolicy:      rollingwriter.WithoutRolling,
				RollingTimePattern: "* * * * *",
				WriterMode:         "lock",
				Compress:           false,
			},
		}},
	}
	Init(cfg)
	for i := 0; i < 1000; i++ {
		SDebugf("测试debug%d", i)
		SInfof("测试info%d", i)
		SWarnf("测试warn%d", i)
		SErrorf("测试错误%d", i)
		//Panic("测试panic",
		//	zap.Int("panic", i),
		//)
		//Fatal("测试fatal",
		//	zap.Int("fatal", i),
		//)
		time.Sleep(time.Duration(1) * time.Second)
	}
}
