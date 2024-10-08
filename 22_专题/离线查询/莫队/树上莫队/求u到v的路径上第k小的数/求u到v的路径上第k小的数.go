// https://ei1333.github.io/library/other/mo-tree.hpp
// https://oi-wiki.org/misc/mo-algo-on-tree/
// https://github.com/EndlessCheng/codeforces-go/blob/53262fb81ffea176cd5f039cec71e3bd266dce83/copypasta/mo.go#L301
// 处理树上的路径相关的离线查询.
// 一般的莫队只能处理线性问题，我们要把树强行压成序列。
// 通过欧拉序转化成序列上的查询，然后用莫队解决。

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"sort"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	// 树上路径第k小.
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2270
	// 对每个查询，求u到v的路径上第k小的数(1-based)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	values := make([]int32, n) // 顶点权值
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree[u-1] = append(tree[u-1], v-1)
		tree[v-1] = append(tree[v-1], u-1)
	}
	queries := make([][3]int32, q) // u, v, k
	for i := int32(0); i < q; i++ {
		var u, v, k int32
		fmt.Fscan(in, &u, &v, &k)
		queries[i] = [3]int32{u - 1, v - 1, k}
	}

	// 离散化顶点权值
	newValues, origin := Discretize32(values)

	mo := NewMoOnTree(tree, 0)
	for _, q := range queries {
		mo.AddQuery(q[0], q[1])
	}

	res := make([]int32, q)
	sl := NewSortedListRangeBlock32(int32(len(newValues)))
	add := func(i int32) {
		sl.Add(newValues[i])
	}
	remove := func(i int32) {
		sl.Remove(newValues[i])
	}
	query := func(qid int32) {
		k := queries[qid][2]
		kth := sl.At(k - 1)
		res[qid] = origin[kth]
	}

	mo.Run(add, remove, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoOnTree32 struct {
	root    int32
	in, vs  []int32
	tree    [][]int32
	queries [][2]int32
}

func NewMoOnTree(tree [][]int32, root int32) *MoOnTree32 {
	return &MoOnTree32{tree: tree, root: root}
}

// 添加从顶点u到顶点v的查询.
func (mo *MoOnTree32) AddQuery(u, v int32) { mo.queries = append(mo.queries, [2]int32{u, v}) }

// 处理每个查询.
//
//	add: 将数据添加到窗口.
//	remove: 将数据从窗口移除.
//	query: 查询窗口内的数据.
func (mo *MoOnTree32) Run(add func(rootId int32), remove func(rootId int32), query func(qid int32)) {
	n := int32(len(mo.tree))

	vs := make([]int32, 0, 2*n)
	tin := make([]int32, n)
	tout := make([]int32, n)

	var initTime func(v, fa int32)
	initTime = func(v, fa int32) {
		tin[v] = int32(len(vs))
		vs = append(vs, v)
		for _, to := range mo.tree[v] {
			if to != fa {
				initTime(to, v)
			}
		}
		tout[v] = int32(len(vs))
		vs = append(vs, v)
	}
	initTime(mo.root, -1)

	lca := _offlineLCA32(mo.tree, mo.queries, mo.root)
	// blockSize := int(math.Round(math.Pow(float64(2*n), 2.0/3)))
	blockSize := int32(math.Ceil(float64(2*n) / math.Sqrt(float64(len(mo.queries)))))
	type Q struct{ lb, l, r, lca, qid int32 }
	qs := make([]Q, len(mo.queries))
	for i := int32(0); i < int32(len(qs)); i++ {
		v, w := mo.queries[i][0], mo.queries[i][1]
		if tin[v] > tin[w] {
			v, w = w, v
		}
		if lca_ := lca[i]; lca_ != v {
			qs[i] = Q{tout[v] / blockSize, tout[v], tin[w] + 1, lca_, i}
		} else {
			qs[i] = Q{tin[v] / blockSize, tin[v], tin[w] + 1, -1, i}
		}
	}

	sort.Slice(qs, func(i, j int) bool {
		a, b := qs[i], qs[j]
		if a.lb != b.lb {
			return a.lb < b.lb
		}
		if a.lb&1 == 0 {
			return a.r < b.r
		}
		return a.r > b.r
	})

	flip := make([]bool, n)
	f := func(u int32) {
		flip[u] = !flip[u]
		if flip[u] {
			add(u)
		} else {
			remove(u)
		}
	}

	l, r := int32(0), int32(0)
	for _, q := range qs {
		for ; r < q.r; r++ {
			f(vs[r])
		}
		for ; l < q.l; l++ {
			f(vs[l])
		}
		for l > q.l {
			l--
			f(vs[l])
		}
		for r > q.r {
			r--
			f(vs[r])
		}
		if q.lca >= 0 {
			f(q.lca)
		}
		query(q.qid)
		if q.lca >= 0 {
			f(q.lca)
		}
	}
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

// LCA离线.
func _offlineLCA32(tree [][]int32, queries [][2]int32, root int32) []int32 {
	n := int32(len(tree))
	ufa := NewUnionFindArray(n)
	st, mark, ptr, res := make([]int32, n), make([]int32, n), make([]int32, n), make([]int32, len(queries))
	for i := 0; i < len(queries); i++ {
		res[i] = -1
	}
	top := 0
	st[top] = root
	for _, q := range queries {
		mark[q[0]]++
		mark[q[1]]++
	}
	q := make([][][2]int32, n)
	for i := int32(0); i < n; i++ {
		q[i] = make([][2]int32, 0, mark[i])
		mark[i] = -1
		ptr[i] = int32(len(tree[i]))
	}
	for i := int32(0); i < int32(len(queries)); i++ {
		u, v := queries[i][0], queries[i][1]
		q[u] = append(q[u], [2]int32{v, i})
		q[v] = append(q[v], [2]int32{u, i})
	}
	run := func(u int32) bool {
		for ptr[u] != 0 {
			v := tree[u][ptr[u]-1]
			ptr[u]--
			if mark[v] == -1 {
				top++
				st[top] = v
				return true
			}
		}
		return false
	}

	for top != -1 {
		u := st[top]
		if mark[u] == -1 {
			mark[u] = u
		} else {
			ufa.Union(u, tree[u][ptr[u]])
			mark[ufa.Find(u)] = u
		}

		if !run(u) {
			for _, v := range q[u] {
				if mark[v[0]] != -1 && res[v[1]] == -1 {
					res[v[1]] = mark[ufa.Find(v[0])]
				}
			}
			top--
		}
	}

	return res
}

type _unionFindArray struct {
	data []int32
}

func NewUnionFindArray(n int32) *_unionFindArray {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &_unionFindArray{data: data}
}

func (ufa *_unionFindArray) Union(key1, key2 int32) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *_unionFindArray) Find(key int32) int32 {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

const INF int = 1e18

type SortedListRangeBlock32 struct {
	_blockSize  int32   // 每个块的大小.
	_len        int32   // 所有数的个数.
	_counter    []int32 // 每个数出现的次数.
	_blockCount []int32 // 每个块的数的个数.
	_belong     []int32 // 每个数所在的块.
	_blockSum   []int   // 每个块的和.
}

// 值域分块模拟SortedList.
// `O(1)`add/remove，`O(sqrt(n))`查询.
// 一般配合莫队算法使用.
//
//	max:值域的最大值.0 <= max <= 1e6.
//	iterable:初始值.
func NewSortedListRangeBlock32(max int32, nums ...int32) *SortedListRangeBlock32 {
	max += 5
	size := int32(math.Sqrt(float64(max)))
	count := 1 + (max / size)
	sl := &SortedListRangeBlock32{
		_blockSize:  size,
		_counter:    make([]int32, max),
		_blockCount: make([]int32, count),
		_belong:     make([]int32, max),
		_blockSum:   make([]int, count),
	}
	for i := int32(0); i < max; i++ {
		sl._belong[i] = i / size
	}
	if len(nums) > 0 {
		sl.Update(nums...)
	}
	return sl
}

// O(1).
func (sl *SortedListRangeBlock32) Add(value int32) {
	sl._counter[value]++
	pos := sl._belong[value]
	sl._blockCount[pos]++
	sl._blockSum[pos] += int(value)
	sl._len++
}

// O(1).
func (sl *SortedListRangeBlock32) Remove(value int32) {
	sl._counter[value]--
	pos := sl._belong[value]
	sl._blockCount[pos]--
	sl._blockSum[pos] -= int(value)
	sl._len--
}

// O(1).
func (sl *SortedListRangeBlock32) Discard(value int32) bool {
	if !sl.Has(value) {
		return false
	}
	sl.Remove(value)
	return true
}

// O(1).
func (sl *SortedListRangeBlock32) Has(value int32) bool {
	return sl._counter[value] > 0
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) At(index int32) int32 {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	for i := int32(0); i < int32(len(sl._blockCount)); i++ {
		count := sl._blockCount[i]
		if index < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if index < numCount {
					return num
				}
				index -= numCount
				num++
			}
		}
		index -= count
	}
	panic("unreachable")
}

// 严格小于 value 的元素个数.
// 也即第一个大于等于 value 的元素的下标.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) BisectLeft(value int32) int32 {
	pos := sl._belong[value]
	res := int32(0)
	for i := int32(0); i < pos; i++ {
		res += sl._blockCount[i]
	}
	for v := pos * sl._blockSize; v < value; v++ {
		res += sl._counter[v]
	}
	return res
}

// 小于等于 value 的元素个数.
// 也即第一个大于 value 的元素的下标.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) BisectRight(value int32) int32 {
	return sl.BisectLeft(value + 1)
}

func (sl *SortedListRangeBlock32) Count(value int32) int32 {
	return sl._counter[value]
}

// 返回范围 `[min, max]` 内数的个数.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) CountRange(min, max int32) int32 {
	if min > max {
		return 0
	}

	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := int32(0)
		for i := min; i <= max; i++ {
			res += sl._counter[i]
		}
		return res
	}

	res := int32(0)
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += sl._counter[v]
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockCount[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += sl._counter[v]
	}
	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Lower(value int32) (res int32, ok bool) {
	pos := sl._belong[value]
	start := pos * sl._blockSize
	for v := value - 1; v >= start; v-- {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos - 1; i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i + 1) * sl._blockSize
		for {
			if sl._counter[num] > 0 {
				return num, true
			}
			num--
		}
	}

	return
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Higher(value int32) (res int32, ok bool) {
	pos := sl._belong[value]
	end := (pos + 1) * sl._blockSize
	for v := value + 1; v < end; v++ {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos + 1; i < int32(len(sl._blockCount)); i++ {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := i * sl._blockSize
		for {
			if sl._counter[num] > 0 {
				return num, true
			}
			num++
		}
	}

	return
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Floor(value int32) (res int32, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Lower(value)
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Ceiling(value int32) (res int32, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Higher(value)
}

// 返回区间 `[start, end)` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) SumSlice(start, end int32) int {
	if start < 0 {
		start += sl._len
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl._len
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return 0
	}

	res := 0
	remain := end - start
	cur, index := sl._findKth(start)
	sufCount := sl._counter[cur] - index
	if sufCount >= remain {
		return int(remain) * int(cur)
	}

	res += int(sufCount) * int(cur)
	remain -= sufCount
	cur++

	// 当前块内的和
	blockEnd := (sl._belong[cur] + 1) * sl._blockSize
	for remain > 0 && cur < blockEnd {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		res += int(real) * int(cur)
		remain -= real
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for pos < int32(len(sl._blockCount)) && remain >= sl._blockCount[pos] {
		res += sl._blockSum[pos]
		remain -= sl._blockCount[pos]
		pos++
		cur += sl._blockSize
	}

	// 剩余的
	for remain > 0 {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		res += int(real) * int(cur)
		remain -= real
		cur++
	}

	return res
}

// 返回范围 `[min, max]` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock32) SumRange(min, max int32) int {
	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := 0
		for i := min; i <= max; i++ {
			res += int(sl._counter[i]) * int(i)
		}
		return res
	}

	res := 0
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += int(sl._counter[v]) * int(v)
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockSum[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += int(sl._counter[v]) * int(v)
	}
	return res
}

func (sl *SortedListRangeBlock32) ForEach(f func(value, index int32), reverse bool) {
	if reverse {
		ptr := int32(0)
		for i := int32(len(sl._counter) - 1); i >= 0; i-- {
			count := sl._counter[i]
			for j := int32(0); j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	} else {
		ptr := int32(0)
		for i := int32(0); i < int32(len(sl._counter)); i++ {
			count := sl._counter[i]
			for j := int32(0); j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	}
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Pop(index int32) int32 {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	value := sl.At(index)
	sl.Remove(value)
	return value
}

func (sl *SortedListRangeBlock32) Slice(start, end int32) []int32 {
	if start < 0 {
		start += sl._len
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl._len
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return nil
	}

	res := make([]int32, end-start)
	count := int32(0)
	sl.Enumerate(start, end, func(value int32) {
		res[count] = value
		count++
	}, false)

	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock32) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedListRangeBlock32) Enumerate(start, end int32, f func(value int32), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return
	}

	remain := end - start
	cur, index := sl._findKth(start)
	sufCount := sl._counter[cur] - index
	real := sufCount
	if real > remain {
		real = remain
	}
	if f != nil {
		for i := int32(0); i < real; i++ {
			f(cur)
		}
	}
	if erase {
		for i := int32(0); i < real; i++ {
			sl.Remove(cur)
		}
	}
	remain -= sufCount
	if remain == 0 {
		return
	}
	cur++

	// 当前块内
	blockEnd := (sl._belong[cur] + 1) * sl._blockSize
	for remain > 0 && cur < blockEnd {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		remain -= real
		if f != nil {
			for i := int32(0); i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := int32(0); i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for count := sl._blockCount[pos]; remain >= count; {
		remain -= count
		if f != nil {
			for v := cur; v < cur+sl._blockSize; v++ {
				c := sl._counter[v]
				for i := int32(0); i < c; i++ {
					f(v)
				}
			}
		}
		if erase {
			for v := cur; v < cur+sl._blockSize; v++ {
				sl._counter[v] = 0
			}
			sl._len -= count
			sl._blockCount[pos] = 0
			sl._blockSum[pos] = 0
		}
		pos++
		cur += sl._blockSize
	}

	// 剩余的
	for remain > 0 {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		remain -= real
		if f != nil {
			for i := int32(0); i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := int32(0); i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}
}

func (sl *SortedListRangeBlock32) Clear() {
	for i := range sl._counter {
		sl._counter[i] = 0
	}
	for i := range sl._blockCount {
		sl._blockCount[i] = 0
	}
	for i := range sl._blockSum {
		sl._blockSum[i] = 0
	}
	sl._len = 0
}

func (sl *SortedListRangeBlock32) Update(values ...int32) {
	for _, value := range values {
		sl.Add(value)
	}
}

func (sl *SortedListRangeBlock32) Merge(other *SortedListRangeBlock32) {
	other.ForEach(func(value, _ int32) {
		sl.Add(value)
	}, false)
}

func (sl *SortedListRangeBlock32) String() string {
	sb := make([]string, 0, sl._len)
	sl.ForEach(func(value, _ int32) {
		sb = append(sb, fmt.Sprintf("%d", value))
	}, false)
	return fmt.Sprintf("SortedListRangeBlock{%s}", strings.Join(sb, ", "))
}

func (sl *SortedListRangeBlock32) Len() int32 {
	return sl._len
}

func (sl *SortedListRangeBlock32) Min() int32 {
	return sl.At(0)
}

func (sl *SortedListRangeBlock32) Max() int32 {
	if sl._len == 0 {
		panic("empty")
	}

	for i := int32(len(sl._blockCount) - 1); i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i+1)*sl._blockSize - 1
		for {
			if sl._counter[num] > 0 {
				return num
			}
			num--
		}
	}

	panic("unreachable")
}

// 返回索引在`kth`处的元素的`value`,以及该元素是`value`中的第几个(`index`).
func (sl *SortedListRangeBlock32) _findKth(kth int32) (value, index int32) {
	for i := int32(0); i < int32(len(sl._blockCount)); i++ {
		count := sl._blockCount[i]
		if kth < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if kth < numCount {
					return num, kth
				}
				kth -= numCount
				num++
			}
		}
		kth -= count
	}

	panic("unreachable")
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize32(nums []int32) (newNums []int32, origin []int32) {
	newNums = make([]int32, len(nums))
	origin = make([]int32, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
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
