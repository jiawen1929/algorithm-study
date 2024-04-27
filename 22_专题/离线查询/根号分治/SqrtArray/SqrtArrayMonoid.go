// 分块数组/SqrtArray/SqrtArrayWithSum
//
// api:
//  1.Insert(index int32, v V)
//  2.Pop(index int32) V
//  3.Set(index int32, v V)
//  4.Get(index int32) V
//  5.Sum(start, end int32) V
//    SumAll() V
//  6.Clear()
//  7.Len() int32
//  8.GetAll() []V
//  9.ForEach(f func(i int32, v V) bool)

package main

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand"
	"time"
)

func main() {
	// demo()
	// test()
	testTime()
}

func demo() {
	bv := NewSqrtArrayMonoid(10, func(i int32) int32 { return 0 }, -1)
	for i := int32(0); i < 10; i++ {
		bv.Insert(i, 1)
	}
	bv.Set(3, 0)
	bv.Set(8, 1)
	bv.Insert(3, 1)
	bv.Pop(0)
	bv.Pop(0)
	bv.Pop(0)
	bv.Pop(0)
	fmt.Println(bv.GetAll())
}

type E = int32

func (*SqrtArrayMonoid) e() E        { return 0 }
func (*SqrtArrayMonoid) op(a, b E) E { return max32(a, b) }

// 使用分块+树状数组维护的动态数组.
type SqrtArrayMonoid struct {
	n                 int32
	blockSize         int32
	threshold         int32
	shouldRebuildTree bool
	blocks            [][]E
	blockSum          []E
	tree              []int32 // 每个块块长的前缀和
}

func NewSqrtArrayMonoid(n int32, f func(i int32) E, blockSize int32) *SqrtArrayMonoid {
	if blockSize == -1 {
		blockSize = int32(math.Sqrt(float64(n))) + 1
	}

	res := &SqrtArrayMonoid{n: n, blockSize: blockSize, threshold: blockSize << 1, shouldRebuildTree: true}
	blockCount := (n + blockSize - 1) / blockSize
	blocks, blockSum := make([][]E, blockCount), make([]E, blockCount)
	for bid := int32(0); bid < blockCount; bid++ {
		start, end := bid*blockSize, (bid+1)*blockSize
		if end > n {
			end = n
		}
		bucket := make([]E, end-start)
		sum := res.e()
		for i := start; i < end; i++ {
			bucket[i-start] = f(i)
			sum = res.op(sum, bucket[i-start])
		}
		blocks[bid], blockSum[bid] = bucket, sum
	}
	res.blocks, res.blockSum = blocks, blockSum
	return res
}

func (sl *SqrtArrayMonoid) Insert(index int32, value E) {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []E{value})
		sl.blockSum = append(sl.blockSum, value)
		sl.shouldRebuildTree = true
		sl.n++
		return
	}

	if index < 0 {
		index += sl.n
	}
	if index < 0 {
		index = 0
	}
	if index > sl.n {
		index = sl.n
	}

	pos, startIndex := sl._findKth(index)
	sl._updateTree(pos, true)
	sl.blockSum[pos] = sl.op(sl.blockSum[pos], value)
	sl.blocks[pos] = append(sl.blocks[pos], sl.e())
	copy(sl.blocks[pos][startIndex+1:], sl.blocks[pos][startIndex:])
	sl.blocks[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); n > sl.threshold {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][sl.blockSize:] // !注意max的设置(为了让左右互不影响)
		sl.blocks[pos] = sl.blocks[pos][:sl.blockSize:sl.blockSize]
		sl.blockSum = append(sl.blockSum, sl.e())
		copy(sl.blockSum[pos+2:], sl.blockSum[pos+1:])
		sl._updateSum(pos)
		sl._updateSum(pos + 1)
		sl.shouldRebuildTree = true
	}

	sl.n++
	return
}

func (sl *SqrtArrayMonoid) Pop(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.n--
	sl._updateTree(pos, false)

	copy(sl.blocks[pos][startIndex:], sl.blocks[pos][startIndex+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	sl._updateSum(pos)

	if len(sl.blocks[pos]) == 0 {
		// !delete block
		copy(sl.blocks[pos:], sl.blocks[pos+1:])
		sl.blocks = sl.blocks[:len(sl.blocks)-1]
		copy(sl.blockSum[pos:], sl.blockSum[pos+1:])
		sl.blockSum = sl.blockSum[:len(sl.blockSum)-1]
		sl.shouldRebuildTree = true
	}
	return value
}

func (sl *SqrtArrayMonoid) Get(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SqrtArrayMonoid) Set(index int32, value E) {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	oldValue := sl.blocks[pos][startIndex]
	if oldValue == value {
		return
	}
	sl.blocks[pos][startIndex] = value
	sl._updateSum(pos)
}

func (sl *SqrtArrayMonoid) Sum(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > sl.n {
		end = sl.n
	}
	if start >= end {
		return sl.e()
	}

	res := sl.e()
	pos, index := sl._findKth(start)
	count := end - start
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		bl := int32(len(block))
		endIndex := min32(bl, index+count)
		curCount := endIndex - index
		if curCount == bl {
			res = sl.op(res, sl.blockSum[pos])
		} else {
			for j := index; j < endIndex; j++ {
				res = sl.op(res, block[j])
			}
		}
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SqrtArrayMonoid) SumAll() E {
	res := sl.e()
	for _, v := range sl.blockSum {
		res = sl.op(res, v)
	}
	return res
}

func (sl *SqrtArrayMonoid) Len() int32 {
	return sl.n
}

func (sl *SqrtArrayMonoid) Clear() {
	sl.n = 0
	sl.shouldRebuildTree = true
	sl.blocks = sl.blocks[:0]
	sl.blockSum = sl.blockSum[:0]
	sl.tree = sl.tree[:0]
}

func (sl *SqrtArrayMonoid) GetAll() []E {
	res := make([]E, 0, sl.n)
	for _, block := range sl.blocks {
		res = append(res, block...)
	}
	return res
}

func (sl *SqrtArrayMonoid) ForEach(f func(i int32, v E) (shouldBreak bool)) {
	ptr := int32(0)
	for _, block := range sl.blocks {
		for _, v := range block {
			if f(ptr, v) {
				return
			}
			ptr++
		}
	}
}

func (sl *SqrtArrayMonoid) _rebuildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := int32(0); i < m; i++ {
		j := i | (i + 1)
		if j < m {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SqrtArrayMonoid) _updateTree(index int32, addOne bool) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	m := int32(len(tree))
	if addOne {
		for i := index; i < m; i |= i + 1 {
			tree[i]++
		}
	} else {
		for i := index; i < m; i |= i + 1 {
			tree[i]--
		}
	}
}

func (sl *SqrtArrayMonoid) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.n {
		return last, lastLen
	}
	if k >= sl.n-lastLen {
		return last, k + lastLen - sl.n
	}
	if sl.shouldRebuildTree {
		sl._rebuildTree()
	}
	tree := sl.tree
	pos = -1
	m := int32(len(tree))
	bitLen := int8(bits.Len32(uint32(m)))
	for d := bitLen - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *SqrtArrayMonoid) _updateSum(pos int32) {
	sum := sl.e()
	for _, v := range sl.blocks[pos] {
		sum = sl.op(sum, v)
	}
	sl.blockSum[pos] = sum
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

func test() {
	for i := int32(0); i < 100; i++ {
		n := rand.Int31n(10000) + 1000
		nums := make([]int, n)
		for i := int32(0); i < n; i++ {
			nums[i] = rand.Intn(100)
		}
		seg := NewSqrtArrayMonoid(n, func(i int32) E { return E(nums[i]) }, -1)

		for j := 0; j < 1000; j++ {
			// Get
			index := rand.Int31n(n)
			if seg.Get(index) != E(nums[index]) {
				fmt.Println("Get Error")
				panic("Get Error")
			}

			// Set
			index = rand.Int31n(n)
			value := rand.Intn(100)
			nums[index] = value
			seg.Set(index, E(value))
			if seg.Get(index) != E(value) {
				fmt.Println("Set Error")
				panic("Set Error")
			}

			// Query
			start, end := rand.Int31n(n), rand.Int31n(n)
			if start > end {
				start, end = end, start
			}
			sum_ := E(0)
			for i := start; i < end; i++ {
				sum_ = seg.op(sum_, E(nums[i]))
			}
			if seg.Sum(start, end) != sum_ {
				fmt.Println("Query Error")
				panic("Query Error")
			}

			// QueryAll
			sum_ = E(0)
			for _, v := range nums {
				sum_ = seg.op(sum_, E(v))
			}
			if seg.SumAll() != sum_ {
				fmt.Println("QueryAll Error")
				panic("QueryAll Error")
			}

			// GetAll
			all := seg.GetAll()
			for i, v := range all {
				if v != E(nums[i]) {
					fmt.Println("GetAll Error")
					panic("GetAll Error")
				}
			}

			// Insert
			index = rand.Int31n(n)
			value = rand.Intn(100)
			nums = append(nums, 0)
			copy(nums[index+1:], nums[index:])
			nums[index] = value
			seg.Insert(index, E(value))

			// Pop
			index = rand.Int31n(n)
			value = nums[index]
			nums = append(nums[:index], nums[index+1:]...)
			if seg.Pop(index) != E(value) {
				fmt.Println("Pop Error")
				panic("Pop Error")
			}

			// ForEach
			sum_ = E(0)
			seg.ForEach(func(i int32, v E) bool {
				sum_ = seg.op(sum_, v)
				return false
			})
			if sum_ != seg.SumAll() {
				fmt.Println("ForEach Error")
				panic("ForEach Error")
			}
		}
	}
	fmt.Println("Pass")
}

func testTime() {
	// 2e5
	n := int32(2e5)
	nums := make([]int, n)
	for i := 0; i < int(n); i++ {
		nums[i] = rand.Intn(5)
	}

	time1 := time.Now()
	seg := NewSqrtArrayMonoid(n, func(i int32) int32 { return E(nums[i]) }, -1)

	for i := int32(0); i < n; i++ {
		seg.Get(i)
		seg.Set(i, i)
		seg.Sum(i, n)
		seg.SumAll()
		seg.Insert(i, i)
		if i&1 == 0 {
			seg.Pop(i)
		}
		seg.SumAll()
	}
	fmt.Println("Time1", time.Since(time1)) // Time1 336.550792ms
}
