package singleton

import (
	"sync"
	"testing"

	"github.com/go-playground/assert/v2"
)

const parCount = 100

func TestSingleton(t *testing.T) {
	t.Run("Singlaton", func(t *testing.T) {
		ins1 := GetInstance()
		ins2 := GetInstance()
		assert.Equal(t, ins1, ins2)
	})
	t.Run("parallel singleton", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(parCount)
		instances := [parCount]*Singleton{}
		for i := 0; i < parCount; i++ {
			go func(index int) {
				instances[index] = GetInstance()
				wg.Done()
			}(i)
		}
		wg.Wait()
		for i := 1; i < parCount; i++ {
			assert.Equal(t, instances[i], instances[i-1])
		}
	})
}
