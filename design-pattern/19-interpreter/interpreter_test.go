package interpreter

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestInterpreter(t *testing.T) {
	p := &Parse{}
	p.Parse("1 + 2 + 3 - 4 + 5 - 6")
	res := p.Result().Interpret()
	expect := 1
	assert.Equal(t, res, expect)
}
