package mediator

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestMediator(t *testing.T) {
	mediator := GetMediatorInstance()
	mediator.CD = &CDDriver{}
	mediator.CPU = &CPU{}
	mediator.Video = &VideoCard{}
	mediator.Sound = &SoundCard{}

	mediator.CD.ReadData()
	assert.Equal(t, mediator.CD.Data, "music,image")
	assert.Equal(t, mediator.CPU.Sound, "music")
	assert.Equal(t, mediator.CPU.Video, "image")
	assert.Equal(t, mediator.Video.Data, "image")
	assert.Equal(t, mediator.Sound.Data, "music")
}
