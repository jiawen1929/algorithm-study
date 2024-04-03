// https://ei1333.github.io/library/graph/mst/boruvka.hpp
// Boruvka(最小全域木) 在线最小生成树
// 不预先给出图，
// 而是给定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
)

const INF int = 2e18

func main() {
	// SpeedRunMSTEasy()
	SpeedRunMSTHard()
}

// P - MST (Easy)
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_g
// 给定数组A，长度为N。
// 给定一张n个顶点的完全图，边(u,v)的权值为A[u]*A[v]。
// 求这张图的最小生成树的权值。
// n<=2e5.
func SpeedRunMSTEasy() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	sort.Ints(A) // !从小到大排序

	set := NewFastSetFrom(n, func(i int) bool { return true })
	setUsed := func(u int) { set.Erase(u) }
	setUnused := func(u int) { set.Insert(u) }
	findUnused := func(u int) (v int, cost int) {
		min_ := set.Next(0)
		if min_ == n {
			return -1, -1
		}
		max_ := set.Prev(n)
		best := min_
		if A[max_]*A[u] < A[best]*A[u] {
			best = max_
		}
		return best, A[best] * A[u]
	}

	edges := OnlineMST(n, setUsed, setUnused, findUnused)
	res := 0
	for _, e := range edges {
		res += e[2]
	}
	fmt.Fprintln(out, res)
}

// P - MST (Hard)
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
// 给定两个数组A和B，长度为N。
// 给定一张n个顶点的完全图，边(u,v)的权值为A[u]*A[v]+B[u]*B[v]。
// 求这张图的最小生成树的权值。
// n<=5e4.
//
// cht维护二维点集，每次求出给定(x,y)条件下使得ax+by最小的点(a,b).
// 总时间复杂度O(n(logn)^2)
func SpeedRunMSTHard() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &B[i])
	}

	var cht *LineContainer2DWithId
	init := func() {
		if cht == nil {
			cht = NewLineContainer2DWithId(n)
		} else {
			cht.Clear()
		}
	}
	add := func(u int) { cht.Add(A[u], B[u], int32(u)) }
	find := func(u int) (v int, cost int) {
		f, id := cht.QueryMin(A[u], B[u])
		return int(id), f
	}
	edges := OnlineMSTIncremental(n, init, add, find)
	res := 0
	for _, e := range edges {
		res += e[2]
	}
	fmt.Fprintln(out, res)
}

// Brouvka
//
//	不预先给出图，而是指定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。
//	findUnused(u)：返回 unused 中与 u 权值最小的点 v 和边权 cost
//	              如果不存在，返回 (-1,*)
func OnlineMST(
	n int,
	setUsed func(u int), setUnused func(u int), findUnused func(u int) (v int, cost int),
) (res [][3]int) {
	uf := newUnionFindArraySimple(n)
	res = make([][3]int, 0, n-1)
	for {
		updated := false
		groups := make([][]int, n)
		cand := make([][3]int, n) // [u, v, cost]
		for v := 0; v < n; v++ {
			cand[v] = [3]int{-1, -1, -1}
		}

		for v := 0; v < n; v++ {
			leader := uf.Find(v)
			groups[leader] = append(groups[leader], v)
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			for _, x := range groups[v] {
				setUsed(x)
			}
			for _, x := range groups[v] {
				y, cost := findUnused(x)
				if y == -1 {
					continue
				}
				a, c := cand[v][0], cand[v][2]
				if a == -1 || cost < c {
					cand[v] = [3]int{x, y, cost}
				}
			}
			for _, x := range groups[v] {
				setUnused(x)
			}
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			a, b, c := cand[v][0], cand[v][1], cand[v][2]
			if a == -1 {
				continue
			}
			updated = true
			if uf.Union(a, b) {
				res = append(res, [3]int{a, b, c})
			}
		}

		if !updated {
			break
		}
	}

	return res
}

// Brouvka 在线最小生成树，配合incremental(可增量计算的)数据结构使用.
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
func OnlineMSTIncremental(
	n int,
	init func(), // 初始化数据结构，大约调用 2logN 次.
	add func(u int), // 添加搜索点
	find func(u int) (v int, cost int), // (v,cost) or (-1,*)
) (res [][3]int) {
	uf := newUnionFindArraySimple(n)
	res = make([][3]int, 0, n-1)

	for uf.Part > 1 {
		updated := false
		groups := make([][]int, n)
		for v := 0; v < n; v++ {
			leader := uf.Find(v)
			groups[leader] = append(groups[leader], v)
		}
		weight := make([]int, n)
		for i := 0; i < n; i++ {
			weight[i] = INF
		}
		who := make([][2]int, n)
		for i := 0; i < n; i++ {
			who[i] = [2]int{-1, -1}
		}

		init()
		for i := 0; i < n; i++ {
			for _, v := range groups[i] {
				w, x := find(v)
				if x < weight[i] {
					weight[i] = x
					who[i] = [2]int{v, w}
				}
			}
			for _, v := range groups[i] {
				add(v)
			}
		}

		init()
		for i := n - 1; i >= 0; i-- {
			for _, v := range groups[i] {
				w, x := find(v)
				if x < weight[i] {
					weight[i] = x
					who[i] = [2]int{v, w}
				}
			}
			for _, v := range groups[i] {
				add(v)
			}
		}

		for i := 0; i < n; i++ {
			a, b := who[i][0], who[i][1]
			if a == -1 {
				continue
			}
			if uf.Union(a, b) {
				updated = true
				res = append(res, [3]int{a, b, weight[i]})
			}
		}

		if !updated {
			break
		}
	}

	return res
}

// O(n^2) Prim求完全图最小生成树.
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
func Prim(n int, cost func(u, v int) int) [][3]int {
	res := make([][3]int, 0, n-1)
	weight := make([]int, n)
	for i := 0; i < n; i++ {
		weight[i] = INF
	}
	to := make([]int, n)
	add := func(v int) {
		for w := 0; w < n; w++ {
			if to[w] != -1 && chmin(&weight[w], cost(v, w)) {
				to[w] = v
			}
		}
		weight[v] = INF
		to[v] = -1
	}
	add(0)
	for i := 0; i < n-1; i++ {
		argMin := 0
		for j := 1; j < n; j++ {
			if weight[j] < weight[argMin] {
				argMin = j
			}
		}
		res = append(res, [3]int{to[argMin], argMin, weight[argMin]})
		add(argMin)
	}
	return res
}

type unionFindArraySimple struct {
	Part int
	n    int
	data []int32
}

func newUnionFindArraySimple(n int) *unionFindArraySimple {
	data := make([]int32, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &unionFindArraySimple{Part: n, n: n, data: data}
}

func (u *unionFindArraySimple) Union(key1 int, key2 int) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *unionFindArraySimple) Find(key int) int {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = int32(u.Find(int(u.data[key])))
	return int(u.data[key])
}

func (u *unionFindArraySimple) GetSize(key int) int {
	return int(-u.data[u.Find(key)])
}

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

func chmin(a *int, b int) bool {
	if *a > b {
		*a = b
		return true
	}
	return false
}

type Line struct {
	k, b   int
	p1, p2 int // p=p1/p2
}

type LineContainer2DWithId struct {
	minCHT, maxCHT       *_LineContainer
	kMax, kMin           int
	bMax, bMin           int
	kMaxIndex, kMinIndex int32
	bMaxIndex, bMinIndex int32
	mp                   map[[2]int]int32
	capacity             int
}

func NewLineContainer2DWithId(capacity int) *LineContainer2DWithId {
	return &LineContainer2DWithId{
		minCHT: _NewLineContainer(true, capacity),
		maxCHT: _NewLineContainer(false, capacity),
		kMax:   -INF, kMin: INF, bMax: -INF, bMin: INF,
		kMaxIndex: -1, kMinIndex: -1, bMaxIndex: -1, bMinIndex: -1,
		mp:       make(map[[2]int]int32, capacity),
		capacity: capacity,
	}
}

// 追加 a*x + b*y.
func (lc *LineContainer2DWithId) Add(a, b int, id int32) {
	lc.minCHT.Add(b, a)
	lc.maxCHT.Add(b, a)
	pair := [2]int{a, b}
	lc.mp[pair] = id

	if a > lc.kMax {
		lc.kMax = a
		lc.kMaxIndex = id
	}
	if a < lc.kMin {
		lc.kMin = a
		lc.kMinIndex = id
	}
	if b > lc.bMax {
		lc.bMax = b
		lc.bMaxIndex = id
	}
	if b < lc.bMin {
		lc.bMin = b
		lc.bMinIndex = id
	}
}

// 查询 x=xi,y=yi 时的最大值 max_{a,b} (ax + by)和对应的点id.
func (lc *LineContainer2DWithId) QueryMax(x, y int) (int, int32) {
	if lc.minCHT.Size() == 0 {
		return -INF, -1
	}

	if x == 0 {
		if y > 0 {
			return lc.bMax * y, lc.bMaxIndex
		}
		return lc.bMin * y, lc.bMinIndex
	}
	if y == 0 {
		if x > 0 {
			return lc.kMax * x, lc.kMaxIndex
		}
		return lc.kMin * x, lc.kMinIndex
	}

	// y/x
	if x > 0 {
		l := lc.maxCHT.sl.BisectLeftByPair(y, x)
		line := lc.maxCHT.sl.At(l)
		a := line.b
		b := line.k
		return a*x + b*y, lc.mp[[2]int{a, b}]
	}
	l := lc.minCHT.sl.BisectLeftByPair(y, x)
	line := lc.minCHT.sl.At(l)
	a := -line.b
	b := -line.k
	return a*x + b*y, lc.mp[[2]int{a, b}]
}

// 查询 x=xi,y=yi 时的最小值 min_{a,b} (ax + by).
func (lc *LineContainer2DWithId) QueryMin(x, y int) (int, int32) {
	v, i := lc.QueryMax(-x, -y)
	return -v, i
}

func (lc *LineContainer2DWithId) Clear() {
	lc.minCHT.Clear()
	lc.maxCHT.Clear()
	lc.kMax, lc.kMin = -INF, INF
	lc.bMax, lc.bMin = -INF, INF
	lc.kMaxIndex, lc.kMinIndex = -1, -1
	lc.bMaxIndex, lc.bMinIndex = -1, -1
	lc.mp = make(map[[2]int]int32, lc.capacity)
}

type _LineContainer struct {
	minimize bool
	sl       *SpecializedSortedList
}

func _NewLineContainer(minimize bool, capacity int) *_LineContainer {
	return &_LineContainer{
		minimize: minimize,
		sl:       NewSpecializedSortedList(func(a, b S) bool { return a.k < b.k }),
	}
}

func (lc *_LineContainer) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}

	newLine := &Line{k: k, b: m}
	lc.sl.Add(newLine)
	it1 := lc.sl.BisectRightByK(newLine.k) - 1
	it2 := it1
	line2 := lc.sl.At(it2)
	it1++
	it3 := it2
	for lc.insect(line2, lc.sl.At(it1)) {
		lc.sl.Pop(it1)
	}

	if it3 != 0 {
		it3--
		line3 := lc.sl.At(it3)
		if lc.insect(line3, line2) {
			lc.sl.Pop(it2)
			lc.insect(line3, lc.sl.At(it2))
		}
	}

	if it3 == 0 {
		return
	}

	dp1, dp2 := lc.sl.At(it3-1), lc.sl.At(it3)
	for it3 != 0 {
		it2 := it3
		if lessPair(dp1.p1, dp1.p2, dp2.p1, dp2.p2) {
			break
		}
		it3--
		lc.sl.Pop(it2)
		lc.insect(dp1, lc.sl.At(it2))
		dp1, dp2 = lc.sl.At(it3-1), dp1
	}
}

// 查询 kx + m 的最小值（或最大值).
func (lc *_LineContainer) Query(x int) int {
	if lc.sl.Len() == 0 {
		panic("empty container")
	}
	pos := lc.sl.BisectLeftByPair(x, 1)
	line := lc.sl.At(pos)
	v := line.k*x + line.b
	if lc.minimize {
		return -v
	}
	return v
}

func (lc *_LineContainer) Size() int32 { return lc.sl.Len() }

func (lc *_LineContainer) Clear() { lc.sl.Clear() }

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *_LineContainer) insect(line1, line2 *Line) bool {
	if line2 == nil {
		line1.p1 = INF
		line1.p2 = 1
		return false
	}
	if line1.k == line2.k {
		if line1.b > line2.b {
			line1.p1 = INF
			line1.p2 = 1
		} else {
			line1.p1 = INF
			line1.p2 = -1
		}
	} else {
		// lc_div
		line1.p1 = line2.b - line1.b
		line1.p2 = line1.k - line2.k
	}
	return !lessPair(line1.p1, line1.p2, line2.p1, line2.p2)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 分母不为0的分数比较大小
//
//	a1/b1 < a2/b2
func lessPair(a1, b1, a2, b2 int) bool {
	if a1 == INF || a2 == INF { // 有一个是+-INF
		return a1/b1 < a2/b2
	}
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	return diff^mul < 0
}

const _LOAD int32 = 50 // 75/100/150/200

type S = *Line

type SpecializedSortedList struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSpecializedSortedList(less func(a, b S) bool, elements ...S) *SpecializedSortedList {
	elements = append(elements[:0:0], elements...)
	res := &SpecializedSortedList{less: less}
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

func (sl *SpecializedSortedList) Add(value S) *SpecializedSortedList {
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
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *SpecializedSortedList) Pop(index int32) {
	pos, startIndex := sl._findKth(index)
	sl._delete(pos, startIndex)
}

func (sl *SpecializedSortedList) At(index int32) S {
	if index < 0 || index >= sl.size {
		return nil
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SpecializedSortedList) BisectRightByK(k int) int32 {
	pos, index := sl._locRightByK(k)
	return sl._queryTree(pos) + index
}

func (sl *SpecializedSortedList) BisectLeftByPair(a, b int) int32 {
	pos, index := sl._locLeftByPair(a, b)
	return sl._queryTree(pos) + index
}

func (sl *SpecializedSortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SpecializedSortedList) Len() int32 {
	return sl.size
}

func (sl *SpecializedSortedList) _delete(pos, index int32) {
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

func (sl *SpecializedSortedList) _locLeftByPair(a, b int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !lessPair(sl.mins[mid].p1, sl.mins[mid].p2, a, b) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		last := block[len(block)-1]
		if !lessPair(last.p1, last.p2, a, b) {
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
		if !lessPair(cur[mid].p1, cur[mid].p2, a, b) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _locRight(value S) (pos, index int32) {
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

func (sl *SpecializedSortedList) _locRightByK(k int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if k < sl.mins[mid].k {
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
		if k < cur[mid].k {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _buildTree() {
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

func (sl *SpecializedSortedList) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < int32(len(tree)); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SpecializedSortedList) _queryTree(end int32) int32 {
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

func (sl *SpecializedSortedList) _findKth(k int32) (pos, index int32) {
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
	m := int32(len(tree))
	bitLength := bits.Len32(uint32(m))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
