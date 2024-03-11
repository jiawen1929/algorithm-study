// !可修改点权的 Link-Cut Tree
// https://ei1333.github.io/library/structure/lct/link-cut-tree-lazy-path.hpp
// https://hitonanode.github.io/cplib-cpp/data_structure/link_cut_tree.hpp

// LinkCutTreeLazy(rev,id,op,mapping,composition): コンストラクタ.
//   rev は要素を反転する演算を指す.
//   id は作用素の単位元を指す.
//   op は 2 つの要素の値をマージする二項演算,
//   mapping は要素と作用素をマージする二項演算,
//   composition は作用素同士をマージする二項演算,
// Build(vs): 各要素の値を vs[i] としたノードを生成し, その配列を返す.
// Evert(t): t を根に変更する.
// LinkEdge(child, parent): child の親を parent にする.如果已经连通则不进行操作。
// CutEdge(u,v) : u と v の間の辺を切り離す.如果边不存在则不进行操作。
// QueryToRoot(u): u から根までのパス上の頂点の値を二項演算でまとめた結果を返す.
// QueryPath(u, v): u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
// QeuryKthAncestor(x, k): x から根までのパスに出現するノードを並べたとき, 0-indexed で k 番目のノードを返す.
// QueryLCA(u, v): u と v の lca を返す. u と v が異なる連結成分なら nullptr を返す.
//   !上記の操作は根を勝手に変えるため、根を固定したい場合は Evert で根を固定してから操作する.
// IsConnected(u, v): u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
// Alloc(v): 要素の値を v としたノードを生成する.
// UpdateToRoot(t, lazy): t から根までのパス上の頂点に作用素 lazy を加える.
// Update(u, v, lazy): u から v までのパス上の頂点に作用素 lazy を加える.
// Set(t, v): t の値を v に変更する.
// Get(t): t の値を返す.
// GetRoot(t): t の根を返す.
// expose(t): t と根をつなげて, t を splay Tree の根にする.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	// https://judge.yosupo.jp/problem/dynamic_tree_vertex_add_path_sum
	// 0 u v w x  删除 u-v 边, 添加 w-x 边
	// 1 p x  将 p 节点的值加上 x
	// 2 u v  输出 u-v 路径上所有点的值的和(包含端点)
	// n,q=2e5
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	lct := NewLinkCutTreeLazy(true)
	vs := lct.Build(nums)
	for i := 0; i < n-1; i++ { // 连接树边
		var v1, v2 int
		fmt.Fscan(in, &v1, &v2)
		// lct.Evert(vs[v1])
		// lct.Link(vs[v1], vs[v2])
		lct.LinkEdge(vs[v1], vs[v2])
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var u, v, w, x int
			fmt.Fscan(in, &u, &v, &w, &x)
			// lct.Evert(vs[u])
			// lct.Cut(vs[v])
			lct.CutEdge(vs[u], vs[v])
			// lct.Evert(vs[w])
			// lct.Link(vs[w], vs[x])
			lct.LinkEdge(vs[w], vs[x])
		} else if op == 1 {
			var p, x int
			fmt.Fscan(in, &p, &x)
			// lct.Set(vs[p], vs[p].key+x)
			lct.UpdatePath(vs[p], vs[p], x)
		} else {
			var u, v int
			fmt.Fscan(in, &u, &v)
			fmt.Fprintln(out, lct.QueryPath(vs[u], vs[v]))
		}
	}

}

type E = int
type Id = int

// 区间反转
func (*LinkCutTreeLazy) rev(e E) E                               { return e }
func (*LinkCutTreeLazy) id() Id                                  { return 0 }
func (*LinkCutTreeLazy) op(a, b E) E                             { return a + b }
func (*LinkCutTreeLazy) mapping(lazy Id, data E) E               { return data + lazy }
func (*LinkCutTreeLazy) composition(parentLazy, childLazy Id) Id { return parentLazy + childLazy }

type LinkCutTreeLazy struct {
	nodeId int
	edges  map[struct{ u, v int }]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTreeLazy(check bool) *LinkCutTreeLazy {
	return &LinkCutTreeLazy{edges: make(map[struct{ u, v int }]struct{}), check: check}
}

// 各要素の値を vs[i] としたノードを生成し, その配列を返す.
func (lct *LinkCutTreeLazy) Build(vs []E) []*treeNode {
	nodes := make([]*treeNode, len(vs))
	for i, v := range vs {
		nodes[i] = lct.Alloc(v)
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeLazy) Alloc(e E) *treeNode {
	res := newTreeNode(e, lct.id(), lct.nodeId)
	lct.nodeId++
	return res
}

// t を根に変更する.
func (lct *LinkCutTreeLazy) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

// 存在していない辺 uv を新たに張る.
//  すでに存在している辺 uv に対しては何もしない.
func (lct *LinkCutTreeLazy) LinkEdge(child, parent *treeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := struct{ u, v int }{id1, id2}
		lct.edges[tuple] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

// 存在している辺を切り離す.
//  存在していない辺に対しては何もしない.
func (lct *LinkCutTreeLazy) CutEdge(u, v *treeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := struct{ u, v int }{id1, id2}
		if _, has := lct.edges[tuple]; !has {
			return
		}
		delete(lct.edges, tuple)
	}

	lct.Evert(u)
	lct.expose(v)
	parent := v.l
	v.l = nil
	lct.update(v)
	parent.p = nil
	return true
}

// u と v の lca を返す.
//  u と v が異なる連結成分なら nullptr を返す.
//  !上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTreeLazy) QueryLCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeLazy) QueryKthAncestor(x *treeNode, k int) *treeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.sz > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.sz
			}
			if k == 0 {
				return x
			}
			k--
			x = x.l
		}
	}
	return nil
}

// u から根までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTreeLazy) QueryToRoot(u *treeNode) E {
	lct.expose(u)
	return u.sum
}

// u から v までのパス上の頂点の値を二項演算でまとめた結果を返す.
func (lct *LinkCutTreeLazy) QueryPath(u, v *treeNode) E {
	lct.Evert(u)
	return lct.QueryToRoot(v)
}

//  t から根までのパス上の頂点に作用素 lazy を加える.
func (lct *LinkCutTreeLazy) UpdateToRoot(t *treeNode, lazy Id) {
	lct.expose(t)
	lct.propagate(t, lazy)
	lct.push(t)
}

// u から v までのパス上の頂点に作用素 lazy を加える.
func (lct *LinkCutTreeLazy) UpdatePath(u, v *treeNode, lazy Id) {
	lct.Evert(u)
	lct.UpdateToRoot(v, lazy)
}

// t の値を v に変更する.
func (lct *LinkCutTreeLazy) Set(t *treeNode, v E) {
	lct.expose(t)
	t.key = v
	lct.update(t)
}

// t の値を返す.
func (lct *LinkCutTreeLazy) Get(t *treeNode) E {
	return t.key
}

func (lct *LinkCutTreeLazy) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeLazy) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeLazy) expose(t *treeNode) *treeNode {
	rp := (*treeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		cur.r = rp
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeLazy) update(t *treeNode) *treeNode {
	t.sz = 1
	t.sum = t.key
	if t.l != nil {
		t.sz += t.l.sz
		t.sum = lct.op(t.l.sum, t.sum)
	}
	if t.r != nil {
		t.sz += t.r.sz
		t.sum = lct.op(t.sum, t.r.sum)
	}
	return t
}

func (lct *LinkCutTreeLazy) rotr(t *treeNode) {
	x := t.p
	y := x.p
	x.l = t.r
	if t.r != nil {
		t.r.p = x
	}
	t.r = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTreeLazy) rotl(t *treeNode) {
	x := t.p
	y := x.p
	x.r = t.l
	if t.l != nil {
		t.l.p = x
	}
	t.l = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTreeLazy) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.sum = lct.rev(t.sum)
	t.rev = !t.rev
}

func (lct *LinkCutTreeLazy) propagate(t *treeNode, lazy Id) {
	t.lazy = lct.composition(lazy, t.lazy)
	t.key = lct.mapping(lazy, t.key)
	t.sum = lct.mapping(lazy, t.sum)
}

func (lct *LinkCutTreeLazy) push(t *treeNode) {
	if t.lazy != lct.id() {
		if t.l != nil {
			lct.propagate(t.l, t.lazy)
		}
		if t.r != nil {
			lct.propagate(t.r, t.lazy)
		}
		t.lazy = lct.id()
	}

	if t.rev {
		if t.l != nil {
			lct.toggle(t.l)
		}
		if t.r != nil {
			lct.toggle(t.r)
		}
		t.rev = false
	}
}

func (lct *LinkCutTreeLazy) splay(t *treeNode) {
	lct.push(t)
	for !t.IsRoot() {
		q := t.p
		if q.IsRoot() {
			lct.push(q)
			lct.push(t)
			if q.l == t {
				lct.rotr(t)
			} else {
				lct.rotl(t)
			}
		} else {
			r := q.p
			lct.push(r)
			lct.push(q)
			lct.push(t)
			if r.l == q {
				if q.l == t {
					lct.rotr(q)
					lct.rotr(t)
				} else {
					lct.rotl(t)
					lct.rotr(t)
				}
			} else {
				if q.r == t {
					lct.rotl(q)
					lct.rotl(t)
				} else {
					lct.rotr(t)
					lct.rotl(t)
				}
			}
		}
	}
}

type treeNode struct {
	l, r, p  *treeNode
	key, sum E
	lazy     Id
	rev      bool
	sz       int
	id       int
}

func newTreeNode(v E, lazy Id, id int) *treeNode {
	return &treeNode{key: v, sum: v, lazy: lazy, sz: 1, id: id}
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) String() string {
	return fmt.Sprintf("key: %v, sum: %v, sz: %v, rev: %v", n.key, n.sum, n.sz, n.rev)
}
