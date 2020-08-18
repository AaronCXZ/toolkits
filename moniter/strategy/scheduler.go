package strategy

import (
	"sort"
	"time"
)

const MAXJOBNUM = 1000

type Scheduler struct {
	strategies [MAXJOBNUM]*Strategy
	size       int
	alerts     chan Alerter
}

var defauleScheruler = NewScheduler()

func NewScheduler() *Scheduler {
	return &Scheduler{
		strategies: [MAXJOBNUM]*Strategy{},
		size:       0,
		alerts:     make(chan Alerter, 1),
	}
}

// 获取所有任务的方法
func (s *Scheduler) Strategies() []*Strategy {
	return s.strategies[:s.size]
}

// 实现sort接口
func (s *Scheduler) Len() int {
	return s.size
}

func (s *Scheduler) Swap(i, j int) {
	s.strategies[i], s.strategies[j] = s.strategies[j], s.strategies[i]
}

func (s *Scheduler) Less(i, j int) bool {
	return s.strategies[j].nextRun.Unix() >= s.strategies[i].nextRun.Unix()
}

// 需要立即运行的策略
func (s *Scheduler) getRunningStrategies() (runningStrategies [MAXJOBNUM]*Strategy, n int) {
	sort.Sort(s)
	for i := 0; i < s.size; i++ {
		if s.strategies[i].shouldRun() {
			runningStrategies[n] = s.strategies[i]
			n++
		} else {
			break
		}
	}
	return runningStrategies, n
}

// 下次运行的策略及其运行时间
func (s *Scheduler) NextRun() (*Strategy, time.Time) {
	if s.size == 0 {
		return nil, time.Now()
	}
	sort.Sort(s)
	return s.strategies[0], s.strategies[0].nextRun
}

// 运行策略
func (s *Scheduler) RunPending() {
	runningStrategies, n := s.getRunningStrategies()
	if n != 0 {
		for i := 0; i < n; i++ {
			go runningStrategies[i].do()
			runningStrategies[i].lastRun = time.Now()
			runningStrategies[i].scheduleNextRun()
			runningStrategies[i].updateCmd()
		}
	}
}

// 添加策略
func (s *Scheduler) Add(strategies ...*Strategy) {
	for _, strategy := range strategies {
		if err := strategy.checkStrategy(); err != nil {
			s.strategies[s.size] = strategy
			s.size += 1
		}
	}
}

// 清空策略
func (s *Scheduler) Clear() {
	for i := 0; i < s.size; i++ {
		s.strategies[i] = nil
	}
	s.size = 0
}

// 启动
func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case alert := <-s.alerts:
				alert.Send()
			case <-ticker.C:
				s.RunPending()
			case <-stopped:
				ticker.Stop()
				return
			}
		}
	}()
	return stopped
}

func Add(strategies ...*Strategy) {
	defauleScheruler.Add(strategies...)
}

func Start() chan bool {
	return defauleScheruler.Start()
}

func Clear() {
	defauleScheruler.Clear()
}
