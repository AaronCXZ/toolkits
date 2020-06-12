package singleNumber

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSingleNumber(t *testing.T) {
	nums := []int{4, 1, 2, 1, 2}
	got := SingleNumber(nums)
	want := 4
	assert.Equal(t, want, got)
}
