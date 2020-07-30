package adapter

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var expect = "adaptee method"

func TestAdaptee(t *testing.T) {
	adaptee := NewAdaptee()
	target := Newadapter(adaptee)
	got := target.Request()
	assert.Equal(t, got, expect)
}
