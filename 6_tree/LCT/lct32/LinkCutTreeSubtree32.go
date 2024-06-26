// 动态子树和查询
// n,q=2e5
// !子树和还可以用来查询联通分量大小

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/dynamic_tree_vertex_add_subtree_sum
	// 0 u v w x  删除 u-v 边, 添加 w-x 边
	// 1 p x  将 p 节点的值加上 x
	// 2 u p  对于边(u,p) 查询结点v的子树的和,p为u的父节点

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	lct := NewLinkCutTreeSubTree(true)
	vs := lct.Build(int32(n), func(i int32) E { return int(nums[i]) })

	for i := 0; i < n-1; i++ { // 连接树边
		var v1, v2 int32
		fmt.Fscan(in, &v1, &v2)
		lct.LinkEdge(vs[v1], vs[v2])
	}

	for i := 0; i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var u, v, w, x int32
			fmt.Fscan(in, &u, &v, &w, &x)
			lct.CutEdge(vs[u], vs[v])
			lct.LinkEdge(vs[w], vs[x])
		} else if op == 1 {
			var root, delta int32
			fmt.Fscan(in, &root, &delta)
			lct.Set(vs[root], lct.Get(vs[root])+int(delta))
		} else {
			var root, parent int32
			fmt.Fscan(in, &root, &parent)
			lct.Evert(vs[parent]) // !注意查询子树前要先把父节点旋转到根节点
			fmt.Fprintln(out, lct.QuerySubTree(vs[root]))
		}
	}
}

type E = int // 子树和

func (*treeNode) e() E                { return 0 }
func (*treeNode) op(this, other E) E  { return this + other }
func (*treeNode) inv(this, other E) E { return this - other }

type LinkCutTreeSubTree struct {
	nodeId int32
	edges  map[[2]int32]struct{}
	check  bool
}

// check: AddEdge/RemoveEdge で辺の存在チェックを行うかどうか.
func NewLinkCutTreeSubTree(check bool) *LinkCutTreeSubTree {
	return &LinkCutTreeSubTree{edges: make(map[[2]int32]struct{}), check: check}
}

func (lct *LinkCutTreeSubTree) Build(n int32, f func(i int32) E) []*treeNode {
	nodes := make([]*treeNode, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = lct.Alloc(f(i))
	}
	return nodes
}

// 要素の値を v としたノードを生成する.
func (lct *LinkCutTreeSubTree) Alloc(key E) *treeNode {
	res := newTreeNode(key, lct.nodeId)
	lct.nodeId++
	lct.update(res)
	return res
}

// t を根に変更する.
func (lct *LinkCutTreeSubTree) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

func (lct *LinkCutTreeSubTree) LinkEdge(child, parent *treeNode) (ok bool) {
	if lct.check {
		if lct.IsConnected(child, parent) {
			return
		}
		id1, id2 := child.id, parent.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := [2]int32{id1, id2}
		lct.edges[tuple] = struct{}{}
	}

	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
	return true
}

func (lct *LinkCutTreeSubTree) CutEdge(u, v *treeNode) (ok bool) {
	if lct.check {
		id1, id2 := u.id, v.id
		if id1 > id2 {
			id1, id2 = id2, id1
		}
		tuple := [2]int32{id1, id2}
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
//
//	u と v が異なる連結成分なら nullptr を返す.
//	!上記の操作は根を勝手に変えるため, 事前に Evert する必要があるかも.
func (lct *LinkCutTreeSubTree) QueryLCA(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeSubTree) KthAncestor(x *treeNode, k int32) *treeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.cnt > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.cnt
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

func (lct *LinkCutTreeSubTree) GetParent(t *treeNode) *treeNode {
	lct.expose(t)
	p := t.l
	if p == nil {
		return nil
	}
	for {
		lct.push(p)
		if p.r == nil {
			return p
		}
		p = p.r
	}
}

func (lct *LinkCutTreeSubTree) Jump(from, to *treeNode, k int32) *treeNode {
	lct.Evert(to)
	lct.expose(from)
	for {
		lct.push(from)
		rs := int32(0)
		if from.r != nil {
			rs = from.r.cnt
		}
		if k < rs {
			from = from.r
			continue
		}
		if k == rs {
			break
		}
		k -= rs + 1
		from = from.l
	}
	lct.splay(from)
	return from
}

func (lct *LinkCutTreeSubTree) Dist(u, v *treeNode) int32 {
	lct.Evert(u)
	lct.expose(v)
	return v.cnt - 1
}

// t を根とする部分木の要素の値の和を返す.
//
//	!Evert を忘れない！
func (lct *LinkCutTreeSubTree) QuerySubTree(t *treeNode) E {
	lct.expose(t)
	return t.op(t.key, t.sub)
}

// t の値を v に変更する.
func (lct *LinkCutTreeSubTree) Set(t *treeNode, key E) *treeNode {
	lct.expose(t)
	t.key = key
	lct.update(t)
	return t
}

// t の値を返す.
func (lct *LinkCutTreeSubTree) Get(t *treeNode) E {
	return t.key
}

// u と v が同じ連結成分に属する場合は true, そうでなければ false を返す.
func (lct *LinkCutTreeSubTree) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeSubTree) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

func (lct *LinkCutTreeSubTree) expose(t *treeNode) *treeNode {
	rp := (*treeNode)(nil)
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		if cur.r != nil {
			cur.Add(cur.r)
		}
		cur.r = rp
		if cur.r != nil {
			cur.Erase(cur.r)
		}
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeSubTree) update(t *treeNode) {
	t.cnt = 1
	if t.l != nil {
		t.cnt += t.l.cnt
	}
	if t.r != nil {
		t.cnt += t.r.cnt
	}

	t.Merge(t.l, t.r)
}

func (lct *LinkCutTreeSubTree) rotr(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) rotl(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.rev = !t.rev
}

func (lct *LinkCutTreeSubTree) push(t *treeNode) {
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

func (lct *LinkCutTreeSubTree) splay(t *treeNode) {
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
	key, sum, sub E
	rev           bool
	cnt           int32
	id            int32
	l, r, p       *treeNode
}

func newTreeNode(key E, id int32) *treeNode {
	res := &treeNode{key: key, sum: key, cnt: 1, id: id}
	res.sub = res.e()
	return res
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) Add(other *treeNode)   { n.sub = n.op(n.sub, other.sum) }
func (n *treeNode) Erase(other *treeNode) { n.sub = n.inv(n.sub, other.sum) }
func (n *treeNode) Merge(n1, n2 *treeNode) {
	var tmp1, tmp2 E
	if n1 != nil {
		tmp1 = n1.sum
	} else {
		tmp1 = n.e()
	}

	if n2 != nil {
		tmp2 = n2.sum
	} else {
		tmp2 = n.e()
	}

	n.sum = n.op(n.op(tmp1, n.key), n.op(n.sub, tmp2))
}
