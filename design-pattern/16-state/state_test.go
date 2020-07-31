package state

import "testing"

func ExampleWeek() {
	ctx := NewDayContext()
	todayAndNext := func() {
		ctx.Today()
		ctx.Next(ctx)
	}
	for i := 0; i < 8; i++ {
		todayAndNext()
	}
}

func TestWeek(t *testing.T) {
	ExampleWeek()
}
