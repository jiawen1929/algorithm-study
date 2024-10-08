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
	"sort"
)

func main() {
	// https://www.luogu.com.cn/problem/SP10707
	// https://www.acwing.com/problem/content/description/2536/
	// 对每个查询，求u到v的路径上顶点种类数
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	values := make([]int32, n) // 顶点权值
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &values[i])

	}
	pool := make(map[int32]int32)
	{
		id := func(o int32) int32 {
			if v, ok := pool[o]; ok {
				return v
			}
			v := int32(len(pool))
			pool[o] = v
			return v
		}
		for i := int32(0); i < n; i++ {
			values[i] = id(values[i])
		}
	}

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	queries := make([][2]int32, q) // u, v
	for i := int32(0); i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		queries[i] = [2]int32{u, v}
	}

	mo := NewMoOnTree(tree, 0)
	for _, q := range queries {
		mo.AddQuery(q[0], q[1])
	}

	res := make([]int32, q)
	counter, kind := make([]int32, int32(len(pool))), int32(0)
	add := func(i int32) {
		counter[values[i]]++
		if counter[values[i]] == 1 {
			kind++
		}
	}
	remove := func(i int32) {
		counter[values[i]]--
		if counter[values[i]] == 0 {
			kind--
		}
	}
	query := func(qid int32) {
		res[qid] = kind
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

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
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
