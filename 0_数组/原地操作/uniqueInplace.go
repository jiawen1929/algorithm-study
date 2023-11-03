package main

import "fmt"

// 有序数组原地去重.
func UniqueInplace(sorted *([]int)) {
	nums := *sorted
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	*sorted = nums[:slow+1]
}

func main() {
	nums := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	UniqueInplace(&nums)
	fmt.Println(nums)
}
