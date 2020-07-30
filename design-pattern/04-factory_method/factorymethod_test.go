package factorymethod

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func compute(factory OperatorFactory, a, b int) int {
	op := factory.Create()
	op.SetA(a)
	op.SetB(b)
	return op.Result()
}

func TestOperator(t *testing.T) {
	var factory OperatorFactory
	t.Run("plus", func(t *testing.T) {
		factory = PlusOperatorFactory{}
		got := compute(factory, 1, 2)
		assert.Equal(t, got, 3)
	})
	t.Run("minus", func(t *testing.T) {
		factory = MinusOperatorFactory{}
		got := compute(factory, 4, 2)
		assert.Equal(t, got, 2)
	})

}
