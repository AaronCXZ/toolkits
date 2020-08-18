package strategy

import (
	"strconv"
	"strings"
)

func makeCounts(req int) (counts []string) {
	str := []string{"_", "m"}
	for i := 1; i <= req; i++ {
		count := strings.Join(str, strconv.Itoa(i))
		counts = append(counts, count)
	}
	return counts
}
