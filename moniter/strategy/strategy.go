package strategy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Muskchen/toolkits/ssh"
)

var (
	ErrNotStrategyName = errors.New("策略名称为空")
)

const (
	PATTERN_EXCLUDE_PARTITION = "```EXCLUDE```"
	TIME_FORMAT               = "2006-01-02 15:04"
)

type Strategy struct {
	ID int64 `json:"id"`
	// 配置文件读取
	Name      string   `json:"name"`
	Hosts     []string `json:"hosts"`
	User      string   `json:"user,omitempty"`
	Password  string   `json:"password,omitempty"`
	FileName  string   `json:"filename"`  // 文件路径
	Interval  int      `json:"interval"`  // 频率
	Threshold int      `json:"threshold"` // 阈值
	Pattern   string   `json:"pattern"`   // 用户正则
	Phones    string   `json:"phones"`
	Mails     string   `json:"mails"`
	Enable    bool     `json:"enable,-"`
	//
	Exclude string `json:"exclude,omitempty"`
	lastRun time.Time
	nextRun time.Time
	cmd     string
}

func (st *Strategy) shouldRun() bool {
	return time.Now().Unix() >= st.nextRun.Unix()
}

func (st *Strategy) run(host string) (count int, err error) {
	fileExit := fmt.Sprintf("if [ ! -f %s ];then\n echo \"1\"\n else echo \"2\"\n fi\n", st.FileName)
	newSSH := ssh.NewSSH(st.User, st.Password, host, 22)
	if err := newSSH.Start(); err != nil {
		return 0, err
	}
	defer newSSH.Stop()
	exitOut, err := newSSH.Run(fileExit)
	if err != nil {
		return 0, err
	}
	exit, _ := strconv.Atoi(strings.Replace(exitOut, "\n", "", -1))
	if exit == 1 {
		return
	}
	out, err := newSSH.Run(st.cmd)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.Replace(out, "\n", "", -1))
}

func (st *Strategy) do() {
	for _, host := range st.Hosts {
		count, err := st.run(host)
		if err != nil {
			return
		}
		if alert := st.makeAlert(count, host); alert != nil {
			defauleScheruler.alerts <- alert
		}
	}
}

func (st *Strategy) scheduleNextRun() {
	now := time.Now()
	if st.lastRun == time.Unix(0, 0) {
		st.lastRun = now
	}
	st.nextRun = st.lastRun.Add(st.periodDuration())
	if st.nextRun.Before(now) || st.nextRun.Before(st.lastRun) {
		st.nextRun = st.nextRun.Add(st.periodDuration())
	}
}

func (st *Strategy) parseKeyword() string {
	if st.Exclude == "" {
		return fmt.Sprintf("egrep '%s'", st.Pattern)
	} else {
		return fmt.Sprintf("egrep '%s' | egrep -v '%s'", st.Pattern, st.Exclude)
	}
}

func (st *Strategy) periodDuration() time.Duration {
	return time.Duration(st.Interval) * time.Minute
}

func (st *Strategy) checkStrategy() error {
	if st.Name == "" {
		return ErrNotStrategyName
	}
	if st.User == "" {
		st.User = "monitor"
	}
	if st.Password == "" {
		st.Password = "monitor"
	}
	if st.Interval == 0 {
		st.Interval = 5
	}
	return nil
}

func (st *Strategy) updateCmd() {
	var times []string
	counts := makeCounts(st.Interval)
	for _, count := range counts {
		h, _ := time.ParseDuration(count)
		TIME := st.nextRun.Add(h).Format(TIME_FORMAT)
		times = append(times, TIME)
	}
	grepKey := st.parseKeyword()
	st.cmd = fmt.Sprintf("cat %s|egrep '%s'| %s |wc -l", st.FileName, strings.Join(times, "|"), grepKey)
}

func (st *Strategy) makeAlert(count int, host string) (alert Alerter) {
	if count >= st.Threshold {
		context := fmt.Sprintf("服务器：%s\n项目：%s\n监控间隔：%d分钟\n表达式：%s\n出现次数：%d次\n阈值：%d\n", host, st.Name, st.Interval, st.Pattern, count, st.Threshold)
		sub := fmt.Sprintf("%s 出现%s", st.Name, st.Pattern)
		phones := st.Phones
		mails := st.Mails
		return NewAlerter(phones, mails, context, sub)
	}
	return nil
}
