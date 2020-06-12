package maxProfit

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestMaxProfit(t *testing.T) {
	prices := []int{7, 6, 4, 3, 1}
	got := MaxProfit(prices)
	want := 0
	assert.Equal(t, want, got)
}

func BenchmarkMaxProfit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prices := []int{7, 6, 4, 3, 1}
		MaxProfit(prices)
	}
}
