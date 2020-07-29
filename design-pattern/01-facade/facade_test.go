package facade

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var (
	expectAPI = "A module running\nB module running"
	expectA   = "A module running"
	expectB   = "B module running"
)

func TestFacadeAPI(t *testing.T) {
	t.Run("api", func(t *testing.T) {
		api := NewAPI()
		ret := api.Test()
		assert.Equal(t, ret, expectAPI)
	})
	t.Run("a module", func(t *testing.T) {
		a := NewAModuleAPI()
		aRet := a.TestA()
		assert.Equal(t, aRet, expectA)
	})
	t.Run("b module", func(t *testing.T) {
		b := NewBModuleAPI()
		bRet := b.TestB()
		assert.Equal(t, bRet, expectB)
	})
}
