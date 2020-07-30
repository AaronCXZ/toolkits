package prototype

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var manager *PrototypeManager

func init() {
	manager = NewPrototypeManager()
	t1 := &Type1{
		name: "type1",
	}
	manager.Set("t1", t1)
}

type Type1 struct {
	name string
}

func (t *Type1) Clone() Cloneable {
	tc := *t
	return &tc
}

type Type2 struct {
	name string
}

func (t *Type2) Clone() Cloneable {
	tc := *t
	return &tc
}

func TestClone(t *testing.T) {
	t.Run("get t1", func(t *testing.T) {
		t1 := manager.Get("t1")
		t2 := t1.Clone()
		assert.Equal(t, t1, t2)
	})
	t.Run("manager get", func(t *testing.T) {
		c := manager.Get("t1").Clone()
		t1 := c.(*Type1)
		assert.Equal(t, t1.name, "type1")
	})
}
