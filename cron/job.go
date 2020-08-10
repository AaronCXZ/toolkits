package cron

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"
)

var (
	ErrTimeFormat           = errors.New("time format error")
	ErrParamsNotAdapted     = errors.New("the number of params is not adapted")
	ErrNotAFunction         = errors.New("only functions can be schedule into the job queue")
	ErrPeriodNotSpecified   = errors.New("unspecified job period")
	ErrParameterCannotBeNil = errors.New("nil paramaters cannot be used with reflection")
)

type Job struct {
	interval uint64                   // 两次任务时间间隔
	jobFunc  string                   // 执行函数的名称
	unit     timeUnit                 // 间隔时间的单位
	atTime   time.Duration            // 任务的运行时间
	err      error                    // 错误
	loc      *time.Location           // 时区
	lastRun  time.Time                // 上次运行的时间
	nextRun  time.Time                // 下次运行的时间
	startDay time.Weekday             // 从周几开始
	funcs    map[string]interface{}   // 执行的函数
	fparams  map[string][]interface{} // 执行函数的参数
	lock     bool                     // 任务的锁
	tags     []string                 // 任务的标签
}

// 创建新的任务
func NewJob(interval uint64) *Job {
	return &Job{
		interval: interval,
		loc:      loc,
		lastRun:  time.Unix(0, 0),
		nextRun:  time.Unix(0, 0),
		startDay: time.Sunday,
		funcs:    make(map[string]interface{}),
		fparams:  make(map[string][]interface{}),
		tags:     []string{},
	}
}

// 当前时间是否大于等于下次运行时间
func (j *Job) shouldRun() bool {
	return time.Now().Unix() >= j.nextRun.Unix()
}

// 运行任务
func (j *Job) run() ([]reflect.Value, error) {
	// 当有锁时
	if j.lock {
		if locker == nil {
			return nil, fmt.Errorf("trying to lock %s with nil locker", j.jobFunc)
		}
		// 函数名称hash之后作为key
		key := getFunctionKey(j.jobFunc)

		locker.Lock(key)
		defer locker.Unlock(key)
	}
	result, err := callJobFuncWithParams(j.funcs[j.jobFunc], j.fparams[j.jobFunc])
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 错误信息
func (j *Job) Err() error {
	return j.err
}

// 添加定时任务的函数以及对应的参数
func (j *Job) Do(jobFun interface{}, params ...interface{}) error {
	if j.err != nil {
		return j.err
	}
	// 判断是否为函数
	typ := reflect.ValueOf(jobFun)
	if typ.Kind() != reflect.Func {
		return ErrNotAFunction
	}
	// 获取函数名称
	fname := getFunctionName(jobFun)
	j.funcs[fname] = jobFun
	j.fparams[fname] = params
	j.jobFunc = fname
	// 判断任务的下次执行时间
	now := time.Now().In(j.loc)
	if !j.nextRun.After(now) {
		j.scheduleNextRun()
	}
	return nil
}

// panic处理
func (j *Job) DoSafely(jobFunc interface{}, params ...interface{}) error {
	recoveryWrapperFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Internal panic occurred: %s", r)
			}
		}()

		callJobFuncWithParams(jobFunc, params)
	}
	return j.Do(recoveryWrapperFunc)
}

// 任务运行的具体时间
// 时分秒
func (j *Job) At(t string) *Job {
	hour, min, sec, err := formatTime(t)
	if err != nil {
		j.err = ErrTimeFormat
		return j
	}

	j.atTime = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute + time.Duration(sec)*time.Second
	return j
}

// 格式化任务的具体运行时间
func (j *Job) GetAt() string {
	return fmt.Sprintf("%d:%d", j.atTime/time.Hour, (j.atTime%time.Hour)/time.Minute)
}

// 更新时区
func (j *Job) Loc(loc *time.Location) *Job {
	j.loc = loc
	return j
}

// 添加标签，一个或多个
func (j *Job) Tag(t string, others ...string) {
	j.tags = append(j.tags, t)
	j.tags = append(j.tags, others...)
}

// 删除一个标签
func (j *Job) Untag(t string) {
	var newTags []string
	for _, tag := range j.tags {
		if t != tag {
			newTags = append(newTags, tag)
		}
	}
	j.tags = newTags
}

// 获取标签
func (j *Job) Tags() []string {
	return j.tags
}

// 任务间隔时间
func (j *Job) periodDuration() (time.Duration, error) {
	interval := time.Duration(j.interval)
	var periodDuration time.Duration

	switch j.unit {
	case seconds:
		periodDuration = interval * time.Second
	case minutes:
		periodDuration = interval * time.Minute
	case hours:
		periodDuration = interval * time.Hour
	case days:
		periodDuration = interval * time.Hour * 24
	case weeks:
		periodDuration = interval * time.Hour * 24 * 7
	default:
		return 0, ErrPeriodNotSpecified
	}
	return periodDuration, nil
}

// 新的一天开始,时分秒信息归零
func (j *Job) roundToMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, j.loc)
}

// 更新任务下次的运行时间
func (j *Job) scheduleNextRun() error {
	now := time.Now()
	// 判断是否运行过任务
	if j.lastRun == time.Unix(0, 0) {
		j.lastRun = now
	}
	// 任务间隔时间
	periodDuration, err := j.periodDuration()
	if err != nil {
		return err
	}

	switch j.unit {
	case seconds, minutes, hours:
		j.nextRun = j.lastRun.Add(periodDuration)
	case days:
		j.nextRun = j.roundToMidnight(j.lastRun)
		j.nextRun = j.nextRun.Add(j.atTime)
	case weeks:
		j.nextRun = j.roundToMidnight(j.lastRun)
		dayDiff := int(j.startDay)
		dayDiff -= int(j.nextRun.Weekday())
		if dayDiff != 0 {
			j.nextRun = j.nextRun.Add(time.Duration(dayDiff) * 24 * time.Hour)
		}
		j.nextRun = j.nextRun.Add(j.atTime)
	}
	// 更新任务下次运行时间
	for j.nextRun.Before(now) || j.nextRun.Before(j.lastRun) {
		j.nextRun = j.nextRun.Add(periodDuration)
	}
	return nil
}

// 获取任务的下次运行时间
func (j *Job) NexScheduledTime() time.Time {
	return j.nextRun
}

// 判断任务间隔时间
func (j *Job) mustInterval(i uint64) error {
	if j.interval != i {
		return fmt.Errorf("interval must be %d", i)
	}
	return nil
}

// 指定下次运行时间
func (j *Job) From(t *time.Time) *Job {
	j.nextRun = *t
	return j
}

// 指定时间间隔的单位
func (j *Job) setUnit(unit timeUnit) *Job {
	j.unit = unit
	return j
}

// 指定时间间隔为秒
func (j *Job) Seconds() *Job {
	return j.setUnit(seconds)
}

func (j *Job) Minutes() *Job {
	return j.setUnit(minutes)
}

func (j *Job) Hours() *Job {
	return j.setUnit(hours)
}

func (j *Job) Days() *Job {
	return j.setUnit(days)
}

func (j *Job) Weeks() *Job {
	return j.setUnit(weeks)
}

// 时间间隔为1
func (j *Job) Second() *Job {
	j.mustInterval(1)
	return j.Seconds()
}

func (j *Job) Minute() *Job {
	j.mustInterval(1)
	return j.Minutes()
}

func (j *Job) Hour() *Job {
	j.mustInterval(1)
	return j.Hours()
}

func (j *Job) Day() *Job {
	j.mustInterval(1)
	return j.Days()
}

func (j *Job) Week() *Job {
	j.mustInterval(1)
	return j.Weeks()
}

// 每周几运行
func (j *Job) Weekday(startDay time.Weekday) *Job {
	j.mustInterval(1)
	j.startDay = startDay
	return j.Weeks()
}

func (j *Job) GetWeekday() time.Weekday {
	return j.startDay
}

func (j *Job) Monday() *Job {
	return j.Weekday(time.Monday)
}

func (j *Job) Tuesday() *Job {
	return j.Weekday(time.Tuesday)
}

func (j *Job) Wednesday() *Job {
	return j.Weekday(time.Wednesday)
}

func (j *Job) Thursday() *Job {
	return j.Weekday(time.Thursday)
}

func (j *Job) Friday() *Job {
	return j.Weekday(time.Friday)
}

func (j *Job) Saturday() *Job {
	return j.Weekday(time.Saturday)
}

func (j *Job) Sunday() *Job {
	return j.Weekday(time.Sunday)
}

func (j *Job) Lock() *Job {
	j.lock = true
	return j
}
