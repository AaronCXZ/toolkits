package str

import (
	"fmt"
	"testing"
)

func TestRand(t *testing.T) {
	t.Run("RandLetters", func(t *testing.T) {
		s := RandLetters(16)
		fmt.Println(s)
	})
	t.Run("RandDigits", func(t *testing.T) {
		s := RandDigits(16)
		fmt.Println(s)
	})
}
