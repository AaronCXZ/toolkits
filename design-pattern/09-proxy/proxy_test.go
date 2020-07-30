package proxy

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestProxy(t *testing.T) {
	var sub Subject
	sub = &Proxy{}

	res := sub.Do()
	assert.Equal(t, res, "pre:real:after")
}
