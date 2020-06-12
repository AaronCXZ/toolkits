package containsDuplicate

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestContainsDuplicate(t *testing.T) {
	nums := []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}
	got := ContainsDuplicate(nums)
	assert.Equal(t, true, got)
}
