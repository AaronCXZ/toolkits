package decorator

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestDecorator(t *testing.T) {
	var c Component = &ConcreteComponent{}
	c = WarpAddDecorator(c, 10)
	c = WarpMulDecorator(c, 8)
	c = WarpAddDecorator(c, 2)
	res := c.Calc()
	assert.Equal(t, res, 82)

}
