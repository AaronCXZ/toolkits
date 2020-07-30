package builder

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestBuilder(t *testing.T) {
	t.Run("builder1", func(t *testing.T) {
		builder := &Builder1{}
		director := NewDirector(builder)
		director.Construct()
		res := builder.GetResult()
		assert.Equal(t, res, "123")
	})
	t.Run("builder2", func(t *testing.T) {
		builder := &Builder2{}
		director := NewDirector(builder)
		director.Construct()
		res := builder.GetResult()
		assert.Equal(t, res, 6)
	})
}
