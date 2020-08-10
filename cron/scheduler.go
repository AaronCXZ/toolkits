package cron

import (
	"sort"
	"time"
)

// 按下次执行时间进行排序的任务列表
type Scheduler struct {
	jobs [MAXJOBNUM]*Job // 所有的任务
	size int             // 任务的数量
	loc  *time.Location  // 时区
}

var (
	defaultScheduler = NewScheduler()
)

// 创建新的任务列表
func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: [10000]*Job{},
		size: 0,
		loc:  loc,
	}
}

// 所有的任务
func (s *Scheduler) Jobs() []*Job {
	return s.jobs[:s.size]
}

// 实现sort接口
func (s *Scheduler) Len() int {
	return s.size
}

func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[j].nextRun.Unix() >= s.jobs[i].nextRun.Unix()
}

// 修改时区
func (s *Scheduler) ChangeLoc(newLocation *time.Location) {
	s.loc = newLocation
}

// 需要立即运行的任务以及个数
func (s *Scheduler) getRunningJobs() (runningJobs [MAXJOBNUM]*Job, n int) {
	// 根据任务的下次运行时间排序
	sort.Sort(s)
	for i := 0; i < s.size; i++ {
		if s.jobs[i].shouldRun() {
			runningJobs[n] = s.jobs[i]
			n++
		} else {
			break
		}
	}
	return runningJobs, n
}

// 下一个运行的任务以及运行时间
func (s *Scheduler) NextRun() (*Job, time.Time) {
	if s.size <= 0 {
		return nil, time.Now()
	}
	sort.Sort(s)
	return s.jobs[0], s.jobs[0].nextRun
}

// 根据时间间隔创建新任务
func (s *Scheduler) Every(interval uint64) *Job {
	job := NewJob(interval).Loc(s.loc)
	s.jobs[s.size] = job
	s.size++
	return job
}

// 运行任务
func (s *Scheduler) RunPending() {
	runningJobs, n := s.getRunningJobs()
	if n != 0 {
		// 逐个运行任务并更新最后运行时间和下次运行时间
		for i := 0; i < n; i++ {
			go runningJobs[i].run()
			runningJobs[i].lastRun = time.Now()
			runningJobs[i].scheduleNextRun()
		}
	}
}

// 延迟运行所有任务
func (s *Scheduler) RunAllWithDelay(d int) {
	for i := 0; i < s.size; i++ {
		go s.jobs[i].run()
		if 0 != d {
			time.Sleep(time.Duration(d))
		}
	}
}

// 立即运行所有任务
func (s *Scheduler) RunAll() {
	s.RunAllWithDelay(0)
}

// 删除任务
func (s *Scheduler) Remove(j interface{}) {
	s.removeByCondition(func(someJob *Job) bool {
		return someJob.jobFunc == getFunctionName(j)
	})
}

// 删除任务
func (s *Scheduler) RemoveByRef(j *Job) {
	s.removeByCondition(func(someJob *Job) bool {
		return someJob == j
	})
}

// 根据函数删除任务
func (s *Scheduler) removeByCondition(shouldRemove func(*Job) bool) {
	i := 0
	for {
		found := false
		for ; i < s.size; i++ {
			if shouldRemove(s.jobs[i]) {
				found = true
				break
			}
		}

		if !found {
			return
		}

		for j := i + 1; j < s.size; j++ {
			s.jobs[i] = s.jobs[j]
			i++
		}
		s.size--
		s.jobs[s.size] = nil
	}
}

// 是否存在任务
func (s *Scheduler) Scheduled(j interface{}) bool {
	for _, job := range s.jobs {
		if job.jobFunc == getFunctionName(j) {
			return true
		}
	}
	return false
}

// 清空任务
func (s *Scheduler) Clear() {
	for i := 0; i < s.size; i++ {
		s.jobs[i] = nil
	}
	s.size = 0
}

// 启动任务
func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	// 每秒执行一次 RunPending
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
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

// 新建任务
func Every(interval uint64) *Job {
	return defaultScheduler.Every(interval)
}

// 执行到运行时间的任务
func RunPending() {
	defaultScheduler.RunPending()
}

// 立即执行所有任务
func RunAll() {
	defaultScheduler.RunAll()
}

// 延迟运行所有任务
func RunAllWithDelay(d int) {
	defaultScheduler.RunAllWithDelay(d)
}

// 启动任务
func Start() chan bool {
	return defaultScheduler.Start()
}

// 清空任务
func Clear() {
	defaultScheduler.Clear()
}

// 删除任务
func Remove(j interface{}) {
	defaultScheduler.Remove(j)
}

// 是否存在任务
func Scheduled(j interface{}) bool {
	for _, job := range defaultScheduler.jobs {
		if job.jobFunc == getFunctionName(j) {
			return true
		}
	}
	return false
}

// 下次运行的任务及时间
func NextRun() (job *Job, time time.Time) {
	return defaultScheduler.NextRun()
}
