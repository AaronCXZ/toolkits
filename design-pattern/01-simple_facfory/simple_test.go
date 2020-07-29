package implefacfory

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAPI(t *testing.T) {
	t.Run("type Hi API", func(t *testing.T) {
		api := NewAPI(1)
		got := api.Say("Tom")
		want := "Hi, Tom"
		assert.Equal(t, got, want)
	})
	t.Run("type Hello API", func(t *testing.T) {
		api := NewAPI(2)
		got := api.Say("Tom")
		want := "Hello, Tom"
		assert.Equal(t, got, want)
	})
}
