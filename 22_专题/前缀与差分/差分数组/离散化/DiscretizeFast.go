package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	newNums, originValue, compress := DiscretizeFast(nums)
	fmt.Println(newNums, originValue)
	for i := range nums {
		if originValue[newNums[i]] != nums[i] {
			panic("error")
		}
	}

	fmt.Println(compress(1), compress(2), compress(3), compress(4), compress(5), compress(6), compress(7), compress(8), compress(9), compress(10))
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func DiscretizeFast(nums []int) (newNums []int32, origin []int, f func(v int) int32) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	f = func(v int) int32 { return int32(sort.SearchInts(origin, v)) }
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}
