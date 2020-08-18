package strategy

import (
	"strconv"
	"strings"
)

// 生成计算时间需要的字符串
func makeCounts(req int) (counts []string) {
	str := []string{"_", "m"}
	for i := 1; i <= req; i++ {
		count := strings.Join(str, strconv.Itoa(i))
		counts = append(counts, count)
	}
	return counts
}

// 策略ID生成
func genStrategyID(name, body string) int64 {
	var id int64
	all := name + body
	if len(all) < 1 {
		return id
	}

	id = int64(all[0])

	for i := 1; i < len(all); i++ {
		id += int64(all[i])
		id += int64(all[i] - all[i-1])
	}

	id += 1000000 //避免与web端配置的id冲突
	return id
}
