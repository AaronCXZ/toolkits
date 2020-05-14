package set

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet(1, 2, 3, 4, 5)
	fmt.Println(s)
	t.Run("add", func(t *testing.T) {
		s.Add(6, 7, 8)
		fmt.Println(s)
	})
	t.Run("delete", func(t *testing.T) {
		s.Del(1, 2)
		fmt.Println(s)
	})

	t.Run("len", func(t *testing.T) {
		fmt.Println(s.Len())
	})
	t.Run("Contains", func(t *testing.T) {
		fmt.Println(s.Contains(3))
		fmt.Println(s.Contains(1))
	})

	t.Run("clear", func(t *testing.T) {
		s.Clear()
		fmt.Println(s.Len())
	})
	t.Run("Equal", func(t *testing.T) {
		s1 := NewSet(1, 2, 3)
		s2 := NewSet(2, 3, 1)
		s3 := NewSet(1, 2, 3, 4)
		s4 := NewSet(1, 2, 3)
		fmt.Println(s1.Equal(s2))
		fmt.Println(s1.Equal(s3))
		fmt.Println(s1.Equal(s4))
	})
	t.Run("IsSubset", func(t *testing.T) {
		s1 := NewSet(1, 2, 3)
		s2 := NewSet(3, 1)
		s3 := NewSet(1, 2, 3, 4)
		s4 := NewSet()
		fmt.Println(s1.IsSubset(s2))
		fmt.Println(s1.IsSubset(s3))
		fmt.Println(s1.IsSubset(s4))
	})

}
