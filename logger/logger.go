package logger

import (
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConf struct {
	Name     string `yaml:"name"`
	Dir      string `yaml:"dir"`
	Level    string `yaml:"level"`
	MaxSize  int    `yaml:"maxsize"`
	MaxAge   int    `yaml:"maxage"`
	Compress bool   `yaml:"compress"`
}

var Log *logrus.Logger

func Init(cfg LoggerConf) {
	var err error
	logFile := filepath.Join(cfg.Dir, cfg.Name)
	logger := lumberjack.Logger{
		Filename:  logFile,
		MaxAge:    cfg.MaxAge,
		MaxSize:   cfg.MaxSize,
		LocalTime: true,
	}
	Log = logrus.New()
	if Log.Level, err = logrus.ParseLevel(cfg.Level); err != nil {
		fmt.Printf("log level %s error, %s", cfg.Level, err)
		return
	}
	Log.Out = &logger
	Log.Formatter = &logrus.TextFormatter{}
}
