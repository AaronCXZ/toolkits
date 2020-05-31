package set

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet(1, 2, 3, 4, 5)
	fmt.Println(s)
	t.Run("adds", func(t *testing.T) {
		s.Adds(6, 7, 8)
		fmt.Println(s)
	})
	t.Run("add", func(t *testing.T) {
		s.Add(9)
		fmt.Println(s)
	})
	t.Run("removes", func(t *testing.T) {
		s.Removes(1, 2)
		fmt.Println(s)
	})
	t.Run("remove", func(t *testing.T) {
		s.Remove(3)
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
		fmt.Println(s)
	})

	t.Run("same", func(t *testing.T) {
		s1 := NewSet(1, 2, 3)
		s2 := NewSet(2, 3, 1)
		s3 := NewSet(1, 2, 3, 4)
		s4 := NewSet(1, 2, 3)
		fmt.Println(s1.Same(s2))
		fmt.Println(s1.Same(s3))
		fmt.Println(s1.Same(s4))
	})
	t.Run("Elements", func(t *testing.T) {
		s1 := NewSet(1, 2, 3)
		fmt.Println(s1.Elements())
	})

}
