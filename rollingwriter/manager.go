package rollingwriter

import (
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type manager struct {
	thresholdSize int64
	startAt       time.Time
	fire          chan string
	cr            *cron.Cron
	context       chan int
	wg            sync.WaitGroup
	lock          sync.Mutex
}

func NewManager(c *Config) (Manager, error) {
	m := &manager{
		startAt: time.Now(),
		fire:    make(chan string),
		cr:      cron.New(),
		context: make(chan int),
		wg:      sync.WaitGroup{},
	}

	// 判断日志滚动模式
	switch c.RollingPolicy {
	default:
		fallthrough
	case WithoutRolling:
		return m, nil
	case TimeRolling:
		if _, err := m.cr.AddFunc(c.RollingTimePattern, func() {
			m.fire <- m.GenLogFileName(c)
		}); err != nil {
			return nil, err
		}
		m.cr.Start()
	case VolumeRolling:
		m.ParseVolume(c)
		m.wg.Add(1)
		go func() {
			// 每秒一次的计时器
			timer := time.Tick(time.Duration(Precision) * time.Second)
			filepath := LogFilePath(c)
			var file *os.File
			var err error
			m.wg.Done()

			// 触发滚动或关闭，关闭时退出循环
			for {
				select {
				// 关闭chan
				case <-m.context:
					return
				//	每秒一次检查当前日志文件大小
				case <-timer:
					if file, err = os.Open(filepath); err != nil {
						continue
					}
					// 判断是否触发滚动
					if info, err := file.Stat(); err == nil && info.Size() > m.thresholdSize {
						m.fire <- m.GenLogFileName(c)
					}
					_ = file.Close()
				}
			}
		}()
		m.wg.Wait()
	}
	return m, nil
}

func (m *manager) Fire() chan string {
	return m.fire
}

func (m *manager) Close() {
	close(m.context)
	m.cr.Stop()
}

// 生成新的历史日志文件名称，更新startAt为当前时间
func (m *manager) GenLogFileName(c *Config) (filename string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if c.Compress {
		filename = path.Join(c.LogPath, c.FileName+".log.gz."+m.startAt.Format(c.TimeTagFormat))
	} else {
		filename = path.Join(c.LogPath, c.FileName+".log."+m.startAt.Format(c.TimeTagFormat))
	}
	m.startAt = time.Now()
	return filename
}

// 根据配置更新m.thresholdSize
func (m *manager) ParseVolume(c *Config) {
	// 读取大小滚动策略时的截断大小
	s := []byte(strings.ToUpper(c.RollingVolumeSize))
	// 如果不包含单位，则thresholdSize为1G
	if !(strings.Contains(string(s), "K") || strings.Contains(string(s), "KB") ||
		strings.Contains(string(s), "M") || strings.Contains(string(s), "MB") ||
		strings.Contains(string(s), "G") || strings.Contains(string(s), "GB") ||
		strings.Contains(string(s), "T") || strings.Contains(string(s), "TB")) {

		m.thresholdSize = 1024 * 1024 * 1024
		return
	}

	var unit int64 = 1
	p, _ := strconv.Atoi(string(s[:len(s)-1]))
	unitstr := string(s[len(s)-1])

	if s[len(s)-1] == 'B' {
		p, _ = strconv.Atoi(string(s[:len(s)-2]))
		unitstr = string(s[len(s)-2:])
	}
	// 使用fallthrough的switch，匹配到的case分支后面的case分支都会执行，后面的default分支不执行
	switch unitstr {
	default:
		fallthrough
	case "T", "TB":
		unit *= 1024
		fallthrough
	case "G", "GB":
		unit *= 1024
		fallthrough
	case "M", "MB":
		unit *= 1024
		fallthrough
	case "K", "KB":
		unit *= 1024
	}
	m.thresholdSize = int64(p) * unit
}
