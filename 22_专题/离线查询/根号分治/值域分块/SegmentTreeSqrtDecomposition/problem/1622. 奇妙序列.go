// api:
//  1.Insert(index int32, v V) -> log(n)
//  2.Pop(index int32) V -> sqrt(n)
//  3.Set(index int32, v V) -> sqrt(n)
//  4.Get(index int32) V -> O(log(sqrt(n)))
//  5.Query(start, end int32) V -> O(sqrt(n))
//    QueryAll() V
//  6.Update(start, end int32, lazy Id) -> O(sqrt(n))
//    UpdateAll(lazy Id)
//  7.Clear()
//  8.Len() int32
//  9.GetAll() []V
// 10.ForEach(f func(i int32, v V) bool)
// 11.MaxRight(start int32, predicate func(end int32, sum E) bool) int32
// 12.MinLeft(end int32, predicate func(start int32, sum E) bool) int32

package main

import (
	"math"
	"math/bits"
)

// https://leetcode.cn/problems/fancy-sequence/
type Fancy struct {
	seg *LazySegmentTreeSqrtDecompositionDynamic
}

func Constructor() Fancy {
	return Fancy{seg: NewLazySegmentTreeSqrtDecompositionDynamic(0, nil, 200)}
}

func (this *Fancy) Append(val int) {
	this.seg.Insert(this.seg.Len(), FromElement(val))
}

func (this *Fancy) AddAll(inc int) {
	this.seg.UpdateAll(Id{add: inc, mul: 1})
}

func (this *Fancy) MultAll(m int) {
	this.seg.UpdateAll(Id{add: 0, mul: m})
}

func (this *Fancy) GetIndex(idx int) int {
	if idx >= int(this.seg.Len()) {
		return -1
	}
	return this.seg.Get(int32(idx)).sum
}

/**
 * Your Fancy object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Append(val);
 * obj.AddAll(inc);
 * obj.MultAll(m);
 * param_4 := obj.GetIndex(idx);
 */

const MOD int = 1e9 + 7

type E = struct{ sum, size int }
type Id = struct{ mul, add int }

func FromElement(v int) E { return E{sum: v, size: 1} }

// RangeAddRangeMax

func (*LazySegmentTreeSqrtDecompositionDynamic) e() E   { return E{} }
func (*LazySegmentTreeSqrtDecompositionDynamic) id() Id { return Id{mul: 1} }
func (*LazySegmentTreeSqrtDecompositionDynamic) op(a, b E) E {
	a.sum += b.sum
	if a.sum >= MOD {
		a.sum -= MOD
	}
	a.size += b.size
	return a
}
func (*LazySegmentTreeSqrtDecompositionDynamic) mapping(f Id, a E) E {
	a.sum = (a.sum*f.mul + a.size*f.add) % MOD
	return a
}
func (*LazySegmentTreeSqrtDecompositionDynamic) composition(f, g Id) Id {
	return Id{
		mul: (f.mul * g.mul) % MOD,
		add: (f.mul*g.add + f.add) % MOD,
	}
}

type LazySegmentTreeSqrtDecompositionDynamic struct {
	n                 int32
	blockSize         int32
	threshold         int32
	shouldRebuildTree bool
	data              [][]E
	sum               []E
	lazy              []Id
	tree              []int32
}

// bucketSize 为 -1 时，使用默认值 sqrt(n).
func NewLazySegmentTreeSqrtDecompositionDynamic(n int32, f func(i int32) E, blockSize int32) *LazySegmentTreeSqrtDecompositionDynamic {
	if blockSize == -1 {
		blockSize = int32(math.Sqrt(float64(n))) + 1
	}
	if blockSize < 100 {
		blockSize = 100 // 防止 blockSize 过小
	}
	bucketCount := (n + blockSize - 1) / blockSize
	res := &LazySegmentTreeSqrtDecompositionDynamic{n: n, blockSize: blockSize, threshold: blockSize << 1, shouldRebuildTree: true}
	blocks, blockSum := make([][]E, bucketCount), make([]E, bucketCount)
	blockLazy := make([]Id, bucketCount)
	for bid := int32(0); bid < bucketCount; bid++ {
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
		blockLazy[bid] = res.id()
	}
	res.data, res.sum, res.lazy = blocks, blockSum, blockLazy
	return res
}

func (sl *LazySegmentTreeSqrtDecompositionDynamic) Insert(index int32, value E) {
	if len(sl.data) == 0 {
		sl.data = append(sl.data, []E{value})
		sl.sum = append(sl.sum, value)
		sl.lazy = append(sl.lazy, sl.id())
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
	sl._propagate(pos)
	sl.sum[pos] = sl.op(sl.sum[pos], value)
	sl.data[pos] = append(sl.data[pos], sl.e())
	copy(sl.data[pos][startIndex+1:], sl.data[pos][startIndex:])
	sl.data[pos][startIndex] = value

	// n -> load + (n - load)
	if n := int32(len(sl.data[pos])); n > sl.threshold {
		sl.data = append(sl.data, nil)
		copy(sl.data[pos+2:], sl.data[pos+1:])
		sl.data[pos+1] = sl.data[pos][sl.blockSize:]
		sl.data[pos] = sl.data[pos][:sl.blockSize:sl.blockSize] // !注意max的设置(为了让左右互不影响)

		sl.sum = append(sl.sum, sl.e())
		copy(sl.sum[pos+2:], sl.sum[pos+1:])
		sl._updateSum(pos)
		sl._updateSum(pos + 1)

		sl.lazy = append(sl.lazy, sl.id())
		copy(sl.lazy[pos+2:], sl.lazy[pos+1:])
		// sl.lazy[pos] = sl.id()
		sl.lazy[pos+1] = sl.id()

		sl.shouldRebuildTree = true
	}

	sl.n++
	return
}

func (sl *LazySegmentTreeSqrtDecompositionDynamic) Pop(index int32) E {
	if index < 0 {
		index += sl.n
	}
	pos, startIndex := sl._findKth(index)
	value := sl.data[pos][startIndex]
	if sl.lazy[pos] != sl.id() {
		value = sl.mapping(sl.lazy[pos], value)
	}
	// !delete element
	sl.n--
	sl._updateTree(pos, false)

	copy(sl.data[pos][startIndex:], sl.data[pos][startIndex+1:])
	sl.data[pos] = sl.data[pos][:len(sl.data[pos])-1]
	sl._propagateAndUpdateSum(pos)

	if len(sl.data[pos]) == 0 {
		// !delete block
		copy(sl.data[pos:], sl.data[pos+1:])
		sl.data = sl.data[:len(sl.data)-1]
		copy(sl.sum[pos:], sl.sum[pos+1:])
		sl.sum = sl.sum[:len(sl.sum)-1]
		copy(sl.lazy[pos:], sl.lazy[pos+1:])
		sl.lazy = sl.lazy[:len(sl.lazy)-1]
		sl.shouldRebuildTree = true
	}
	return value
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) Set(index int32, value E) {
	if index < 0 {
		index += st.n
	}
	bid, pos := st._findKth(index)
	st._propagate(bid)
	if st.data[bid][pos] == value {
		return
	}
	st.data[bid][pos] = value
	st._updateSum(bid)
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) Get(index int32) E {
	if index < 0 {
		index += st.n
	}
	bid, pos := st._findKth(index)
	if st.lazy[bid] == st.id() {
		return st.data[bid][pos]
	} else {
		return st.mapping(st.lazy[bid], st.data[bid][pos])
	}
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}

	bid1, startIndex1 := st._findKth(start)
	bid2, startIndex2 := st._findKth(end)
	start, end = startIndex1, startIndex2
	res := st.e()
	if bid1 == bid2 {
		block := st.data[bid1]
		for i := start; i < end; i++ {
			res = st.op(res, block[i])
		}
		if st.lazy[bid1] != st.id() {
			res = st.mapping(st.lazy[bid1], res)
		}
	} else {
		if start < int32(len(st.data[bid1])) {
			block1 := st.data[bid1]
			for i := start; i < int32(len(block1)); i++ {
				res = st.op(res, block1[i])
			}
			if st.lazy[bid1] != st.id() {
				res = st.mapping(st.lazy[bid1], res)
			}
		}

		for i := bid1 + 1; i < bid2; i++ {
			res = st.op(res, st.sum[i])
		}

		if m := int32(len(st.data)); bid2 < m && end > 0 {
			block2 := st.data[bid2]
			tmp := st.e()
			for i := int32(0); i < end; i++ {
				tmp = st.op(tmp, block2[i])
			}
			if st.lazy[bid2] != st.id() {
				tmp = st.mapping(st.lazy[bid2], tmp)
			}
			res = st.op(res, tmp)
		}
	}
	return res
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) QueryAll() E {
	res := st.e()
	for bid, m := int32(0), int32(len(st.data)); bid < m; bid++ {
		res = st.op(res, st.sum[bid])
	}
	return res
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) Update(start, end int32, lazy Id) {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return
	}

	changeData := func(k, l, r int32) {
		st._propagate(k)
		d := st.data[k]
		for i := l; i < r; i++ {
			d[i] = st.mapping(lazy, d[i])
		}
		e := st.e()
		for _, v := range d {
			e = st.op(e, v)
		}
		st.sum[k] = e
	}

	bid1, startIndex1 := st._findKth(start)
	bid2, startIndex2 := st._findKth(end)
	start, end = startIndex1, startIndex2
	m := int32(len(st.data))
	if bid1 == bid2 {
		if bid1 < m {
			changeData(bid1, start, end)
		}
	} else {
		if bid1 < m {
			if start == 0 {
				if st.lazy[bid1] == st.id() {
					st.lazy[bid1] = lazy
				} else {
					st.lazy[bid1] = st.composition(lazy, st.lazy[bid1])
				}
				st.sum[bid1] = st.mapping(lazy, st.sum[bid1])
			} else {
				changeData(bid1, start, int32(len(st.data[bid1])))
			}
		}

		for i := bid1 + 1; i < bid2; i++ {
			if st.lazy[i] == st.id() {
				st.lazy[i] = lazy
			} else {
				st.lazy[i] = st.composition(lazy, st.lazy[i])
			}
			st.sum[i] = st.mapping(lazy, st.sum[i])
		}

		if bid2 < m {
			if end == int32(len(st.data[bid2])) {
				if st.lazy[bid2] == st.id() {
					st.lazy[bid2] = lazy
				} else {
					st.lazy[bid2] = st.composition(lazy, st.lazy[bid2])
				}
				st.sum[bid2] = st.mapping(lazy, st.sum[bid2])
			} else {
				changeData(bid2, 0, end)
			}
		}
	}
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) UpdateAll(lazy Id) {
	for i := int32(0); i < int32(len(st.data)); i++ {
		if st.lazy[i] == st.id() {
			st.lazy[i] = lazy
		} else {
			st.lazy[i] = st.composition(lazy, st.lazy[i])
		}
		st.sum[i] = st.mapping(lazy, st.sum[i])
	}
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) GetAll() []E {
	st._propagateAll()
	res := make([]E, 0, st.n)
	for _, bucket := range st.data {
		res = append(res, bucket...)
	}
	return res
}

func (st LazySegmentTreeSqrtDecompositionDynamic) ForEach(f func(i int32, v E) (shouldBreak bool)) {
	ptr := int32(0)
	for bid, block := range st.data {
		st._propagate(int32(bid))
		for _, v := range block {
			if f(ptr, v) {
				return
			}
			ptr++
		}
	}
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) Len() int32 {
	return st.n
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) Clear() {
	st.n = 0
	st.shouldRebuildTree = true
	st.data = st.data[:0]
	st.sum = st.sum[:0]
	st.lazy = st.lazy[:0]
	st.tree = st.tree[:0]
}

// 查询最大的 end 使得切片 [start:end] 内的值满足 predicate.
func (st *LazySegmentTreeSqrtDecompositionDynamic) MaxRight(start int32, predicate func(end int32, sum E) bool) int32 {
	if start >= st.n {
		return st.n
	}

	curSum := st.e()
	res := start
	bid, startPos := st._findKth(start)

	// 散块内
	{
		st._propagate(bid)
		pos := startPos
		block := st.data[bid]
		m := int32(len(block))
		for ; pos < m; pos++ {
			nextRes, nextSum := res+1, st.op(curSum, block[pos])
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				return res
			}
		}
	}
	bid++

	// 整块跳跃
	{
		m := int32(len(st.data))
		for ; bid < m; bid++ {
			nextRes := res + int32(len(st.data[bid]))
			nextSum := st.op(curSum, st.sum[bid])
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				// 答案在这个块内
				st._propagate(bid)
				block := st.data[bid]
				for _, v := range block {
					nextRes, nextSum = res+1, st.op(curSum, v)
					if predicate(nextRes, nextSum) {
						res, curSum = nextRes, nextSum
					} else {
						return res
					}
				}
			}
		}
	}

	return res
}

// 查询最小的 start 使得切片 [start:end] 内的值满足 predicate.
func (st *LazySegmentTreeSqrtDecompositionDynamic) MinLeft(end int32, predicate func(start int32, sum E) bool) int32 {
	if end <= 0 {
		return 0
	}

	curSum := st.e()
	res := end
	bid, startPos := st._findKth(end - 1)

	// 散块内
	{
		st._propagate(bid)
		pos := startPos
		block := st.data[bid]
		for ; pos >= 0; pos-- {
			nextRes, nextSum := res-1, st.op(block[pos], curSum)
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				return res
			}
		}
	}
	bid--

	// 整块跳跃
	{
		for ; bid >= 0; bid-- {
			nextRes := res - int32(len(st.data[bid]))
			nextSum := st.op(st.sum[bid], curSum)
			if predicate(nextRes, nextSum) {
				res, curSum = nextRes, nextSum
			} else {
				// 答案在这个块内
				st._propagate(bid)
				block := st.data[bid]
				for i := int32(len(block)) - 1; i >= 0; i-- {
					nextRes, nextSum = res-1, st.op(block[i], curSum)
					if predicate(nextRes, nextSum) {
						res, curSum = nextRes, nextSum
					} else {
						return res
					}
				}
			}
		}
	}
	return res
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) _propagate(k int32) {
	if st.lazy[k] == st.id() {
		return
	}
	f := st.lazy[k]
	dk := st.data[k]
	for i := int32(0); i < int32(len(dk)); i++ {
		dk[i] = st.mapping(f, dk[i])
	}
	st.lazy[k] = st.id()
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) _propagateAndUpdateSum(k int32) {
	f := st.lazy[k]
	dk := st.data[k]
	e := st.e()
	if f == st.id() {
		for _, v := range dk {
			e = st.op(e, v)
		}
	} else {
		for i, v := range dk {
			dk[i] = st.mapping(f, v)
			e = st.op(e, dk[i])
		}
	}
	st.sum[k] = e
	st.lazy[k] = st.id()
}

func (st *LazySegmentTreeSqrtDecompositionDynamic) _propagateAll() {
	m := int32(len(st.data))
	for k := int32(0); k < m; k++ {
		st._propagate(k)
	}
	for i := int32(0); i < m; i++ {
		st.lazy[i] = st.id()
	}
}

func (sl *LazySegmentTreeSqrtDecompositionDynamic) _updateSum(pos int32) {
	sum := sl.e()
	for _, v := range sl.data[pos] {
		sum = sl.op(sum, v)
	}
	sl.sum[pos] = sum
}

func (sl *LazySegmentTreeSqrtDecompositionDynamic) _rebuildTree() {
	m := int32(len(sl.data))
	sl.tree = make([]int32, m)
	for i := int32(0); i < m; i++ {
		sl.tree[i] = int32(len(sl.data[i]))
	}
	tree := sl.tree
	for i := int32(0); i < m; i++ {
		j := i | (i + 1)
		if j < m {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *LazySegmentTreeSqrtDecompositionDynamic) _updateTree(index int32, addOne bool) {
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

func (sl *LazySegmentTreeSqrtDecompositionDynamic) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.data[0])) {
		return 0, k
	}
	last := int32(len(sl.data) - 1)
	lastLen := int32(len(sl.data[last]))
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
