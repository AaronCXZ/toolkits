package main

import "testing"

func TestArrayInsert(t *testing.T) {
	capacity := 10
	arr := NewArray(uint(capacity))
	for i := 0; i < capacity-2; i++ {
		err := arr.Insert(uint(i), i+1)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	arr.Print()

	arr.Insert(uint(6), 999)
	arr.Print()

	arr.InsertToTail(666)
	arr.Print()
}

func TestArrayDelete(t *testing.T) {
	capacity := 10
	arr := NewArray(uint(capacity))
	for i := 0; i < capacity; i++ {
		if err := arr.Insert(uint(i), i+1); err != nil {
			t.Fatal(err.Error())
		}
	}
	arr.Print()

	for i := 9; i >= 0; i-- {
		if _, err := arr.Delete(uint(i)); err != nil {
			t.Fatal(err.Error())
		}
		arr.Print()
	}
}

func TestArrayFind(t *testing.T) {
	capacity := 10
	arr := NewArray(uint(capacity))
	for i := 0; i < capacity; i++ {
		if err := arr.Insert(uint(i), i+1); err != nil {
			t.Fatal(err.Error())
		}
	}
	arr.Print()

	t.Log(arr.Find(0))
	t.Log(arr.Find(9))
	t.Log(arr.Find(11))
}
