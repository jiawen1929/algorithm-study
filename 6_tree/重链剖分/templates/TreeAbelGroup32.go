// 单点修改，路径和查询，子树和查询

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupoVertexAddPathSum()
	// yuki1641()
}

func demo() {

	{
		//    0
		//   / \
		//  1   2
		//     / \
		//    3   4

		tree := NewTree32(5)
		tree.AddEdge(0, 1, 0)
		tree.AddEdge(0, 2, 0)
		tree.AddEdge(2, 3, 0)
		tree.AddEdge(2, 4, 0)
		tree.Build(0)

		S := NewTreeAbleGroup(tree, false, true, true)
		S.Build(func(vidOrEid int32) E { return int(vidOrEid) })
		fmt.Println(S.QuerySubtree(0))          // 7
		fmt.Println(S.QuerySubtree(1))          // 1
		fmt.Println(S.QuerySubtree(2))          // 9
		fmt.Println(S.QuerySubtree(3))          // 3
		fmt.Println(S.QuerySubtree(4))          // 4
		fmt.Println(S.QueryPath(1, 3))          // 4
		fmt.Println(S.QuerySubtreeRooted(0, 3)) // 4
		S.Add(3, 10)
		fmt.Println(S.QuerySubtree(0))          // 20
		fmt.Println(S.QuerySubtreeRooted(4, 3)) // 20
	}
}

// https://judge.yosupo.jp/problem/vertex_add_path_sum
func yosupoVertexAddPathSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := 0; i < int(n); i++ {
		fmt.Fscan(in, &weights[i])
	}
	tree := NewTree32(n)
	for i := 1; i < int(n); i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree.AddEdge(u, v, 0)
	}
	tree.Build(0)

	S := NewTreeAbleGroup(tree, false, true, false)
	S.Build(func(vidOrEid int32) E { return weights[vidOrEid] })
	for i := 0; i < int(q); i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 0 {
			var v, x int32
			fmt.Fscan(in, &v, &x)
			S.Add(v, int(x))
		} else {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			fmt.Fprintln(out, S.QueryPath(u, v))
		}
	}
}

// https://yukicoder.me/problems/no/1641
func yuki1641() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	weights := make([]int, n)
	for i := 0; i < int(n); i++ {
		fmt.Fscan(in, &weights[i])
	}
	tree := NewTree32(n)
	for i := 1; i < int(n); i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree.AddEdge(u, v, 0)
	}
	tree.Build(0)

	S := NewTreeAbleGroup(tree, false, false, true)
	S.Build(func(vidOrEid int32) E { return weights[vidOrEid] })
	for i := 0; i < int(q); i++ {
		var t, x, y int32
		fmt.Fscan(in, &t, &x, &y)
		x--
		if t == 1 {
			S.Add(x, int(y))
		}
		if t == 2 {
			fmt.Fprintln(out, S.QuerySubtree(x))
		}
	}
}

type E = int

func e() E { return 0 }

func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

// func op(a, b E) E { return a ^ b }
// func inv(a E) E   { return a }

type TreeAbelGroup struct {
	pathQuery, subtreeQuery bool
	edge                    int32
	n                       int32
	tree                    *Tree32
	bit, bitSubtree         *bitGroup32
}

func NewTreeAbleGroup(tree *Tree32, edge bool, pathQuery, subtreeQuery bool) *TreeAbelGroup {
	var edgeValue int32
	if edge {
		edgeValue = 1
	}
	return &TreeAbelGroup{pathQuery: pathQuery, subtreeQuery: subtreeQuery, edge: edgeValue, n: tree.n, tree: tree}
}

func (tag *TreeAbelGroup) Build(f func(vidOrEid int32) E) {
	bitRaw1 := make([]E, 2*tag.n)
	bitRaw2 := make([]E, tag.n)
	tree := tag.tree
	for v := int32(0); v < tag.n; v++ {
		var x E
		if tag.edge == 0 {
			x = f(v)
		} else {
			if v == 0 {
				x = e()
			} else {
				x = f(tree.vToE[v])
			}
		}
		bitRaw1[tree.ELid(v)] = x
		bitRaw1[tree.ERid(v)] = inv(x)
		bitRaw2[tree.Lid[v]] = x
	}
	if tag.pathQuery {
		tag.bit = newBITGroup32From(2*tag.n, func(index int32) E {
			return bitRaw1[index]
		})
	}
	if tag.subtreeQuery {
		tag.bitSubtree = newBITGroup32From(tag.n, func(index int32) E {
			return bitRaw2[index]
		})
	}
}

func (tag *TreeAbelGroup) Add(i int32, x E) {
	v := i
	if tag.edge != 0 {
		v = tag.tree.EToV(i)
	}
	if tag.pathQuery {
		tag.bit.Update(tag.tree.ELid(v), x)
		tag.bit.Update(tag.tree.ERid(v), inv(x))
	}
	if tag.subtreeQuery {
		tag.bitSubtree.Update(tag.tree.Lid[v], x)
	}
}

func (tag *TreeAbelGroup) QueryPath(from, to int32) E {
	if !tag.pathQuery {
		panic("path query not enabled")
	}
	lca := tag.tree.Lca(from, to)
	x1 := tag.bit.QueryRange(tag.tree.ELid(lca)+1, tag.tree.ELid(from)+1)
	x2 := tag.bit.QueryRange(tag.tree.ELid(lca)+tag.edge, tag.tree.ELid(to)+1)
	return op(x1, x2)
}

func (tag *TreeAbelGroup) QuerySubtree(u int32) E {
	return tag.QuerySubtreeRooted(u, -1)
}

func (tag *TreeAbelGroup) QuerySubtreeRooted(u, root int32) E {
	if !tag.subtreeQuery {
		panic("subtree query not enabled")
	}
	l, r := tag.tree.Lid[u], tag.tree.Rid[u]
	if root == -1 {
		return tag.bitSubtree.QueryRange(l+tag.edge, r)
	}
	if root == u {
		return tag.bitSubtree.QueryAll()
	}
	if tag.tree.InSubtree(u, root) {
		return tag.bitSubtree.QueryRange(l+tag.edge, r)
	}
	return op(tag.bitSubtree.QueryRange(0, l+1), tag.bitSubtree.QueryRange(r, tag.n))
}

type bitGroup32 struct {
	n     int32
	data  []E
	total E
}

func newBITGroup32(n int32) *bitGroup32 {
	data := make([]E, n)
	for i := range data {
		data[i] = e()
	}
	return &bitGroup32{n: n, data: data, total: e()}
}

func newBITGroup32From(n int32, f func(index int32) E) *bitGroup32 {
	total := e()
	data := make([]E, n)
	for i := range data {
		data[i] = f(int32(i))
		total = op(total, data[i])
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &bitGroup32{n: n, data: data, total: total}
}

func (fw *bitGroup32) Update(i int32, x E) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *bitGroup32) QueryAll() E { return fw.total }

// [0, end)
func (fw *bitGroup32) QueryPrefix(end int32) E {
	if end > fw.n {
		end = fw.n
	}
	res := e()
	for end > 0 {
		res = op(res, fw.data[end-1])
		end &= end - 1
	}
	return res
}

// [start, end)
func (fw *bitGroup32) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	if start > end {
		return e()
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, fw.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
}

type neighbor = struct {
	to   int32
	eid  int32
	cost int
}

type Tree32 struct {
	Lid, Rid      []int32
	IdToNode      []int32
	Depth         []int32
	DepthWeighted []int
	Parent        []int32
	Head          []int32 // 重链头
	Tree          [][]neighbor
	Edges         [][2]int32
	vToE          []int32 // 节点v的父边的id
	n             int32
}

func NewTree32(n int32) *Tree32 {
	res := &Tree32{Tree: make([][]neighbor, n), Edges: make([][2]int32, 0, n-1), n: n}
	return res
}

func (t *Tree32) AddEdge(u, v int32, w int) {
	eid := int32(len(t.Edges))
	t.Tree[u] = append(t.Tree[u], neighbor{to: v, eid: eid, cost: w})
	t.Tree[v] = append(t.Tree[v], neighbor{to: u, eid: eid, cost: w})
	t.Edges = append(t.Edges, [2]int32{u, v})
}

func (t *Tree32) AddDirectedEdge(from, to int32, cost int) {
	eid := int32(len(t.Edges))
	t.Tree[from] = append(t.Tree[from], neighbor{to: to, eid: eid, cost: cost})
	t.Edges = append(t.Edges, [2]int32{from, to})
}

func (t *Tree32) Build(root int32) {
	if root != -1 && int32(len(t.Edges)) != t.n-1 {
		panic("edges count != n-1")
	}
	n := t.n
	t.Lid = make([]int32, n)
	t.Rid = make([]int32, n)
	t.IdToNode = make([]int32, n)
	t.Depth = make([]int32, n)
	t.DepthWeighted = make([]int, n)
	t.Parent = make([]int32, n)
	t.Head = make([]int32, n)
	t.vToE = make([]int32, n)
	for i := int32(0); i < n; i++ {
		t.Depth[i] = -1
		t.Head[i] = root
		t.vToE[i] = -1
	}
	if root != -1 {
		t._dfsSize(root, -1)
		time := int32(0)
		t._dfsHld(root, &time)
	} else {
		time := int32(0)
		for i := int32(0); i < n; i++ {
			if t.Depth[i] == -1 {
				t._dfsSize(i, -1)
				t._dfsHld(i, &time)
			}
		}
	}
}

// 从v开始沿着重链向下收集节点.
func (t *Tree32) HeavyPathAt(v int32) []int32 {
	path := []int32{v}
	for {
		a := path[len(path)-1]
		for _, e := range t.Tree[a] {
			if e.to != t.Parent[a] && t.Head[e.to] == v {
				path = append(path, e.to)
				break
			}
		}
		if path[len(path)-1] == a {
			break
		}
	}
	return path
}

// 返回重儿子，如果没有返回 -1.
func (t *Tree32) HeavyChild(v int32) int32 {
	k := t.Lid[v] + 1
	if k == t.n {
		return -1
	}
	w := t.IdToNode[k]
	if t.Parent[w] == v {
		return w
	}
	return -1
}

// 从v开始向上走k步.
func (t *Tree32) KthAncestor(v, k int32) int32 {
	if k > t.Depth[v] {
		return -1
	}
	for {
		u := t.Head[v]
		if t.Lid[v]-k >= t.Lid[u] {
			return t.IdToNode[t.Lid[v]-k]
		}
		k -= t.Lid[v] - t.Lid[u] + 1
		v = t.Parent[u]
	}
}

func (t *Tree32) Lca(u, v int32) int32 {
	for {
		if t.Lid[u] > t.Lid[v] {
			u, v = v, u
		}
		if t.Head[u] == t.Head[v] {
			return u
		}
		v = t.Parent[t.Head[v]]
	}
}

func (t *Tree32) LcaRooted(u, v, root int32) int32 {
	return t.Lca(u, v) ^ t.Lca(u, root) ^ t.Lca(v, root)
}

func (t *Tree32) Dist(a, b int32) int32 {
	c := t.Lca(a, b)
	return t.Depth[a] + t.Depth[b] - 2*t.Depth[c]
}

func (t *Tree32) DistWeighted(a, b int32) int {
	c := t.Lca(a, b)
	return t.DepthWeighted[a] + t.DepthWeighted[b] - 2*t.DepthWeighted[c]
}

// c 是否在 p 的子树中.c和p不能相等.
func (t *Tree32) InSubtree(c, p int32) bool {
	return t.Lid[p] <= t.Lid[c] && t.Lid[c] < t.Rid[p]
}

// 从 a 开始走 k 步到 b.
func (t *Tree32) Jump(a, b, k int32) int32 {
	if k == 1 {
		if a == b {
			return -1
		}
		if t.InSubtree(b, a) {
			return t.KthAncestor(b, t.Depth[b]-t.Depth[a]-1)
		}
		return t.Parent[a]
	}
	c := t.Lca(a, b)
	dac := t.Depth[a] - t.Depth[c]
	dbc := t.Depth[b] - t.Depth[c]
	if k > dac+dbc {
		return -1
	}
	if k <= dac {
		return t.KthAncestor(a, k)
	}
	return t.KthAncestor(b, dac+dbc-k)
}

func (t *Tree32) SubtreeSize(v int32) int32 {
	return t.Rid[v] - t.Lid[v]
}

func (t *Tree32) SubtreeSizeRooted(v, root int32) int32 {
	if v == root {
		return t.n
	}
	x := t.Jump(v, root, 1)
	if t.InSubtree(v, x) {
		return t.Rid[v] - t.Lid[v]
	}
	return t.n - t.Rid[x] + t.Lid[x]
}

func (t *Tree32) CollectChild(v int32) []int32 {
	var res []int32
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			res = append(res, e.to)
		}
	}
	return res
}

// 收集与 v 相邻的轻边.
func (t *Tree32) CollectLight(v int32) []int32 {
	var res []int32
	skip := true
	for _, e := range t.Tree[v] {
		if e.to != t.Parent[v] {
			if !skip {
				res = append(res, e.to)
			}
			skip = false
		}
	}
	return res
}

func (tree *Tree32) RestorePath(from, to int32) []int32 {
	res := []int32{}
	composition := tree.GetPathDecomposition(from, to, 0)
	for _, e := range composition {
		a, b := e[0], e[1]
		if a <= b {
			for i := a; i <= b; i++ {
				res = append(res, tree.IdToNode[i])
			}
		} else {
			for i := a; i >= b; i-- {
				res = append(res, tree.IdToNode[i])
			}
		}
	}
	return res
}

// 返回沿着`路径顺序`的 [起点,终点] 的 欧拉序 `左闭右闭` 数组.
//
//	!eg:[[2 0] [4 4]] 沿着路径顺序但不一定沿着欧拉序.
func (tree *Tree32) GetPathDecomposition(u, v int32, edge int32) [][2]int32 {
	up, down := [][2]int32{}, [][2]int32{}
	lid, head, parent := tree.Lid, tree.Head, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			down = append(down, [2]int32{lid[head[v]], lid[v]})
			v = parent[head[v]]
		} else {
			up = append(up, [2]int32{lid[u], lid[head[u]]})
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		down = append(down, [2]int32{lid[u] + edge, lid[v]})
	} else if lid[v]+edge <= lid[u] {
		up = append(up, [2]int32{lid[u], lid[v] + edge})
	}
	for i := 0; i < len(down)/2; i++ {
		down[i], down[len(down)-1-i] = down[len(down)-1-i], down[i]
	}
	return append(up, down...)
}

// 遍历路径上的 `[起点,终点)` 欧拉序 `左闭右开` 区间.
func (tree *Tree32) EnumeratePathDecomposition(u, v int32, edge int32, f func(start, end int32)) {
	head, lid, parent := tree.Head, tree.Lid, tree.Parent
	for {
		if head[u] == head[v] {
			break
		}
		if lid[u] < lid[v] {
			a, b := lid[head[v]], lid[v]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			v = parent[head[v]]
		} else {
			a, b := lid[u], lid[head[u]]
			if a > b {
				a, b = b, a
			}
			f(a, b+1)
			u = parent[head[u]]
		}
	}
	if lid[u] < lid[v] {
		a, b := lid[u]+edge, lid[v]
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	} else if lid[v]+edge <= lid[u] {
		a, b := lid[u], lid[v]+edge
		if a > b {
			a, b = b, a
		}
		f(a, b+1)
	}
}

// 返回 root 的欧拉序区间, 左闭右开, 0-indexed.
func (tree *Tree32) Id(root int32) (int32, int32) {
	return tree.Lid[root], tree.Rid[root]
}

// 返回返回边 u-v 对应的 边id.
func (tree *Tree32) Eid(u, v int32) int32 {
	if tree.Parent[u] != v {
		u, v = v, u
	}
	return tree.vToE[u]
}

// 点v对应的父边的边id.如果v是根节点则返回-1.
func (tre *Tree32) VToE(v int32) int32 {
	return tre.vToE[v]
}

// 第i条边对应的深度更深的那个节点.
func (tree *Tree32) EToV(i int32) int32 {
	u, v := tree.Edges[i][0], tree.Edges[i][1]
	if tree.Parent[u] == v {
		return u
	}
	return v
}

func (tree *Tree32) ELid(u int32) int32 {
	return 2*tree.Lid[u] - tree.Depth[u]
}

func (tree *Tree32) ERid(u int32) int32 {
	return 2*tree.Rid[u] - tree.Depth[u] - 1
}

func (t *Tree32) _dfsSize(cur, pre int32) {
	size := t.Rid
	t.Parent[cur] = pre
	if pre != -1 {
		t.Depth[cur] = t.Depth[pre] + 1
	} else {
		t.Depth[cur] = 0
	}
	size[cur] = 1
	nexts := t.Tree[cur]
	for i := int32(len(nexts)) - 2; i >= 0; i-- {
		e := nexts[i+1]
		if t.Depth[e.to] == -1 {
			nexts[i], nexts[i+1] = nexts[i+1], nexts[i]
		}
	}
	hldSize := int32(0)
	for i, e := range nexts {
		to := e.to
		if t.Depth[to] == -1 {
			t.DepthWeighted[to] = t.DepthWeighted[cur] + e.cost
			t.vToE[to] = e.eid
			t._dfsSize(to, cur)
			size[cur] += size[to]
			if size[to] > hldSize {
				hldSize = size[to]
				if i != 0 {
					nexts[0], nexts[i] = nexts[i], nexts[0]
				}
			}
		}
	}
}

func (t *Tree32) _dfsHld(cur int32, times *int32) {
	t.Lid[cur] = *times
	*times++
	t.Rid[cur] += t.Lid[cur]
	t.IdToNode[t.Lid[cur]] = cur
	heavy := true
	for _, e := range t.Tree[cur] {
		to := e.to
		if t.Depth[to] > t.Depth[cur] {
			if heavy {
				t.Head[to] = t.Head[cur]
			} else {
				t.Head[to] = to
			}
			heavy = false
			t._dfsHld(to, times)
		}
	}
}

// 路径 [a,b] 与 [c,d] 的交集.
// 如果为空则返回 {-1,-1}，如果只有一个交点则返回 {x,x}，如果有两个交点则返回 {x,y}.
func (t *Tree32) PathIntersection(a, b, c, d int32) (int32, int32) {
	ab := t.Lca(a, b)
	ac := t.Lca(a, c)
	ad := t.Lca(a, d)
	bc := t.Lca(b, c)
	bd := t.Lca(b, d)
	cd := t.Lca(c, d)
	x := ab ^ ac ^ bc // meet(a,b,c)
	y := ab ^ ad ^ bd // meet(a,b,d)
	if x != y {
		return x, y
	}
	z := ac ^ ad ^ cd
	if x != z {
		x = -1
	}
	return x, x
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
