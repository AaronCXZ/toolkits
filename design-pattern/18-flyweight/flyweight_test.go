package flyweight

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func ExampleFlyweight() {
	viewer := NewImageViewer("image1.jpg")
	viewer.Display()
}

func TestFlyweight(t *testing.T) {
	ExampleFlyweight()
	viewer1 := NewImageViewer("image1.jpg")
	viewer2 := NewImageViewer("image1.jpg")
	assert.Equal(t, viewer1.ImageFlyweight, viewer2.ImageFlyweight)
}
