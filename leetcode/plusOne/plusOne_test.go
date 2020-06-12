package plusOne

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPlusOne(t *testing.T) {
	nums := []int{4, 3, 2, 1}
	got := PlusOne(nums)
	want := []int{4, 3, 2, 2}
	assert.Equal(t, got, want)
}
