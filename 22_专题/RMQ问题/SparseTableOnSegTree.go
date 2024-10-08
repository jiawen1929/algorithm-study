// SparseTableOnSegTree
// 二维rmq/线段树套st表/线段树套rmq
// O(r*c*log(c))构建,O(log(r))查询.
// https://maspypy.github.io/library/ds/sparse_table/sparse_table_on_segtree.hpp

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	yuki866()
	// cf713D()
}

// No.866 レベルKの正方形 (包含k个不同字符的正方形数目)
// https://yukicoder.me/problems/no/866
// 题意: 给定一个 H×W 的地图 G，Gij 为 a-z 的字符。
// 问存在多少个正方形，正方形内的字符种类数恰好为 K。
//
// H,W<=2000.
// !二维查询+二分，枚举左上角，二分边长O(HlogHWlogW).
func yuki866() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W, K int
	fmt.Fscan(in, &H, &W, &K)
	var words []string
	for i := 0; i < H; i++ {
		var word string
		fmt.Fscan(in, &word)
		words = append(words, word)
	}

	grid := make([][]uint32, H)
	for i := 0; i < H; i++ {
		grid[i] = make([]uint32, W)
		for j := 0; j < W; j++ {
			grid[i][j] = 1 << (words[i][j] - 'a')
		}
	}

	res := 0
	seg := NewSparseTableOnSegTreeFrom(grid, func() uint32 { return 0 }, func(a, b uint32) uint32 { return a | b })
	for i := 0; i < H; i++ {
		for j := 0; j < W; j++ {
			maxSize := min(H-1-i, W-1-j)
			v1 := MaxRight(0, func(n int) bool { return bits.OnesCount32(seg.Query(int32(i), int32(i+n), int32(j), int32(j+n))) < K }, maxSize+1)
			v2 := MaxRight(0, func(n int) bool { return bits.OnesCount32(seg.Query(int32(i), int32(i+n), int32(j), int32(j+n))) <= K }, maxSize+1)
			res += v2 - v1
		}
	}

	fmt.Fprintln(out, res)
}

// Animals and Puzzle (区间最大正方形)
// https://www.luogu.com.cn/problem/CF713D
// 题意：给定一个 n×m 的地图 a，ai​ 为 0 或 1。
// 有 q 次询问，每次询问给定一个矩形，求出这个矩形中最大的由 1 构成的正方形的边长是多少。
//
// 1.dp 预处理出每个点为左上角的最大正方形边长.
// 2.线段树套st表,二分答案查询二维区间最大值.
func cf713D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row, col int32
	fmt.Fscan(in, &row, &col)
	grid := make([][]int32, row)
	for i := int32(0); i < row; i++ {
		grid[i] = make([]int32, col)
		for j := int32(0); j < col; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}
	dp := make([][]int32, row)
	for i := int32(0); i < row; i++ {
		dp[i] = make([]int32, col)
	}
	for i := row - 1; i >= 0; i-- {
		for j := col - 1; j >= 0; j-- {
			cur := min32(row-i, col-j)
			if i+1 < row {
				cur = min32(cur, dp[i+1][j]+1)
			}
			if j+1 < col {
				cur = min32(cur, dp[i][j+1]+1)
			}
			if i+1 < row && j+1 < col {
				cur = min32(cur, dp[i+1][j+1]+1)
			}
			if grid[i][j] == 0 {
				cur = 0
			}
			dp[i][j] = cur
		}
	}

	st := NewSparseTableOnSegTreeFrom(dp, func() int32 { return 0 }, max32)

	query := func(rowStart, rowEnd, colStart, colEnd int32) int32 {
		check := func(mid int32) bool {
			return st.Query(rowStart, rowEnd-mid+1, colStart, colEnd-mid+1) >= mid
		}

		left, right := int32(1), min32(rowEnd-rowStart, colEnd-colStart)
		for left <= right {
			mid := (left + right) / 2
			if check(mid) {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
		return right
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var x1, y1, x2, y2 int32
		fmt.Fscan(in, &x1, &y1, &x2, &y2)
		x1--
		y1--
		fmt.Fprintln(out, query(x1, x2, y1, y2))
	}
}

// 3148. 矩阵中的最大得分
// https://leetcode.cn/problems/maximum-difference-score-in-a-grid/description/
// 给你一个由 正整数 组成、大小为 m x n 的矩阵 grid。你可以从矩阵中的任一单元格移动到另一个位于正下方或正右侧的任意单元格（不必相邻）。从值为 c1 的单元格移动到值为 c2 的单元格的得分为 c2 - c1 。
// 你可以从 任一 单元格开始，并且必须至少移动一次。
// 返回你能得到的 最大 总得分。

const INF32 int32 = 1e9 + 10

func maxScore(grid [][]int) int {
	grid32 := make([][]int32, len(grid))
	for i := range grid {
		grid32[i] = make([]int32, len(grid[i]))
		for j := range grid[i] {
			grid32[i][j] = int32(grid[i][j])
		}
	}

	ROW, COL := int32(len(grid32)), int32(len(grid32[0]))
	res := -INF32
	st := NewSparseTableOnSegTreeFrom(grid32, func() int32 { return INF32 }, min32)
	for i := int32(0); i < ROW; i++ {
		for j := int32(0); j < COL; j++ {
			// up
			if i > 0 {
				upMin := st.Query(0, i, j, j+1)
				res = max32(res, grid32[i][j]-upMin)
			}
			// left
			if j > 0 {
				leftMin := st.Query(i, i+1, 0, j)
				res = max32(res, grid32[i][j]-leftMin)
			}

			// left and up
			if i > 0 && j > 0 {
				upLeftMin := st.Query(0, i, 0, j)
				res = max32(res, grid32[i][j]-upLeftMin)
			}
		}
	}

	return int(res)
}

// 更快的 SparseTable2DFast.
type SparseTableOnSegTree[E any] struct {
	row, col int32
	e        func() E
	op       func(E, E) E
	data     []*SparseTable[E]
}

func NewSparseTableOnSegTreeFrom[E any](grid [][]E, e func() E, op func(E, E) E) *SparseTableOnSegTree[E] {
	row := int32(len(grid))
	col := int32(0)
	if row > 0 {
		col = int32(len(grid[0]))
	}
	data := make([]*SparseTable[E], 2*row)
	for i := int32(0); i < row; i++ {
		data[row+i] = NewSparseTableFrom(grid[i], e, op)
	}
	for i := row - 1; i > 0; i-- {
		data[i] = NewSparseTable(
			col,
			func(j int32) E {
				x := data[2*i].Query(j, j+1)
				y := data[2*i+1].Query(j, j+1)
				return op(x, y)
			},
			e, op,
		)
	}
	return &SparseTableOnSegTree[E]{row: row, col: col, e: e, op: op, data: data}
}

func (st *SparseTableOnSegTree[E]) Query(rowStart, rowEnd, colStart, colEnd int32) E {
	if !(0 <= rowStart && rowStart <= rowEnd && rowEnd <= st.row) {
		return st.e()
	}
	if !(0 <= colStart && colStart <= colEnd && colEnd <= st.col) {
		return st.e()
	}
	res := st.e()
	rowStart += st.row
	rowEnd += st.row
	for rowStart < rowEnd {
		if rowStart&1 != 0 {
			res = st.op(res, st.data[rowStart].Query(colStart, colEnd))
			rowStart++
		}
		if rowEnd&1 != 0 {
			rowEnd--
			res = st.op(res, st.data[rowEnd].Query(colStart, colEnd))
		}
		rowStart >>= 1
		rowEnd >>= 1
	}
	return res
}

type SparseTable[E any] struct {
	st [][]E
	e  func() E
	op func(E, E) E
	n  int32
}

func NewSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *SparseTable[E] {
	res := &SparseTable[E]{}

	b := int32(bits.Len32(uint32(n)))
	st := make([][]E, b)
	for i := range st {
		st[i] = make([]E, n)
	}
	for i := int32(0); i < n; i++ {
		st[0][i] = f(i)
	}
	for i := int32(1); i < b; i++ {
		for j := int32(0); j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	res.st = st
	res.e = e
	res.op = op
	res.n = n
	return res
}

func NewSparseTableFrom[E any](leaves []E, e func() E, op func(E, E) E) *SparseTable[E] {
	return NewSparseTable(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	b := int32(bits.Len32(uint32(end-start))) - 1
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func max32(a, b int32) int32 {
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
