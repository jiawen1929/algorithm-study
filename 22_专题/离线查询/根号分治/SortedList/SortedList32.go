package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	abc241_d()
}

const INF int = 1e18

// https://atcoder.jp/contests/abc241/tasks/abc241_d
func abc241_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	sl := NewSortedList32(func(a, b S) bool { return a < b })
	for i := 0; i < q; i++ {
		var t, x int
		fmt.Fscan(in, &t, &x)
		if t == 1 {
			sl.Add(x)
			continue
		}

		// A の
		// x 以下の要素のうち、大きい方から
		// k 番目の値を出力する。(
		// k は
		// 5 以下)
		var k int32
		fmt.Fscan(in, &k)
		if t == 2 {
			pos := sl.BisectRight(x) - k
			if pos < 0 {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, sl.At(pos))
			}
		}

		// A の
		// x 以上の要素のうち、小さい方から
		// k 番目の値を出力する。(
		// k は
		// 5 以下)
		if t == 3 {
			pos := sl.BisectLeft(x) + k - 1
			if pos >= sl.Len() {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, sl.At(pos))
			}
		}
	}
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int32 = 200

type S = int

var EMPTY S

// 使用分块+树状数组维护的有序序列.
type SortedList32 struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSortedList32(less func(a, b S) bool, elements ...S) *SortedList32 {
	elements = append(elements[:0:0], elements...)
	res := &SortedList32{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := int32(len(elements))
	blocks := [][]S{}
	for start := int32(0); start < n; start += _LOAD {
		end := min32(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SortedList32) Add(value S) *SortedList32 {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:]
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		sl.mins = append(sl.mins, EMPTY)
		copy(sl.mins[pos+2:], sl.mins[pos+1:])
		sl.mins[pos+1] = sl.blocks[pos+1][0]
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SortedList32) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < int32(len(sl.blocks[pos])) && sl.blocks[pos][index] == value
}

func (sl *SortedList32) Discard(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locRight(value)
	if index > 0 && sl.blocks[pos][index-1] == value {
		sl._delete(pos, index-1)
		return true
	}
	return false
}

func (sl *SortedList32) Pop(index int32) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

func (sl *SortedList32) At(index int32) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedList32) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedList32) Lower(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList32) Higher(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedList32) Floor(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList32) Ceiling(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedList32) BisectLeft(value S) int32 {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedList32) BisectRight(value S) int32 {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedList32) Count(value S) int32 {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedList32) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedList32) ForEach(f func(value S, index int32) bool, reverse bool) {
	if !reverse {
		count := int32(0)
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				if f(block[j], count) {
					return
				}
				count++
			}
		}
		return
	}
	count := int32(0)
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			if f(block[j], count) {
				return
			}
			count++
		}
	}
}

func (sl *SortedList32) Enumerate(start, end int32, f func(value S), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	pos, startIndex := sl._findKth(start)
	count := end - start
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endIndex := min32(int32(len(block)), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == int32(len(block)) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = append(sl.blocks[pos][:startIndex], sl.blocks[pos][endIndex:]...)
				sl.mins[pos] = sl.blocks[pos][0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SortedList32) Slice(start, end int32) []S {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]S, 0, count)
	pos, index := sl._findKth(start)
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endPos := min32(int32(len(block)), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedList32) Range(min, max S) []S {
	if sl.less(max, min) {
		return nil
	}
	res := []S{}
	pos := sl._locBlock(min)
	m := int32(len(sl.blocks))
	for i := pos; i < m; i++ {
		block := sl.blocks[i]
		for j := 0; j < len(block); j++ {
			x := block[j]
			if sl.less(max, x) {
				return res
			}
			if !sl.less(x, min) {
				res = append(res, x)
			}
		}
	}
	return res
}

func (sl *SortedList32) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedList32) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedList32) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value S, index int32) bool {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
		return false
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedList32) Len() int32 {
	return sl.size
}

func (sl *SortedList32) _delete(pos, index int32) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	copy(sl.blocks[pos][index:], sl.blocks[pos][index+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	copy(sl.blocks[pos:], sl.blocks[pos+1:])
	sl.blocks = sl.blocks[:len(sl.blocks)-1]
	copy(sl.mins[pos:], sl.mins[pos+1:])
	sl.mins = sl.mins[:len(sl.mins)-1]
	sl.shouldRebuildTree = true
}

func (sl *SortedList32) _locLeft(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(cur[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedList32) _locRight(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedList32) _locBlock(value S) int32 {
	left, right := int32(-1), int32(len(sl.blocks)-1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	return right
}

func (sl *SortedList32) _buildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	for i := 0; i < len(tree); i++ {
		j := i | (i + 1)
		if j < len(tree) {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SortedList32) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := index; i < m; i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedList32) _queryTree(end int32) int32 {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SortedList32) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < int32(len(tree)) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func (sl *SortedList32) _iteratorAt(pos, index int32) *Iterator {
	return &Iterator{sl: sl, pos: pos, index: index}
}

func (sl *SortedList32) IteratorAt(index int32) *Iterator {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("Index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl._iteratorAt(pos, startIndex)
}

// 返回一个迭代器，指向键值>= value的第一个元素.
func (sl *SortedList32) LowerBound(value S) *Iterator {
	pos, index := sl._locLeft(value)
	return sl._iteratorAt(pos, index)
}

// 返回一个迭代器，指向键值> value的第一个元素.
func (sl *SortedList32) UpperBound(value S) *Iterator {
	pos, index := sl._locRight(value)
	return sl._iteratorAt(pos, index)
}

type Iterator struct {
	sl    *SortedList32
	pos   int32
	index int32
}

func (it *Iterator) HasNext() bool {
	b := it.sl.blocks
	m := int32(len(b))
	if it.pos < m-1 {
		return true
	}
	return it.pos == m-1 && it.index < int32(len(b[it.pos]))-1
}

func (it *Iterator) Next() {
	it.index++
	if it.index == int32(len(it.sl.blocks[it.pos])) {
		it.pos++
		it.index = 0
	}
}

func (it *Iterator) HasPrev() bool {
	if it.pos > 0 {
		return true
	}
	return it.pos == 0 && it.index > 0
}

func (it *Iterator) Prev() {
	it.index--
	if it.index == -1 {
		it.pos--
		it.index = int32(len(it.sl.blocks[it.pos]) - 1)
	}
}

func (it *Iterator) IsBegin() bool {
	return it.pos == 0 && it.index == 0
}

func (it *Iterator) IsEnd() bool {
	m := int32(len(it.sl.blocks))
	return it.pos == m && it.index == 0
}

// GetMut
func (it *Iterator) Value() S {
	return it.sl.blocks[it.pos][it.index]
}

func (it *Iterator) NextValue() S {
	newPos, newIndex := it.pos, it.index
	newIndex++
	if newIndex == int32(len(it.sl.blocks[it.pos])) {
		newPos++
		newIndex = 0
	}
	return it.sl.blocks[newPos][newIndex]
}

func (it *Iterator) PrevValue() S {
	newPos, newIndex := it.pos, it.index
	newIndex--
	if newIndex == -1 {
		newPos--
		newIndex = int32(len(it.sl.blocks[newPos]) - 1)
	}
	return it.sl.blocks[newPos][newIndex]
}

func (it *Iterator) ToIndex() int32 {
	res := it.sl._queryTree(it.pos)
	return res + it.index
}

func (it *Iterator) Copy() *Iterator {
	return &Iterator{sl: it.sl, pos: it.pos, index: it.index}
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

// hack.如果已知元素是最大/最小的, 可以使用下面的方法.
func (sl *SortedList32) _appendFirst(value S) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return
	}
	pos := int32(0)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos], EMPTY)
	copy(sl.blocks[pos][1:], sl.blocks[pos])
	sl.blocks[pos][0] = value
	sl._adjust(pos)
	return
}
func (sl *SortedList32) _appendLast(value S) {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return
	}
	pos := int32(len(sl.blocks) - 1)
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos], value)
	sl._adjust(pos)
	return
}
func (sl *SortedList32) _popFirst() S {
	pos, startIndex := int32(0), int32(0)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}
func (sl *SortedList32) _popLast() S {
	pos := int32(len(sl.blocks) - 1)
	startIndex := len(sl.blocks[pos]) - 1
	value := sl.blocks[pos][startIndex]
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		return value
	}

	// !delete block
	sl.blocks = sl.blocks[:pos]
	sl.mins = sl.mins[:pos]
	sl.shouldRebuildTree = true // TODO: 能否不重建树
	return value
}
func (sl *SortedList32) _adjust(pos int32) {
	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:]
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		sl.mins = append(sl.mins, EMPTY)
		copy(sl.mins[pos+2:], sl.mins[pos+1:])
		sl.mins[pos+1] = sl.blocks[pos+1][0]
		sl.shouldRebuildTree = true
	}
}
