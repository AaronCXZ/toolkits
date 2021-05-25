package holmes

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
)

func parseUtil(s string, base, bitSize int) (uint64, error) {
	v, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, err
	}
	return v, nil
}

func readUint(path string) (uint64, error) {
	v, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return parseUtil(strings.TrimSpace(string(v)), 10, 64)
}

func trimResult(buffer bytes.Buffer) string {
	arr := strings.Split(buffer.String(), "\n\n")
	if len(arr) > 10 {
		arr = arr[:10]
	}
	return strings.Join(arr, "\n\n")
}

func getUsageCGroup() (float64, float64, int, int, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuPercent, err := p.Percent(1 * time.Second)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuPercent = cpuPercent / float64(runtime.GOMAXPROCS(-1))

	mem, err := p.MemoryInfo()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	memLimit, err := readUint(cgroupMemLimitPath)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	memPercent := float64(mem.RSS) * 100 / float64(memLimit)

	gNum := runtime.NumGoroutine()

	tNum := getThreadNum()

	return cpuPercent, memPercent, gNum, tNum, nil
}

func getUsageNormal() (float64, float64, int, int, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0, 0, 0, 0, err
	}

	cpuPercent, err := p.Percent(1 * time.Second)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	cpuPercent = cpuPercent / float64(runtime.GOMAXPROCS(-1))

	mem, err := p.MemoryPercent()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	gNum := runtime.NumGoroutine()
	tNum := getThreadNum()

	return cpuPercent, float64(mem), gNum, tNum, nil
}

func getThreadNum() int {
	return pprof.Lookup("threadcreate").Count()
}

var getUsage func() (float64, float64, int, int, error)

func collect() (int, int, int, int, error) {
	cpu, mem, gNum, tNum, err := getUsage()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return int(cpu), int(mem), gNum, tNum, nil
}

func matchRule(history ring, curVal, ruleMin, ruleAbs, ruleDiff int) bool {
	if curVal < ruleMin {
		return false
	}

	if curVal > ruleAbs {
		return true
	}

	avg := history.avg()
	return curVal >= avg*(100+ruleDiff)/100
}

func getBinaryFileName(filePath string, dumpType configureType) string {
	var (
		binarySuffix = time.Now().Format("20060102150405") + ".bin"
	)
	return path.Join(filePath, type2name[dumpType]+"."+binarySuffix)
}
