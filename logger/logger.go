package logger

import (
	"fmt"
	"os"
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
	Stdout   bool   `yaml:"stdout"`
}

var log *logrus.Logger

func Init(cfg LoggerConf) {
	var err error
	logFile := filepath.Join(cfg.Dir, cfg.Name)
	logger := lumberjack.Logger{
		Filename:  logFile,
		MaxAge:    cfg.MaxAge,
		MaxSize:   cfg.MaxSize,
		LocalTime: true,
	}
	fmt.Println(logger.Filename)
	log = logrus.New()
	if log.Level, err = logrus.ParseLevel(cfg.Level); err != nil {
		fmt.Printf("log level %s error, %s", cfg.Level, err)
		return
	}
	if cfg.Stdout {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(&logger)
	}
	log.Formatter = &logrus.TextFormatter{}
}

func Trace(args ...interface{}) {
	log.Trace(args)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args)
}

func Traceln(args ...interface{}) {
	log.Traceln(args)
}

func Debug(args ...interface{}) {
	log.Debug(args)
}

func Degbugf(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Debugln(args ...interface{}) {
	log.Debugln(args)
}

func Info(args ...interface{}) {
	log.Info(args)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Infoln(args ...interface{}) {
	log.Infoln(args)
}

func Warn(args ...interface{}) {
	log.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args)
}

func Warnln(args ...interface{}) {
	log.Warnln(args)
}

func Error(args ...interface{}) {
	log.Error(args)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
}

func Errorln(args ...interface{}) {
	log.Errorln(args)
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
}

func Fatalln(args ...interface{}) {
	log.Fatalln(args)
}

func Panic(args ...interface{}) {
	log.Panic(args)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args)
}

func Panicln(args ...interface{}) {
	log.Panicln(args)
}
