// 检查边长度限制的路径是否存在
// https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths/solution/zai-xian-zuo-fa-shu-shang-bei-zeng-lca-b-lzjq/

// 1.求出森林中的多个最小生成树是最优的
// 2.求出多个最小生成树后，再求出两个点路径上的最大边权

package main

import (
	"math/bits"
	"sort"
)

func distanceLimitedPathsExist(n int, edgeList [][]int, queries [][]int) []bool {
	forest, _ := Kruskal(n, edgeList)
	uf := NewUnionFindArray(n)
	for _, e := range edgeList {
		uf.Union(e[0], e[1])
	}
	lca := NewLCA(n)
	for i := 0; i < n; i++ {
		for _, e := range forest[i] {
			lca.AddDirectedEdge(i, e.to, e.weight)
		}
	}
	lca.Build(-1)
	res := make([]bool, len(queries))
	for i, q := range queries {
		if uf.IsConnected(q[0], q[1]) {
			res[i] = lca.QueryMaxWeight(q[0], q[1], true) < q[2]
		}
	}
	return res
}

const INF int = 1e18

type LCA struct {
	Tree          [][]edge
	Depth         []int
	DepthWeighted []int
	n             int
	bitLen        int
	dp            [][]int // 节点j向上跳2^i步的父节点
	dpWeight1     [][]int // 节点j向上跳2^i步经过的最大边权
	dpWeight2     [][]int // 节点j向上跳2^i步经过的最小边权
}

type edge struct{ to, weight int }

func NewLCA(n int) *LCA {
	depth := make([]int, n)
	for i := range depth {
		depth[i] = -1
	}
	lca := &LCA{
		Tree:          make([][]edge, n),
		Depth:         depth,
		DepthWeighted: make([]int, n),
		n:             n,
		bitLen:        bits.Len(uint(n)),
	}
	return lca
}

// 添加权值为w的无向边(u, v).
func (lca *LCA) AddEdge(u, v, w int) {
	lca.Tree[u] = append(lca.Tree[u], edge{v, w})
	lca.Tree[v] = append(lca.Tree[v], edge{u, w})
}

// 添加权值为w的有向边(u, v).
func (lca *LCA) AddDirectedEdge(u, v, w int) {
	lca.Tree[u] = append(lca.Tree[u], edge{v, w})
}

// root:0-based
//  当root设为-1时，会从0开始遍历未访问过的连通分量(森林)
func (lca *LCA) Build(root int) {
	lca.dp, lca.dpWeight1, lca.dpWeight2 = makeDp(lca)
	if root != -1 {
		lca.dfsAndInitDp(root, -1, 0, 0)
	} else {
		for i := 0; i < lca.n; i++ {
			if lca.Depth[i] == -1 {
				lca.dfsAndInitDp(i, -1, 0, 0)
			}
		}
	}

	lca.fillDp()
}

// 查询树节点两点的最近公共祖先
func (lca *LCA) QueryLCA(root1, root2 int) int {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, lca.Depth[root2])
	if root1 == root2 {
		return root1
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}
	return lca.dp[0][root1]
}

// 查询树节点两点间距离
//  weighted: 是否将边权计入距离
func (lca *LCA) QueryDist(root1, root2 int, weighted bool) int {
	if weighted {
		return lca.DepthWeighted[root1] + lca.DepthWeighted[root2] - 2*lca.DepthWeighted[lca.QueryLCA(root1, root2)]
	}
	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.QueryLCA(root1, root2)]
}

// 查询树节点两点路径上最大边权(倍增的时候维护其他属性)
//  isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCA) QueryMaxWeight(root1, root2 int, isEdge bool) int {
	res := -INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	toDepth := lca.Depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root1]-toDepth)&(1<<i) > 0 {
			res = max(res, lca.dpWeight1[i][root1])
			root1 = lca.dp[i][root1]
		}
	}
	if root1 == root2 {
		return res
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			res = max(res, max(lca.dpWeight1[i][root1], lca.dpWeight1[i][root2]))
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}
	res = max(res, max(lca.dpWeight1[0][root1], lca.dpWeight1[0][root2]))
	if !isEdge {
		lca_ := lca.dp[0][root1]
		res = max(res, lca.dpWeight1[0][lca_])
	}
	return res
}

// 查询树节点两点路径上最小边权(倍增的时候维护其他属性)
//  isEdge 为true表示查询路径上边权,为false表示查询路径上点权
func (lca *LCA) QueryMinWeight(root1, root2 int, isEdge bool) int {
	res := INF
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	toDepth := lca.Depth[root2]
	for i := lca.bitLen - 1; i >= 0; i-- { // upToDepth
		if (lca.Depth[root1]-toDepth)&(1<<i) > 0 {
			res = min(res, lca.dpWeight2[i][root1])
			root1 = lca.dp[i][root1]
		}
	}
	if root1 == root2 {
		return res
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if lca.dp[i][root1] != lca.dp[i][root2] {
			res = min(res, min(lca.dpWeight2[i][root1], lca.dpWeight2[i][root2]))
			root1 = lca.dp[i][root1]
			root2 = lca.dp[i][root2]
		}
	}
	res = min(res, min(lca.dpWeight2[0][root1], lca.dpWeight2[0][root2]))
	if !isEdge {
		lca_ := lca.dp[0][root1]
		res = min(res, lca.dpWeight2[0][lca_])
	}
	return res
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *LCA) QueryKthAncestor(root, k int) int {
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root = lca.dp[bit][root]
			if root == -1 {
				return -1
			}
		}
		bit++
		k >>= 1
	}
	return root
}

// 从 root 开始向上跳到指定深度 toDepth,toDepth<=dep[v],返回跳到的节点
func (lca *LCA) UpToDepth(root, toDepth int) int {
	if toDepth >= lca.Depth[root] {
		return root
	}
	for i := lca.bitLen - 1; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.dp[i][root]
		}
	}
	return root
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *LCA) Jump(start, target, step int) int {
	lca_ := lca.QueryLCA(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.QueryKthAncestor(start, step)
	}
	return lca.QueryKthAncestor(target, dist-step)
}

func (lca *LCA) dfsAndInitDp(cur, pre, dep, dist int) {
	lca.Depth[cur] = dep
	lca.dp[0][cur] = pre
	lca.DepthWeighted[cur] = dist
	for _, e := range lca.Tree[cur] {
		if next := e.to; next != pre {
			lca.dpWeight1[0][next] = e.weight
			lca.dpWeight2[0][next] = e.weight
			lca.dfsAndInitDp(next, cur, dep+1, dist+e.weight)
		}
	}
}

func makeDp(lca *LCA) (dp, dpWeight1, dpWeight2 [][]int) {
	dp, dpWeight1, dpWeight2 = make([][]int, lca.bitLen), make([][]int, lca.bitLen), make([][]int, lca.bitLen)
	for i := 0; i < lca.bitLen; i++ {
		dp[i], dpWeight1[i], dpWeight2[i] = make([]int, lca.n), make([]int, lca.n), make([]int, lca.n)
		for j := 0; j < lca.n; j++ {
			dp[i][j] = -1
			dpWeight1[i][j] = -INF
			dpWeight2[i][j] = INF
		}
	}
	return
}

func (lca *LCA) fillDp() {
	for i := 0; i < lca.bitLen-1; i++ {
		for j := 0; j < lca.n; j++ {
			if lca.dp[i][j] == -1 {
				lca.dp[i+1][j] = -1
			} else {
				lca.dp[i+1][j] = lca.dp[i][lca.dp[i][j]]
				lca.dpWeight1[i+1][j] = max(lca.dpWeight1[i][j], lca.dpWeight1[i][lca.dp[i][j]])
				lca.dpWeight2[i+1][j] = min(lca.dpWeight2[i][j], lca.dpWeight2[i][lca.dp[i][j]])
			}
		}
	}

	return
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

func maxWithKey(key func(x int) int, args ...int) int {
	max := args[0]
	for _, v := range args[1:] {
		if key(max) < key(v) {
			max = v
		}
	}
	return max
}

type KruskalEdge struct {
	u, v   int
	weight int
}

// !给定无向图的边，求出一个最小生成树(如果不存在,则求出的是森林中的多个最小生成树)
func Kruskal(n int, edges [][]int) (tree [][]edge, ok bool) {
	sortedEdges := make([]KruskalEdge, len(edges))
	for i := range edges {
		e := edges[i]
		sortedEdges[i] = KruskalEdge{u: e[0], v: e[1], weight: e[2]}
	}
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].weight < sortedEdges[j].weight
	})

	tree = make([][]edge, n)
	uf := NewUnionFindArray(n)
	count := 0
	for i := range sortedEdges {
		e := &sortedEdges[i]
		root1, root2 := uf.Find(e.u), uf.Find(e.v)
		if root1 != root2 {
			uf.Union(e.u, e.v)
			tree[e.u] = append(tree[e.u], edge{to: e.v, weight: e.weight})
			tree[e.v] = append(tree[e.v], edge{to: e.u, weight: e.weight})
			count++
			if count == n-1 {
				return tree, true
			}
		}
	}

	return tree, false
}

func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		size:   n,
		rank:   rank,
		parent: parent,
	}
}

type UnionFindArray struct {
	Part   int
	size   int
	rank   []int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	cb(root2, root1)
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}
