package rotate

import (
	"fmt"
	"testing"
)

func TestRotate(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	Rotate(nums, 3)
	fmt.Println(nums)
}
