package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	bit := NewBITRangeBlock(e, op)
	bit.Build(3, func(i int32) int { return 0 })
	fmt.Println(bit)
	bit.Update(0, 1)
	fmt.Println(bit)
	bit.Update(2, 2)
	fmt.Println(bit, bit.QueryRange(0, 2))
	bit.Update(2, 4)
	fmt.Println(bit, bit.QueryRange(0, 3))
}

// TLE
// https://leetcode.cn/problems/maximize-the-minimum-powered-city/description/
func maxPower(stations []int, r int, k int) int64 {
	n := len(stations)
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	bit := NewBITRangeBlock(e, op)
	check := func(mid int) bool {
		bit.Build(int32(len(stations)), func(i int32) int { return stations[i] })
		curK := k
		for i := 0; i < n; i++ {
			cur := bit.QueryRange(int32(max(0, i-r)), int32(min(i+r+1, n)))
			if cur < mid {
				diff := mid - cur
				bit.Update(min32(int32(i+r), int32(n-1)), diff)
				curK -= diff
				if curK < 0 {
					return false
				}
			}
		}
		return true
	}

	left := 1
	right := int(2e15)
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return int64(right)
}

// 基于分块实现的`树状数组`.
// `O(1)`单点加，`O(sqrt(n))`区间和查询.
// 一般配合莫队算法使用.
type BITRangeBlock[E any] struct {
	_n          int32
	_belong     []int32
	_blockStart []int32
	_blockEnd   []int32
	_nums       []E
	_blockSum   []E
	e           func() E
	op          func(a, b E) E
}

func NewBITRangeBlock[E any](e func() E, op func(a, b E) E) *BITRangeBlock[E] {
	return &BITRangeBlock[E]{e: e, op: op}
}

func (b *BITRangeBlock[E]) Build(n int32, f func(i int32) E) {
	blockSize := int32(math.Sqrt(float64(n)) + 1)
	blockCount := 1 + (n / blockSize)
	belong := make([]int32, n)
	for i := range belong {
		belong[i] = int32(i) / blockSize
	}
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	for i := range blockStart {
		blockStart[i] = int32(i) * blockSize
		tmp := (int32(i) + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	nums := make([]E, n)
	for i := range nums {
		nums[i] = b.e()
	}
	blockSum := make([]E, blockCount)
	for i := range blockSum {
		blockSum[i] = b.e()
	}
	b._n = n
	b._belong = belong
	b._blockStart = blockStart
	b._blockEnd = blockEnd
	b._nums = nums
	b._blockSum = blockSum
	for i := int32(0); i < n; i++ {
		b.Update(i, f(i))
	}
}

func (b *BITRangeBlock[E]) Update(index int32, delta E) {
	b._nums[index] = b.op(b._nums[index], delta)
	bid := b._belong[index]
	b._blockSum[bid] = b.op(b._blockSum[bid], delta)
}

func (b *BITRangeBlock[E]) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > b._n {
		end = b._n
	}
	if start >= end {
		return b.e()
	}
	res := b.e()
	bid1 := b._belong[start]
	bid2 := b._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			res = b.op(res, b._nums[i])
		}
		return res
	}
	for i := start; i < b._blockEnd[bid1]; i++ {
		res = b.op(res, b._nums[i])
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res = b.op(res, b._blockSum[bid])
	}
	for i := b._blockStart[bid2]; i < end; i++ {
		res = b.op(res, b._nums[i])
	}
	return res
}

func (b *BITRangeBlock[E]) String() string {
	sb := strings.Builder{}
	sb.WriteString("BITRangeBlock{")
	for i := 0; i < int(b._n); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", b.QueryRange(int32(i), int32(i+1))))
	}
	sb.WriteString("}")
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
