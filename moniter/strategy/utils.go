package strategy

import (
	"fmt"
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

// 文件名校验，log.xxx.json
func checkName(f string) (err error) {
	if !strings.Contains(f, "log.") {
		err = fmt.Errorf("name illege not contain log. %s", f)
		return
	}

	arr := strings.Split(f, ".")
	if len(arr) < 3 {
		err = fmt.Errorf("name illege %s len:%d len < 3 ", f, len(arr))
		return
	}

	if arr[len(arr)-1] != "json" {
		err = fmt.Errorf("name illege %s not json file", f)
		return
	}

	return
}
