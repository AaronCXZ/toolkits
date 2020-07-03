package sort

// 冒泡排序
func BubbleSort(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// 选择排序
func SelectionSort(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		min := i
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[min] {
				min = j
			}
		}
		arr[i], arr[min] = arr[min], arr[i]
	}
}

// 插入排序
func InsertionSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		preIndex, current := i-1, arr[i]
		for preIndex >= 0 && arr[preIndex] > current {
			arr[preIndex+1] = arr[preIndex]
			preIndex--
		}
		arr[preIndex+1] = current
	}
}

// 希尔排序
func ShellSort(arr []int) {
	for gap := len(arr) / 2; gap > 0; gap /= 2 {
		for i := 0; i < len(arr); i++ {
			j := i
			for j-gap >= 0 && arr[j] < arr[j-gap] {
				arr[j], arr[j-gap] = arr[j-gap], arr[j]
				j -= gap
			}
		}
	}
}

// 归并排序
func MergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	middle := len(arr) / 2
	left := MergeSort(arr[:middle])
	right := MergeSort(arr[middle:])
	return merge(left, right)
}

func merge(left, right []int) (result []int) {
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// 快速排序
func QuickSort(arr []int, left int, right int) {
	if left > right {
		return
	}
	index := partition(arr, left, right)
	QuickSort(arr, left, index-1)
	QuickSort(arr, index+1, right)
}

func partition(arr []int, left, right int) int {
	baseNum := arr[left]
	for left < right {
		for arr[right] >= baseNum && right > left {
			right--
		}
		arr[left] = arr[right]
		for arr[left] <= baseNum && right > left {
			left++
		}
		arr[right] = arr[left]
	}
	arr[right] = baseNum
	return right
}

// 多线程快速排序
func quickSort(arr []int, ch chan int) {
	if len(arr) == 1 {
		ch <- arr[0]
		close(ch)
		return
	}

	if len(arr) == 0 {
		close(ch)
		return
	}

	small := make([]int, 0)
	big := make([]int, 0)
	left := arr[0]
	arr = arr[1:]
	for _, num := range arr {
		switch {
		case num <= left:
			small = append(small, num)
		case num > left:
			big = append(big, num)
		}
	}
	left_ch := make(chan int, len(small))
	right_ch := make(chan int, len(big))
	go quickSort(small, left_ch)
	go quickSort(big, right_ch)

	for i := range left_ch {
		ch <- i
	}
	ch <- left
	for i := range right_ch {
		ch <- i
	}
	close(ch)
}

// 堆排序
func HeapSort(arr []int) {
	sift(arr, 0, len(arr)-1)
	for idx := len(arr) / 2; idx >= 0; idx-- {
		sift(arr, idx, len(arr)-1)
	}
	for idx := len(arr) - 1; idx >= 1; idx-- {
		arr[0], arr[idx] = arr[idx], arr[0]
		sift(arr, 0, idx-1)
	}
}

func sift(arr []int, left, right int) {
	fIdx := left
	sIdx := 2*fIdx + 1
	for sIdx <= right {
		if sIdx < right && arr[sIdx] < arr[sIdx+1] {
			sIdx++
		}
		if arr[fIdx] < arr[sIdx] {
			arr[fIdx], arr[sIdx] = arr[sIdx], arr[fIdx]
			fIdx = sIdx
			sIdx = 2*fIdx + 1
		} else {
			break
		}
	}
}

// 计数排序
func CountingSort(arr []int) {
	if len(arr) == 1 {
		return
	}
	min, max := countMaxMin(arr)
	temp := make([]int, max+1)
	for i := 0; i < len(arr); i++ {
		temp[arr[i]]++
	}
	var index int
	for i := min; i < len(temp); i++ {
		for j := temp[i]; j > 0; j-- {
			arr[index] = i
			index++
		}
	}
}

func countMaxMin(arr []int) (min, max int) {
	min, max = arr[0], arr[0]
	for i := 1; i < len(arr); i++ {
		if min > arr[i] {
			min = arr[i]
		}
		if max < arr[i] {
			max = arr[i]
		}
	}
	return min, max
}

// 桶排序
func BucketSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	_, max := countMaxMin(arr)
	buckets := make([][]int, len(arr))
	// 分配入桶
	for i := 0; i < len(arr); i++ {
		index := arr[i] * (len(arr) - 1) / max
		buckets[index] = append(buckets[index], arr[i])
	}
	// 桶内排序
	tmpPos := 0
	for i := 0; i < len(arr); i++ {
		if len(buckets[i]) > 0 {
			BucketSort(buckets[i])
			copy(arr[tmpPos:], buckets[i])
			tmpPos += len(buckets[i])
		}
	}
}

// 基数排序
func RadixSort(arr []int) {

}
