package removeDuplicates

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRemoveDuplicatesT(t *testing.T) {
	nums := []int{1, 1, 2}
	want := RemoveDuplicates(nums)
	assert.Equal(t, 2, want)
}

func BenchmarkRemoveDuplicates(b *testing.B) {
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	for i := 0; i < b.N; i++ {
		RemoveDuplicates(nums)
	}
}
