// https://atcoder.jp/contests/abc324/editorial/7399
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

// 给定一个数组，数组元素为1-n的排列
// 有两种操作：
// 1.把 A[version]中下标大于等于 x 的元素分裂成一个新的数组 Ai(A[version]中保留x个)。
// 2.把 A[version]中值大于 x 的元素分裂成一个新的数组 Ai。
// 这两种操作都不会改变元素相对顺序。
// 输出每次分裂出的数组大小。
//
// SortedList + Deque 维护.
// 启发式分裂：每次分裂出较小的那一半
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	var q int
	fmt.Fscan(in, &q)

	history := make([]*SortedDeque, q+1)
	history[0] = NewSortedDeque(func(a, b S) bool { return a < b }, nums...)
	for i := 1; i < len(history); i++ {
		history[i] = NewSortedDeque(func(a, b S) bool { return a < b })
	}

	for cur := 1; cur < q+1; cur++ {
		var kind, pre, x int
		fmt.Fscan(in, &kind, &pre, &x)

		if kind == 1 { // 将 A[pre] 中下标大于等于 x 的元素分裂成一个新的数组 Ai
			len1 := x
			len2 := history[pre].Len() - x
			if len1 < len2 { // 前面少，拆到前面
				history[cur], history[pre] = history[pre], history[cur]
				for j := 0; j < len1; j++ {
					history[pre].Append(history[cur].PopLeft())
				}
			} else { // 后面少，拆到后面
				for j := 0; j < len2; j++ {
					history[cur].AppendLeft(history[pre].Pop())
				}
			}
		} else { // 将 A[pre] 中值大于 x 的元素分裂成一个新的数组 Ai
			len1 := history[pre].CountLessOrEqual(x)
			len2 := history[pre].Len() - len1
			if len1 < len2 { // 前面少，拆到前面
				history[cur], history[pre] = history[pre], history[cur]
				for j := 0; j < len1; j++ {
					min_ := history[cur].Min()
					history[cur].Remove(min_)
					history[pre].Append(min_)
				}
			} else { // 后面少，拆到后面
				for j := 0; j < len2; j++ {
					max_ := history[pre].Max()
					history[pre].Remove(max_)
					history[cur].AppendLeft(max_)
				}
			}
		}

		fmt.Fprintln(out, history[cur].Len())
	}
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int = 200

type S = int

// 可删除元素、获取第k小值的双端队列.
// !启用删除功能时，需要保证队列中始终不能有重复元素，且删除的元素必须存在于队列中.
type SortedDeque struct {
	sl *SortedList
	dq *RemovableDeque
}

func NewSortedDeque(less func(a, b S) bool, elements ...S) *SortedDeque {
	elements = append(elements[:0:0], elements...)
	res := &SortedDeque{sl: NewSortedList(less, elements...), dq: NewRemovableDeque(len(elements))}
	for _, v := range elements {
		res.dq.Append(v)
	}
	return res
}

func (sd *SortedDeque) Append(value S) {
	sd.sl.Add(value)
	sd.dq.Append(value)
}

func (sd *SortedDeque) AppendLeft(value S) {
	sd.sl.Add(value)
	sd.dq.AppendLeft(value)
}

func (sd *SortedDeque) Pop() S {
	value := sd.dq.Pop()
	sd.sl.Discard(value)
	return value
}

func (sd *SortedDeque) PopLeft() S {
	value := sd.dq.PopLeft()
	sd.sl.Discard(value)
	return value
}

func (sd *SortedDeque) Head() S {
	return sd.dq.Head()
}

func (sd *SortedDeque) Tail() S {
	return sd.dq.Tail()
}

// 删除队列中所有值为value的元素.
func (sd *SortedDeque) Remove(value S) {
	count := sd.dq.Count(value)
	if count == 0 {
		return
	}
	// fast path
	if count == 1 {
		sd.sl.Discard(value)
		sd.dq.Remove(value)
		return
	}
	start := sd.sl.BisectLeft(value)
	end := start + count
	sd.sl.Erase(start, end)
	sd.dq.Remove(value)
}

func (sd *SortedDeque) Min() S {
	return sd.sl.Min()
}

func (sd *SortedDeque) Max() S {
	return sd.sl.Max()
}

func (sd *SortedDeque) Kth(k int) S {
	return sd.sl.At(k)
}

func (sd *SortedDeque) CountLess(value S) int {
	return sd.sl.BisectLeft(value)
}

func (sd *SortedDeque) CountLessOrEqual(value S) int {
	return sd.sl.BisectRight(value)
}

func (sd *SortedDeque) Len() int {
	return sd.sl.Len()
}

// 使用分块+树状数组维护的有序序列.
type SortedList struct {
	less              func(a, b S) bool
	size              int
	blocks            [][]S
	mins              []S
	tree              []int
	shouldRebuildTree bool
}

func NewSortedList(less func(a, b S) bool, elements ...S) *SortedList {
	elements = append(elements[:0:0], elements...)
	res := &SortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	blocks := [][]S{}
	for start := 0; start < n; start += _LOAD {
		end := min(start+_LOAD, n)
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

func (sl *SortedList) Add(value S) *SortedList {
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
	if n := len(sl.blocks[pos]); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SortedList) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < len(sl.blocks[pos]) && sl.blocks[pos][index] == value
}

func (sl *SortedList) Discard(value S) bool {
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

func (sl *SortedList) Pop(index int) S {
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

func (sl *SortedList) At(index int) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedList) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedList) Lower(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Higher(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedList) Floor(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedList) Ceiling(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedList) BisectLeft(value S) int {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedList) BisectRight(value S) int {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedList) Count(value S) int {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedList) ForEach(f func(value S, index int), reverse bool) {
	if !reverse {
		count := 0
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				f(block[j], count)
				count++
			}
		}
	} else {
		count := 0
		for i := len(sl.blocks) - 1; i >= 0; i-- {
			block := sl.blocks[i]
			for j := len(block) - 1; j >= 0; j-- {
				f(block[j], count)
				count++
			}
		}
	}
}

func (sl *SortedList) Enumerate(start, end int, f func(value S), erase bool) {
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
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endIndex := min(len(block), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == len(block) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				for i := startIndex; i < endIndex; i++ {
					sl._updateTree(pos, -1)
				}
				sl.blocks[pos] = append(block[:startIndex], block[endIndex:]...)
				sl.mins[pos] = sl.blocks[pos][0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SortedList) Slice(start, end int) []S {
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
	for ; count > 0 && pos < len(sl.blocks); pos++ {
		block := sl.blocks[pos]
		endPos := min(len(block), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedList) Range(min, max S) []S {
	if sl.less(max, min) {
		return nil
	}
	res := []S{}
	pos := sl._locBlock(min)
	for i := pos; i < len(sl.blocks); i++ {
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

func (sl *SortedList) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedList) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedList) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value S, index int) {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedList) Len() int {
	return sl.size
}

func (sl *SortedList) _delete(pos, index int) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], sl.blocks[pos][index+1:]...)
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
	sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
	sl.shouldRebuildTree = true
}

func (sl *SortedList) _locLeft(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := -1
	right := len(sl.blocks) - 1
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
	right = len(cur)
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

func (sl *SortedList) _locRight(value S) (pos, index int) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := 0
	right := len(sl.blocks)
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
	right = len(cur)
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

func (sl *SortedList) _locBlock(value S) int {
	left, right := -1, len(sl.blocks)-1
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

func (sl *SortedList) _buildTree() {
	sl.tree = make([]int, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = len(sl.blocks[i])
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

func (sl *SortedList) _updateTree(index, delta int) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < len(tree); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedList) _queryTree(end int) int {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := 0
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SortedList) _findKth(k int) (pos, index int) {
	if k < len(sl.blocks[0]) {
		return 0, k
	}
	last := len(sl.blocks) - 1
	lastLen := len(sl.blocks[last])
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
		if next < len(tree) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

type Value = int

type Pair = struct {
	value     Value
	addedTime int
}

type RemovableDeque struct {
	queue       *dq
	counter     map[Value]int
	removedTime map[Value]int
	length      int
	time        int
}

func NewRemovableDeque(cap int) *RemovableDeque {
	return &RemovableDeque{
		queue:       newDq(cap),
		counter:     make(map[Value]int),
		removedTime: make(map[Value]int),
		length:      0,
		time:        0,
	}
}

func (rq *RemovableDeque) Append(value Value) {
	rq.length++
	rq.queue.Append(Pair{value, rq.time})
	rq.counter[value]++
}

func (rq *RemovableDeque) AppendLeft(value Value) {
	rq.length++
	rq.queue.AppendLeft(Pair{value, rq.time})
	rq.counter[value]++
}

func (rq *RemovableDeque) Pop() Value {
	rq.length--
	rq._normalizeTail()
	res := rq.queue.Pop().value
	if _, ok := rq.counter[res]; ok {
		rq.counter[res]--
		if rq.counter[res] == 0 {
			delete(rq.counter, res)
		}
	}
	return res
}

func (rq *RemovableDeque) PopLeft() Value {
	rq.length--
	rq._normalizeHead()
	res := rq.queue.PopLeft().value
	if _, ok := rq.counter[res]; ok {
		rq.counter[res]--
		if rq.counter[res] == 0 {
			delete(rq.counter, res)
		}
	}
	return res
}

func (rq *RemovableDeque) Head() Value {
	rq._normalizeHead()
	return rq.queue.Head().value
}

func (rq *RemovableDeque) Tail() Value {
	rq._normalizeTail()
	return rq.queue.Tail().value
}

// 删除deque中所有值为value的元素.
func (rq *RemovableDeque) Remove(value Value) {
	if _, ok := rq.counter[value]; ok {
		rq.length -= rq.counter[value]
		delete(rq.counter, value)
		rq.removedTime[value] = rq.time
		rq.time++
	}
}

func (rq *RemovableDeque) Count(value Value) int {
	return rq.counter[value]
}

func (rq *RemovableDeque) Empty() bool {
	return rq.length == 0
}

func (rq *RemovableDeque) Len() int {
	return rq.length
}

func (rq *RemovableDeque) String() string {
	res := make([]Value, 0, rq.length)
	for i := 0; i < rq.length; i++ {
		p := rq.queue.At(i)
		v, t := p.value, p.addedTime
		if _, ok := rq.removedTime[v]; ok && t <= rq.removedTime[v] {
			continue
		}
		res = append(res, v)
	}
	return fmt.Sprint(res)
}

func (rq *RemovableDeque) _normalizeHead() {
	for rq.queue.Size() > 0 {
		p := rq.queue.Head()
		v, t := p.value, p.addedTime
		if _, ok := rq.removedTime[v]; ok && t <= rq.removedTime[v] {
			rq.queue.PopLeft()
		} else {
			break
		}
	}
}

func (rq *RemovableDeque) _normalizeTail() {
	for rq.queue.Size() > 0 {
		p := rq.queue.Tail()
		v, t := p.value, p.addedTime
		if _, ok := rq.removedTime[v]; ok && t <= rq.removedTime[v] {
			rq.queue.Pop()
		} else {
			break
		}
	}
}

type dq struct{ l, r []Pair }

func newDq(cap int) *dq { return &dq{make([]Pair, 0, 1+cap/2), make([]Pair, 0, 1+cap/2)} }

func (q *dq) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q *dq) Size() int {
	return len(q.l) + len(q.r)
}

func (q *dq) AppendLeft(v Pair) {
	q.l = append(q.l, v)
}

func (q *dq) Append(v Pair) {
	q.r = append(q.r, v)
}

func (q *dq) PopLeft() Pair {
	var v Pair
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return v
}

func (q *dq) Pop() Pair {
	var v Pair
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return v
}

func (q *dq) Head() Pair {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q *dq) Tail() Pair {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q *dq) At(i int) Pair {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
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
