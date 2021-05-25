package holmes

type ring struct {
	data   []int
	idx    int
	sum    int
	maxLen int
}

func newRing(maxLen int) ring {
	return ring{
		data:   make([]int, 0, maxLen),
		idx:    0,
		maxLen: maxLen,
	}
}

func (r *ring) Push(i int) {
	if r.maxLen == 0 {
		return
	}
	if r.idx >= r.maxLen {
		r.idx = 0
	}

	if len(r.data) < r.maxLen {
		r.sum += i
		r.data = append(r.data, i)
		return
	}

	r.sum += i - r.data[r.idx]

	r.data[r.idx] = i
	r.idx++
}

func (r *ring) avg() int {
	if r.maxLen == 0 || len(r.data) == 0 {
		return 0
	}
	return r.sum / len(r.data)
}
