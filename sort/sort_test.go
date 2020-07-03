package sort

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var arr []int = []int{5, 99, 9, 41, 7, 52, 61}
var arrSort []int = []int{5, 7, 9, 41, 52, 61, 99}

func TestSort(t *testing.T) {
	t.Run("bubble sort", func(t *testing.T) {
		BubbleSort(arr)
		assert.Equal(t, arr, arrSort)
	})
	t.Run("selection sort", func(t *testing.T) {
		SelectionSort(arr)
		assert.Equal(t, arr, arrSort)
	})
	t.Run("insertion sort", func(t *testing.T) {
		InsertionSort(arr)
		assert.Equal(t, arr, arrSort)
	})
	t.Run("shell sort", func(t *testing.T) {
		ShellSort(arr)
		assert.Equal(t, arr, arrSort)

		a := []int{5, 41, 7, 66, 25, 11, 44, 2}
		want := []int{2, 5, 7, 11, 25, 41, 44, 66}
		ShellSort(a)
		assert.Equal(t, a, want)
	})

	t.Run("merge sort", func(t *testing.T) {
		mergeSort := MergeSort(arr)
		assert.Equal(t, mergeSort, arrSort)
	})
	t.Run("Quick sort", func(t *testing.T) {
		QuickSort(arr, 0, len(arr)-1)
		assert.Equal(t, arr, arrSort)
	})
	t.Run("quick sort", func(t *testing.T) {
		ch := make(chan int)
		go quickSort(arr, ch)
		var a []int
		for value := range ch {
			a = append(a, value)
		}
		assert.Equal(t, a, arrSort)
	})
	t.Run("heap sort", func(t *testing.T) {
		HeapSort(arr)
		assert.Equal(t, arr, arrSort)
	})

	t.Run("counting sort", func(t *testing.T) {
		CountingSort(arr)
		assert.Equal(t, arr, arrSort)
	})

	t.Run("bucket sort", func(t *testing.T) {
		BucketSort(arr)
		assert.Equal(t, arr, arrSort)
	})
}
