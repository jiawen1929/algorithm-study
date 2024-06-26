package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
)

func main() {
	CF371E()
	// test()
}

// Subway Innovation
// https://www.luogu.com.cn/problem/CF371E
// 在直线上给定 n 个点，定义其距离为横坐标的绝对值，请你保留 k 个点，使这些点两两之间的距离和最小，输出这 k 个点的坐标。
func CF371E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	var k int32
	fmt.Fscan(in, &k)

	points := make([][2]int32, 0, k)
	for i := int32(0); i < n; i++ {
		points = append(points, [2]int32{nums[i], i})
	}
	sort.Slice(points, func(i, j int) bool { return points[i][0] < points[j][0] })
	sortedNums := make([]int, n)
	for i := int32(0); i < n; i++ {
		sortedNums[i] = int(points[i][0])
	}

	D := NewMedianFinderSortedNums(sortedNums)
	curSum := DistSumOfAllPairsRange(sortedNums, 0, k)
	bestLeft, bestSum := int32(0), curSum
	for left := int32(1); left+k-1 < n; left++ {
		right := left + k - 1
		// reove left-1, add right
		curSum -= D.DistSumRange(sortedNums[left-1], left-1, right)
		curSum += D.DistSumRange(sortedNums[right], left, right+1)
		if curSum < bestSum {
			bestLeft = left
			bestSum = curSum
		}
	}

	for i := bestLeft; i < bestLeft+k; i++ {
		fmt.Fprint(out, points[i][1]+1, " ")
	}
}

type MedianFinderSortedNums struct {
	sortedNums []int
	preSum     []int
}

func NewMedianFinderSortedNums(sortedNums []int) *MedianFinderSortedNums {
	preSum := make([]int, len(sortedNums)+1)
	for i := 0; i < len(sortedNums); i++ {
		preSum[i+1] = preSum[i] + sortedNums[i]
	}
	return &MedianFinderSortedNums{sortedNums: sortedNums, preSum: preSum}
}

func (mf *MedianFinderSortedNums) Median() int {
	return mf.MedianRange(0, int32(len(mf.sortedNums)))
}

// 返回区间 [start, end) 中的中位数(如果有两个中位数，返回较小的那个).
func (mf *MedianFinderSortedNums) MedianRange(start, end int32) int {
	if start < 0 {
		start = 0
	}
	if n := int32(len(mf.sortedNums)); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	mid := start + (end-start-1)>>1
	return mf.sortedNums[mid]
}

// 有序数组中所有点到`x=k`的距离之和.
func (mf *MedianFinderSortedNums) DistSum(k int) int {
	return mf.DistSumRange(k, 0, int32(len(mf.sortedNums)))
}

// 有序数组切片[start:end)中所有点到`x=k`的距禿之和.
func (mf *MedianFinderSortedNums) DistSumRange(k int, start, end int32) int {
	if start < 0 {
		start = 0
	}
	if n := int32(len(mf.sortedNums)); end > n {
		end = n
	}
	if start >= end {
		return 0
	}
	pos := int32(sort.SearchInts(mf.sortedNums, k))
	if pos <= start {
		return (mf.preSum[end] - mf.preSum[start]) - k*(int(end-start))
	}
	if pos >= end {
		return k*(int(end-start)) - (mf.preSum[end] - mf.preSum[start])
	}
	leftSum := k*int(pos-start) - (mf.preSum[pos] - mf.preSum[start])
	rightSum := mf.preSum[end] - mf.preSum[pos] - k*(int(end-pos))
	return leftSum + rightSum
}

func (mf *MedianFinderSortedNums) DistSumToMedian() int {
	median := mf.Median()
	return mf.DistSum(median)
}

func (mf *MedianFinderSortedNums) DistSumToMedianRange(start, end int32) int {
	median := mf.MedianRange(start, end)
	return mf.DistSumRange(median, start, end)
}

// 有序数组中所有点对两两距离之和.一共有`n*(n-1)//2`对点对.
func DistSumOfAllPairs(sortedNums []int) int {
	n := len(sortedNums)
	res := 0
	preSum := 0
	for i := 0; i < n; i++ {
		res += sortedNums[i]*i - preSum
		preSum += sortedNums[i]
	}
	return res
}

func DistSumOfAllPairsRange(sortedNums []int, start, end int32) int {
	res := 0
	preSum := 0
	for i := start; i < end; i++ {
		res += sortedNums[i]*int(i) - preSum
		preSum += sortedNums[i]
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func test() {

	medianBruteForce := func(nums []int) int {
		if len(nums) == 0 {
			return 0
		}
		sortedNums := append([]int(nil), nums...)
		sort.Ints(sortedNums)
		mid := (len(sortedNums) - 1) / 2
		return sortedNums[mid]
	}

	distSumBruteForce := func(nums []int, to int) int {
		res := 0
		for _, num := range nums {
			res += abs(num - to)
		}
		return res
	}

	distSumRangeBruteForce := func(nums []int, to, start, end int) int {
		res := 0
		for i := start; i < end; i++ {
			res += abs(nums[i] - to)
		}
		return res
	}

	distSumToMedianBruteForce := func(nums []int) int {
		median := medianBruteForce(nums)
		return distSumBruteForce(nums, median)
	}

	distSumToMedianRangeBruteForce := func(nums []int, start, end int) int {
		median := medianBruteForce(nums[start:end])
		return distSumRangeBruteForce(nums, median, start, end)
	}

	for tc := 0; tc < 1000; tc++ {
		n := rand.Intn(1000) + 1
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = rand.Intn(1e6) - 5e5
		}
		sort.Ints(nums)

		mf := NewMedianFinderSortedNums(nums)

		for i := 0; i < 100; i++ {
			start, end := rand.Intn(n), rand.Intn(n)+1
			if start > end {
				start, end = end, start
			}
			to := rand.Intn(1e5) - 5e4

			if mf.MedianRange(int32(start), int32(end)) != medianBruteForce(nums[start:end]) {
				panic("err0")
			}
			if mf.DistSum(int(to)) != distSumBruteForce(nums, to) {
				panic("err1")
			}
			if mf.DistSumRange(int(to), int32(start), int32(end)) != distSumRangeBruteForce(nums, to, start, end) {
				panic("err2")
			}
			if mf.DistSumToMedian() != distSumToMedianBruteForce(nums) {
				panic("err3")
			}
			if mf.DistSumToMedianRange(int32(start), int32(end)) != distSumToMedianRangeBruteForce(nums, start, end) {
				fmt.Println(mf.DistSumToMedianRange(int32(start), int32(end)), distSumToMedianRangeBruteForce(nums, start, end))
				panic("err4")
			}
		}
	}

	fmt.Println("pass")
}
