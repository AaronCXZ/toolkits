package singleNumber

func SingleNumber(nums []int) int {
	ans := 0
	for k := range nums {
		ans ^= nums[k]
	}
	return ans
}
