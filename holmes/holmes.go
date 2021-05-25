package holmes

import (
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
	"sync/atomic"
	"time"
)

type Holmes struct {
	opts *options

	collectCount       int
	threadTriggerCount int
	cpuTriggerCount    int
	memTriggerCount    int
	grTriggerCount     int

	threadCoolDumpTime time.Time
	cpuCoolDumpTime    time.Time
	memCoolDumpTime    time.Time
	grCoolDumpTime     time.Time

	threadStats ring
	cpuStats    ring
	memStats    ring
	grNumStats  ring

	stopped int64
}

func New(opts ...Option) (*Holmes, error) {
	holmes := &Holmes{
		opts: newOptions(),
	}
	for _, opt := range opts {
		if err := opt.apply(holmes.opts); err != nil {
			return nil, err
		}
	}
	return holmes, nil
}

func (h *Holmes) EnableThreadDump() *Holmes {
	h.opts.ThreadOpts.Enable = true
	return h
}

func (h *Holmes) DisableThreadDump() *Holmes {
	h.opts.ThreadOpts.Enable = false
	return h
}

func (h *Holmes) EnableGoroutineDump() *Holmes {
	h.opts.GrOpts.Enable = true
	return h
}

func (h *Holmes) DisableGoroutineDump() *Holmes {
	h.opts.GrOpts.Enable = false
	return h
}

func (h *Holmes) EnableCPUDump() *Holmes {
	h.opts.CPUOpts.Enable = true
	return h
}

func (h *Holmes) DisableCPUDump() *Holmes {
	h.opts.CPUOpts.Enable = false
	return h
}

func (h *Holmes) EnableMemDump() *Holmes {
	h.opts.MemOpts.Enable = true
	return h
}

func (h *Holmes) DisableMemDump() *Holmes {
	h.opts.MemOpts.Enable = false
	return h
}

func (h *Holmes) Start() {
	atomic.StoreInt64(&h.stopped, 0)
	h.initEnvironment()
	go h.startDumpLoop()
}

func (h *Holmes) Stop() {
	atomic.StoreInt64(&h.stopped, 1)
}

func (h *Holmes) startDumpLoop() {
	now := time.Now()
	h.cpuCoolDumpTime = now
	h.memCoolDumpTime = now
	h.grCoolDumpTime = now

	h.cpuStats = newRing(minCollectCyclesBeforeDumpStart)
	h.memStats = newRing(minCollectCyclesBeforeDumpStart)
	h.grNumStats = newRing(minCollectCyclesBeforeDumpStart)
	h.threadStats = newRing(minCollectCyclesBeforeDumpStart)

	ticker := time.NewTicker(h.opts.CollectInterval)
	defer ticker.Stop()
	for range ticker.C {
		if atomic.LoadInt64(&h.stopped) == 1 {
			fmt.Println("[Holmes] dump loop stopped")
			return
		}

		cpu, mem, gNum, tNum, err := collect()
		if err != nil {
			h.logf(err.Error())
			continue
		}

		h.cpuStats.Push(cpu)
		h.memStats.Push(mem)
		h.grNumStats.Push(gNum)
		h.threadStats.Push(tNum)

		h.collectCount++
		if h.collectCount < minCollectCyclesBeforeDumpStart {
			h.logf("[Holmes] warming up cycle : %d", h.collectCount)
			continue
		}
		h.goroutineCheckAndDump(gNum)
		h.memCheckAndDump(mem)
		h.cpuCheckAndDump(cpu)
		h.threadCheckAndDump(tNum)
	}
}

func (h *Holmes) goroutineCheckAndDump(gNum int) {
	if !h.opts.GrOpts.Enable {
		return
	}
	if h.grCoolDumpTime.After(time.Now()) {
		h.logf("[Holmes] goroutine dump is in cooldown")
		return
	}
	if triggered := h.goroutineProfile(gNum); triggered {
		h.grCoolDumpTime = time.Now().Add(h.opts.CoolDown)
		h.grTriggerCount++
	}
}

func (h *Holmes) goroutineProfile(gNum int) bool {
	c := h.opts.GrOpts
	if !matchRule(h.grNumStats, gNum, c.GoroutineTriggerNumMin, c.GoroutineTriggerNumAbs, c.GoroutineTriggerPercentDiff) {
		h.debugUniform("NODUMP", type2name[goroutine],
			c.GoroutineTriggerNumMin, c.GoroutineTriggerPercentDiff, c.GoroutineTriggerNumAbs,
			h.grNumStats.data, gNum)
		return false
	}
	var buf bytes.Buffer
	_ = pprof.Lookup("goroutine").WriteTo(&buf, int(h.opts.DumpProfileType))
	h.writeProfileDataToFile(buf, goroutine, gNum)
	return true
}

func (h *Holmes) memCheckAndDump(mem int) {
	if !h.opts.MemOpts.Enable {
		return
	}
	if h.memCoolDumpTime.After(time.Now()) {
		h.logf("[Holmes] mem dump is in cooldown")
		return
	}
	if triggered := h.memProfile(mem); triggered {
		h.memCoolDumpTime = time.Now().Add(h.opts.CoolDown)
		h.memTriggerCount++
	}
}

func (h *Holmes) memProfile(rss int) bool {
	c := h.opts.MemOpts
	if !matchRule(h.memStats, rss, c.MemTriggerPercentMin, c.MenTriggerPercentAbs, c.MemTriggerPercentDiff) {
		h.debugUniform("NUDUMP", type2name[mem],
			c.MemTriggerPercentMin, c.MemTriggerPercentDiff, c.MenTriggerPercentAbs,
			h.memStats.data, rss)
		return false
	}
	var buf bytes.Buffer
	_ = pprof.Lookup("heap").WriteTo(&buf, int(h.opts.DumpProfileType))
	return true
}

func (h *Holmes) threadCheckAndDump(threadNum int) {
	if !h.opts.ThreadOpts.Enable {
		return
	}
	if h.threadCoolDumpTime.After(time.Now()) {
		h.logf("[Holmes] thread dump is in cooldown")
		return
	}
	if triggered := h.threadProfile(threadNum); triggered {
		h.threadCoolDumpTime = time.Now().Add(h.opts.CoolDown)
		h.threadTriggerCount++
	}
}

func (h *Holmes) threadProfile(curThreadNum int) bool {
	c := h.opts.ThreadOpts
	if !matchRule(h.threadStats, curThreadNum, c.ThreadTriggerPercentMin, c.ThreadTriggerPercentAbs, c.ThreadTriggerPercentDiff) {
		h.debugUniform("NODUMP", type2name[thread],
			c.ThreadTriggerPercentMin, c.ThreadTriggerPercentDiff, c.ThreadTriggerPercentAbs,
			h.threadStats.data, curThreadNum)
		return false
	}

	var buf bytes.Buffer
	pprof.Lookup("threadcreate").WriteTo(&buf, int(h.opts.DumpProfileType))
	pprof.Lookup("goroutine").WriteTo(&buf, int(h.opts.DumpProfileType))

	h.writeProfileDataToFile(buf, thread, curThreadNum)
	return true
}

func (h *Holmes) cpuCheckAndDump(cpu int) {
	if !h.opts.CPUOpts.Enable {
		return
	}

	if h.cpuCoolDumpTime.After(time.Now()) {
		h.logf("[Holmes] cpu dump is in cooldown")
		return
	}

	if triggered := h.cpuProfile(cpu); triggered {
		h.cpuCoolDumpTime = time.Now().Add(h.opts.CoolDown)
		h.cpuTriggerCount++
	}
}

func (h *Holmes) cpuProfile(curCPUUsage int) bool {
	c := h.opts.CPUOpts
	if !matchRule(h.cpuStats, curCPUUsage, c.CPUTriggerPercentMin, c.CPUTriggerPercentAbs, c.CPUTriggerPercentDiff) {
		h.debugUniform("NODUMP", type2name[cpu],
			c.CPUTriggerPercentMin, c.CPUTriggerPercentDiff, c.CPUTriggerPercentAbs,
			h.cpuStats.data, curCPUUsage)
		return false
	}

	binFileName := getBinaryFileName(h.opts.DumpPath, cpu)

	bf, err := os.OpenFile(binFileName, defaultLoggerFlags, defaultLoggerPerm)
	if err != nil {
		h.logf("[Holmes] failed to create cpu profile file: %v", err.Error())
		return false
	}
	defer bf.Close()
	err = pprof.StartCPUProfile(bf)
	if err != nil {
		h.logf("[Holmes] failed to profile cpu: %v", err.Error())
		return false
	}

	time.Sleep(defaultCPUSamplingTime)
	pprof.StopCPUProfile()

	h.infoUniform("pprof dump to log dir", type2name[cpu],
		c.CPUTriggerPercentMin, c.CPUTriggerPercentDiff, c.CPUTriggerPercentAbs,
		h.cpuStats.data, curCPUUsage)

	return true
}

func (h *Holmes) writeProfileDataToFile(data bytes.Buffer, dumpType configureType, currentStat int) {
	binFileName := getBinaryFileName(h.opts.DumpPath, dumpType)
	switch dumpType {
	case mem:
		opts := h.opts.MemOpts
		h.infoUniform("pprof", type2name[dumpType],
			opts.MemTriggerPercentMin, opts.MemTriggerPercentDiff, opts.MenTriggerPercentAbs,
			h.memStats.data, currentStat)
	case goroutine:
		opts := h.opts.GrOpts
		h.infoUniform("pprof", type2name[dumpType],
			opts.GoroutineTriggerNumMin, opts.GoroutineTriggerPercentDiff, opts.GoroutineTriggerNumAbs,
			h.grNumStats.data, currentStat)

	case thread:
		opts := h.opts.ThreadOpts
		h.infoUniform("pporf", type2name[dumpType],
			opts.ThreadTriggerPercentMin, opts.ThreadTriggerPercentDiff, opts.ThreadTriggerPercentAbs,
			h.threadStats.data, currentStat)
	}

	if h.opts.DumpProfileType == textDump {
		var res = data.String()
		if !h.opts.DumpFullStack {
			res = trimResult(data)
		}
		h.logf(res)
	} else {
		bf, err := os.OpenFile(binFileName, defaultLoggerFlags, defaultLoggerPerm)
		if err != nil {
			h.logf("[Holmes] pprof %v write to file failed : %v", type2name[dumpType], err.Error())
			return
		}
		defer bf.Close()

		if _, err = bf.Write(data.Bytes()); err != nil {
			h.logf("[Holmes] pprof %v write to file failed : %v", type2name[dumpType], err.Error())
		}
	}
}

func (h *Holmes) debugUniform(msg string, name string, min, diff, abs int, data []int, cur int) {
	h.debugf("[Holmes] %v %v, config_min : %v, config_diff : %v, config_abs : %v, previous : %v, current: %v",
		msg, name, min, diff, abs, data, cur)
}

func (h *Holmes) infoUniform(msg string, name string, min, diff, abs int, data []int, cur int) {
	h.logf("[Holmes] %v %v, config_min : %v, config_diff : %v, config_abs : %v, previous : %v, current: %v",
		msg, name, min, diff, abs, data, cur)
}

func (h *Holmes) initEnvironment() {
	// choose whether the max memory is limited by cgroup
	if h.opts.UseCGroup {
		// use cgroup
		getUsage = getUsageCGroup
		h.logf("[Holmes] use cgroup to limit memory")
	} else {
		// not use cgroup
		getUsage = getUsageNormal
		h.logf("[Holmes] use the default memory percent calculated by gopsutil")
	}
}
