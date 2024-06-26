// https://www.luogu.com.cn/problem/P5787
// 给定n个节点的图，在k个时间内有m条边会出现后消失。问第i时刻是否是二分图。
// 二分图判定使用可撤销的扩展域并查集维护。

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	mutations := make([][4]int32, m) // (u,v,start,end)
	for i := 0; i < m; i++ {
		var u, v, start, end int32
		fmt.Fscan(in, &u, &v, &start, &end)
		u--
		v--
		mutations[i] = [4]int32{u, v, start, end}
	}

	queries := make([]int32, k)
	for i := range queries {
		queries[i] = int32(i)
	}

	S := NewSegmentTreeDivideAndConquerUndo()
	for id := range mutations {
		start, end := mutations[id][2], mutations[id][3]
		S.AddMutation(int(start), int(end), id)
	}
	for id, time := range queries {
		S.AddQuery(int(time), id)
	}

	res := make([]bool, k)
	checker := NewBipartiteChecker(int32(n)) // 带权并查集

	S.Run(
		func(mutationId int) {
			u, v := mutations[mutationId][0], mutations[mutationId][1]
			checker.Union(u, v)
		},
		func() {
			checker.Undo()
		},
		func(queryId int) {
			res[queryId] = checker.IsBipartite()
		},
	)

	for _, v := range res {
		if v {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

type segMutation struct{ start, end, id int }
type segQuery struct{ time, id int }

// 线段树分治undo版.
// 线段树分治是一种处理动态修改和询问的离线算法.
// 通过将某一元素的出现时间段在线段树上保存到`log(n)`个结点中,
// 我们可以 dfs 遍历整棵线段树，运用可撤销数据结构维护来得到每个时间点的答案.
type SegmentTreeDivideAndConquerUndo struct {
	mutate    func(mutationId int)
	undo      func()
	query     func(queryId int)
	mutations []segMutation
	queries   []segQuery
	nodes     [][]int
}

func NewSegmentTreeDivideAndConquerUndo() *SegmentTreeDivideAndConquerUndo {
	return &SegmentTreeDivideAndConquerUndo{}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerUndo) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.mutations = append(o.mutations, segMutation{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerUndo) AddQuery(time int, id int) {
	o.queries = append(o.queries, segQuery{time, id})
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	mutate: 添加编号为`mutationId`的变更后产生的副作用.
//	undo: 撤销一次`mutate`操作.
//	query: 响应编号为`queryId`的查询.
//	一共调用**O(nlogn)**次`mutate`和`undo`，**O(q)**次`query`.
func (o *SegmentTreeDivideAndConquerUndo) Run(mutate func(mutationId int), undo func(), query func(queryId int)) {
	if len(o.queries) == 0 {
		return
	}
	o.mutate, o.undo, o.query = mutate, undo, query
	times := make([]int, len(o.queries))
	for i := range o.queries {
		times[i] = o.queries[i].time
	}
	sort.Ints(times)
	uniqueInplace(&times)
	usedTimes := make([]bool, len(times)+1)
	usedTimes[0] = true
	for _, e := range o.mutations {
		usedTimes[lowerBound(times, e.start)] = true
		usedTimes[lowerBound(times, e.end)] = true
	}
	for i := 1; i < len(times); i++ {
		if !usedTimes[i] {
			times[i] = times[i-1]
		}
	}
	uniqueInplace(&times)

	n := len(times)
	offset := 1
	for offset < n {
		offset <<= 1
	}
	o.nodes = make([][]int, offset+n)
	for _, e := range o.mutations {
		left := offset + lowerBound(times, e.start)
		right := offset + lowerBound(times, e.end)
		eid := e.id << 1
		for left < right {
			if left&1 == 1 {
				o.nodes[left] = append(o.nodes[left], eid)
				left++
			}
			if right&1 == 1 {
				right--
				o.nodes[right] = append(o.nodes[right], eid)
			}
			left >>= 1
			right >>= 1
		}
	}

	for _, q := range o.queries {
		pos := offset + upperBound(times, q.time) - 1
		o.nodes[pos] = append(o.nodes[pos], (q.id<<1)|1)
	}

	o.dfs(1)
}

func (o *SegmentTreeDivideAndConquerUndo) dfs(now int) {
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.query(id >> 1)
		} else {
			o.mutate(id >> 1)
		}
	}
	if now<<1 < len(o.nodes) {
		o.dfs(now << 1)
	}
	if (now<<1)|1 < len(o.nodes) {
		o.dfs((now << 1) | 1)
	}
	for i := len(curNodes) - 1; i >= 0; i-- {
		if curNodes[i]&1 == 0 {
			o.undo()
		}
	}
}

func uniqueInplace(sorted *[]int) {
	if len(*sorted) == 0 {
		return
	}
	tmp := *sorted
	slow := 0
	for fast := 0; fast < len(tmp); fast++ {
		if tmp[fast] != tmp[slow] {
			slow++
			tmp[slow] = tmp[fast]
		}
	}
	*sorted = tmp[:slow+1]
}

func lowerBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func upperBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

type SetValueStep struct {
	cell     *int32
	oldValue int32
	version  int32
}

func NewSetValueStep(cell *int32, oldValue int32, version int32) *SetValueStep {
	return &SetValueStep{cell: cell, oldValue: oldValue, version: version}
}

// func (s *SetValueStep) Apply()  { *s.cell = s.newValue }
func (s *SetValueStep) Revert() { *s.cell = s.oldValue }

// 在线二分图检测.
type BipartiteChecker struct {
	n            int32
	parent       []int32
	rank         []int32
	color        []int32
	version      int32
	firstViolate int32
	history      []*SetValueStep // plugin
}

func NewBipartiteChecker(n int32) *BipartiteChecker {
	res := &BipartiteChecker{
		n:            n,
		parent:       make([]int32, n),
		rank:         make([]int32, n),
		color:        make([]int32, n),
		firstViolate: -1,
	}
	for i := int32(0); i < n; i++ {
		res.parent[i] = i
	}
	return res
}

func (b *BipartiteChecker) IsBipartite() bool {
	return b.firstViolate == -1
}

// (leader, color)
func (b *BipartiteChecker) Find(x int32) (int32, int32) {
	if x == b.parent[x] {
		return x, 0
	}
	leader, color := b.Find(b.parent[x])
	color ^= b.color[x]
	return leader, color
}

func (b *BipartiteChecker) Union(x, y int32) {
	b.version++
	color := int32(1)
	leaderX, distX := b.Find(x)
	x, color = leaderX, color^distX
	leaderY, distY := b.Find(y)
	y, color = leaderY, color^distY
	if x == y {
		if color == 1 && b.firstViolate == -1 {
			b.firstViolate = b.version
		}
		b.setValue(&b.parent[0], b.parent[0])
		return
	}
	if b.rank[x] < b.rank[y] {
		b.setValue(&b.parent[x], y)
		b.setValue(&b.color[x], color)
	} else {
		b.setValue(&b.parent[y], x)
		b.setValue(&b.color[y], color)
		if b.rank[x] == b.rank[y] {
			b.setValue(&b.rank[x], b.rank[x]+1)
		}
	}
}

func (b *BipartiteChecker) Undo() {
	if len(b.history) == 0 {
		return
	}
	v := b.history[len(b.history)-1].version
	if b.firstViolate == v {
		b.firstViolate = -1
	}
	for len(b.history) > 0 && b.history[len(b.history)-1].version == v {
		b.history[len(b.history)-1].Revert()
		b.history = b.history[:len(b.history)-1]
	}
}

func (b *BipartiteChecker) setValue(cell *int32, newValue int32) {
	step := NewSetValueStep(cell, *cell, b.version)
	*cell = newValue // apply
	b.history = append(b.history, step)
}
