package cron

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Muskchen/toolkits/str"
)

type Locker interface {
	Lock(key string) (bool, error)
	Unlock(key string) error
}

type timeUnit int

const MAXJOBNUM = 10000

const (
	seconds timeUnit = iota + 1
	minutes
	hours
	days
	weeks
)

var (
	loc    = time.Local
	locker Locker
)

func ChangeLoc(newLocation *time.Location) {
	loc = newLocation
	defaultScheduler.ChangeLoc(newLocation)
}

func SetLocker(l Locker) {
	locker = l
}

// 利用反射执行函数，获得具体的结果
func callJobFuncWithParams(jobFunc interface{}, params []interface{}) ([]reflect.Value, error) {
	// 获取具体执行的函数
	f := reflect.ValueOf(jobFunc)
	// 判断函数参数的个数
	if len(params) != f.Type().NumIn() {
		return nil, ErrParamsNotAdapted
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	// 执行函数，相当于执行了f(in...)
	return f.Call(in), nil
}

// 获取函数名称
func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// 函数名称进行hash运算
func getFunctionKey(funcName string) string {
	h, _ := str.Sha256(funcName)
	return h
}

func Jobs() []*Job {
	return defaultScheduler.Jobs()
}

// 时间格式化“hh:mm:ss”,拆分时间表达式
func formatTime(t string) (hour, min, sec int, err error) {
	ts := strings.Split(t, ":")
	if len(ts) != 3 {
		return 0, 0, 0, ErrTimeFormat
	}

	if hour, err = strconv.Atoi(ts[0]); err != nil {
		return 0, 0, 0, err
	}
	if min, err = strconv.Atoi(ts[1]); err != nil {
		return 0, 0, 0, err
	}
	if sec, err = strconv.Atoi(ts[2]); err != nil {
		return 0, 0, 0, err
	}
	if hour < 0 || hour > 23 || min < 0 || min > 59 || sec < 0 || sec > 59 {
		return 0, 0, 0, ErrTimeFormat
	}
	return hour, min, sec, nil
}

// 下次打点时间，一秒之后
func NextTick() *time.Time {
	now := time.Now().Add(time.Second)
	return &now
}
