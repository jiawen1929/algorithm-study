// DividePathOnTreeBinaryLift/DoublingLca
// 倍增拆树上路径`path(from,to)`：倍增拆点将树上的一段路径拆成logn个点
//
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// jump()
	test()
}

func jump() {
	// https://judge.yosupo.jp/problem/jump_on_tree
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)

	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}
	D := NewDoublingLca32Dynamic(n)
	var dfs func(int32, int32)
	dfs = func(node, parent int32) {
		for _, child := range tree[node] {
			if child != parent {
				D.AddLeaf(child, node)
				dfs(child, node)
			}
		}
	}
	dfs(0, -1)

	for i := int32(0); i < q; i++ {
		var from, to, k int32
		fmt.Fscan(in, &from, &to, &k)
		fmt.Fprintln(out, D.Jump(from, to, k))
	}
}

const INF int = 1e18

type DoublingLca32Dynamic struct {
	Depth []int32

	n, log, size int32
	jump         [][]int32 // 节点j向上跳2^i步的父节点
}

// 不预先给出整棵树,而是动态添加叶子节点.
// 初始时只有一个根节点root.
func NewDoublingLca32Dynamic(n int32) *DoublingLca32Dynamic {
	depth := make([]int32, n)
	log := int32(bits.Len32(uint32(n))) - 1
	size := n * (log + 1)
	jump := make([][]int32, log+1)
	for i := int32(0); i <= log; i++ {
		nums := make([]int32, n)
		for j := int32(0); j < n; j++ {
			nums[j] = -1
		}
		jump[i] = nums
	}
	return &DoublingLca32Dynamic{Depth: depth, n: n, log: log, size: size, jump: jump}
}

// 在树中添加一条从parent到leaf的边，要求parent已经存在于树中(或者为根节点)，且leaf不在树中.
func (lca *DoublingLca32Dynamic) AddLeaf(leaf int32, parent int32) {
	lca.Depth[leaf] = lca.Depth[parent] + 1
	lca.jump[0][leaf] = parent
	for i := int32(1); i <= lca.log; i++ {
		if from := lca.jump[i-1][leaf]; from != -1 {
			lca.jump[i][leaf] = lca.jump[i-1][from]
		} else {
			break
		}
	}
}

// 遍历路径(start,target)上的所有jump.
// 倍增拆点，将树上的一段路径拆成logn个点.
func (lca *DoublingLca32Dynamic) EnumerateJump(start, target int32, f func(level, index int32)) {
	if start == target {
		f(0, start)
		return
	}
	if lca.Depth[start] < lca.Depth[target] {
		start, target = target, start
	}
	toDepth := lca.Depth[target]
	if lca.Depth[start] > toDepth {
		for i := lca.log; i >= 0; i-- {
			if (lca.Depth[start]-toDepth)&(1<<i) > 0 {
				f(i, start)
				start = lca.jump[i][start]
			}
		}
	}
	if start == target {
		f(0, start)
		return
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.jump[i][start], lca.jump[i][target]; a != b {
			f(i, start)
			f(i, target)
			start, target = a, b
		}
	}
	f(0, start)
	f(0, target)
	f(0, lca.jump[0][start])
}

// 遍历路径(start,target)上的所有jump.
// !要求运算幂等(idempotent).
func (lca *DoublingLca32Dynamic) EnumerateJumpDangerously(start, target int32, f func(level, index int32)) {
	if start == target {
		f(0, start)
		return
	}

	divide := func(node, ancestor int32, f func(level, index int32)) {
		len_ := lca.Depth[node] - lca.Depth[ancestor] + 1
		k := int32(bits.Len32(uint32(len_))) - 1
		jumpLen := len_ - (1 << k)
		from2 := lca.KthAncestor(node, jumpLen)
		f(k, node)
		f(k, from2)
	}

	if lca.Depth[start] < lca.Depth[target] {
		start, target = target, start
	}
	lca_ := lca.Lca(start, target)
	if lca_ == target {
		divide(start, lca_, f)
	} else {
		divide(start, lca_, f)
		divide(target, lca_, f)
	}
}

// 下推路径信息，更新答案.
// O(n*log(n)).
func (lca *DoublingLca32Dynamic) PushDown(f func(pLevel, pIndex int32, cLevel, cIndex1, cIndex2 int32)) {
	n, log := lca.n, lca.log
	for k := log - 1; k >= 0; k-- {
		for i := int32(0); i < n; i++ {
			// push down jump(i,k+1) to jump(i,k) and jump(jump(i,k),k)
			if to := lca.jump[k][i]; to != -1 {
				f(k+1, i, k, i, to)
			}
		}
	}
}

func (lca *DoublingLca32Dynamic) Size() int32 { return lca.size }
func (lca *DoublingLca32Dynamic) Log() int32  { return lca.log }

func (lca *DoublingLca32Dynamic) Lca(root1, root2 int32) int32 {
	if lca.Depth[root1] < lca.Depth[root2] {
		root1, root2 = root2, root1
	}
	root1 = lca.UpToDepth(root1, lca.Depth[root2])
	if root1 == root2 {
		return root1
	}
	for i := lca.log; i >= 0; i-- {
		if a, b := lca.jump[i][root1], lca.jump[i][root2]; a != b {
			root1, root2 = a, b
		}
	}
	return lca.jump[0][root1]
}

func (lca *DoublingLca32Dynamic) Dist(root1, root2 int32, weighted bool) int32 {
	return lca.Depth[root1] + lca.Depth[root2] - 2*lca.Depth[lca.Lca(root1, root2)]
}

// 查询树节点root的第k个祖先(向上跳k步),如果不存在这样的祖先节点,返回 -1
func (lca *DoublingLca32Dynamic) KthAncestor(root, k int32) int32 {
	if k > lca.Depth[root] {
		return -1
	}
	bit := 0
	for k > 0 {
		if k&1 == 1 {
			root = lca.jump[bit][root]
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
func (lca *DoublingLca32Dynamic) UpToDepth(root, toDepth int32) int32 {
	if toDepth >= lca.Depth[root] {
		return root
	}
	for i := lca.log; i >= 0; i-- {
		if (lca.Depth[root]-toDepth)&(1<<i) > 0 {
			root = lca.jump[i][root]
		}
	}
	return root
}

// 从start节点跳向target节点,跳过step个节点(0-indexed)
// 返回跳到的节点,如果不存在这样的节点,返回-1
func (lca *DoublingLca32Dynamic) Jump(start, target, step int32) int32 {
	lca_ := lca.Lca(start, target)
	dep1, dep2, deplca := lca.Depth[start], lca.Depth[target], lca.Depth[lca_]
	dist := dep1 + dep2 - 2*deplca
	if step > dist {
		return -1
	}
	if step <= dep1-deplca {
		return lca.KthAncestor(start, step)
	}
	return lca.KthAncestor(target, dist-step)
}

func (lca *DoublingLca32Dynamic) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	// `LastTrue`判定条件取反，然后向上跳一层.
	if predicate(start) {
		return start
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp != -1 && !predicate(tmp) {
			start = tmp
		}
	}
	return lca.jump[0][start] // 不存在则返回-1
}

func (lca *DoublingLca32Dynamic) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for k := lca.log; k >= 0; k-- {
		tmp := lca.jump[k][start]
		if tmp != -1 && predicate(tmp) {
			start = tmp
		}
	}
	return start
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

func test() {

	//          0
	//        /   \
	//       1     2
	//      / \     \
	//     3   4     5
	//         /
	//        6

	n := int32(7)
	edges := [][]int32{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {4, 6}}
	tree := make([][]int32, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	lca := NewDoublingLca32Dynamic(n)
	var dfs func(int32, int32)
	dfs = func(node, parent int32) {
		for _, child := range tree[node] {
			if child != parent {
				lca.AddLeaf(child, node)
				dfs(child, node)
			}
		}
	}
	dfs(0, -1)

	fmt.Println(lca.Lca(3, 6)) // 1

	to1 := lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= -1 })
	expect[int32](to1, -1)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 1 })
	expect[int32](to1, 1)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 2 })
	expect[int32](to1, 4)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 3 })
	expect[int32](to1, 6)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] <= 4 })
	expect[int32](to1, 6)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] >= 2 })
	expect[int32](to1, 6)
	to1 = lca.FirstTrue(6, func(i int32) bool { return lca.Depth[i] >= 4 })
	expect[int32](to1, -1)

	to2 := lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 1 })
	expect[int32](to2, 1)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 2 })
	expect[int32](to2, 4)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 3 })
	expect[int32](to2, 6)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] >= 4 })
	expect[int32](to2, -1)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= 2 })
	expect[int32](to2, -1)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= 4 })
	expect[int32](to2, 0)
	to2 = lca.LastTrue(6, func(i int32) bool { return lca.Depth[i] <= -1 })
	expect[int32](to2, -1)

	values := make([]int32, lca.Size())
	lca.EnumerateJump(6, 3, func(level, index int32) {
		values[level*lca.n+index] += 2
	})
	lca.EnumerateJump(3, 5, func(level, index int32) {
		values[level*lca.n+index] += 3
	})
	lca.EnumerateJump(4, 4, func(level, index int32) {
		values[level*lca.n+index] += 5
	})
	lca.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*lca.n+pIndex, cLevel*lca.n+cIndex1, cLevel*lca.n+cIndex2
		values[c1] += values[p]
		values[c2] += values[p]
	})

	expected := []int32{3, 5, 3, 5, 7, 3, 2}
	for i := int32(0); i < n; i++ {
		expect[int32](values[i], expected[i])
	}

	// test EnumerateJumpDangerously by idempoent function
	values = make([]int32, lca.Size())
	lca.EnumerateJumpDangerously(6, 3, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 2)
	})
	lca.EnumerateJumpDangerously(3, 5, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 3)
	})
	lca.EnumerateJumpDangerously(4, 4, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 5)
	})
	lca.EnumerateJumpDangerously(1, 6, func(level, index int32) {
		id := level*lca.n + index
		values[id] = max32(values[id], 7)
	})

	lca.PushDown(func(pLevel, pIndex, cLevel, cIndex1, cIndex2 int32) {
		p, c1, c2 := pLevel*lca.n+pIndex, cLevel*lca.n+cIndex1, cLevel*lca.n+cIndex2
		values[c1] = max32(values[c1], values[p])
		values[c2] = max32(values[c2], values[p])
	})

	expected = []int32{3, 7, 3, 3, 7, 3, 7}
	for i := int32(0); i < n; i++ {
		expect[int32](values[i], expected[i])
	}

	fmt.Println("test passed")
}

func expect[S comparable](actual, expected S) {
	if actual != expected {
		panic(fmt.Sprintf("actual: %v, expected: %v", actual, expected))
	}
}
