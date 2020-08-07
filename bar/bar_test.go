package bar

import (
	"testing"
	"time"
)

func TestBar(t *testing.T) {
	t.Run("option", func(t *testing.T) {
		var bar Bar
		bar.NewOption(1, 100)
		for i := 0; i <= 1000; i++ {
			time.Sleep(100 * time.Millisecond)
			bar.Play(int64(i))
		}
		bar.Finish()
	})
	t.Run("option with graph", func(t *testing.T) {
		var bar Bar
		bar.NewOptionWithGraph(1, 100, "*")
		for i := 0; i <= 1000; i++ {
			time.Sleep(100 * time.Millisecond)
			bar.Play(int64(i))
		}
		bar.Finish()
	})
}
