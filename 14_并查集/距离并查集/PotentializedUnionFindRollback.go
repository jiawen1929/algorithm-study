package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

// UnionfindwithPotential
// https://judge.yosupo.jp/problem/unionfind_with_potential
// 0 u v x: 判断A[u]=A[v]+x(mod Mod)是否成立. 如果与现有信息矛盾,则不进行任何操作,否则将该条件加入.
// 1 u v: 输出A[u]-A[v].如果不能确定,输出-1.
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, q int32
	fmt.Fscan(in, &n, &q)

	e := func() int { return 0 }
	op := func(a, b int) int {
		res := (a + b) % MOD
		if res < 0 {
			res += MOD
		}
		return res
	}
	inv := func(a int) int {
		res := -a % MOD
		if res < 0 {
			res += MOD
		}
		return res
	}

	uf := NewPotentializedUnionFindRollback(n, e, op, inv)

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var u, v int32
			var x int
			fmt.Fscan(in, &u, &v, &x)
			diff, same := uf.Diff(u, v)
			valid := !same || diff == x
			if valid {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
			if !same {
				uf.Union(u, v, x)
			}
		} else {
			var u, v int32
			fmt.Fscan(in, &u, &v)
			if diff, ok := uf.Diff(u, v); ok {
				fmt.Fprintln(out, diff)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}

type item[E any] struct {
	root int32
	diff E
}

// 可撤销势能并查集/距离并查集.
type PotentializedUnionFindRollback[E comparable] struct {
	data *RollbackArray[item[E]]
	e    func() E
	op   func(E, E) E
	inv  func(E) E
}

func NewPotentializedUnionFindRollback[E comparable](n int32, e func() E, op func(E, E) E, inv func(E) E) *PotentializedUnionFindRollback[E] {
	initData := make([]item[E], n)
	for i := int32(0); i < n; i++ {
		initData[i] = item[E]{root: -1, diff: e()}
	}
	return &PotentializedUnionFindRollback[E]{
		data: NewRollbackArrayFrom(initData),
		e:    e, op: op, inv: inv,
	}
}

func (uf *PotentializedUnionFindRollback[E]) GetTime() int32 {
	return uf.data.GetTime()
}

func (uf *PotentializedUnionFindRollback[E]) Rollback(time int32) {
	uf.data.Rollback(time)
}

// P[a] - P[b] = x
func (uf *PotentializedUnionFindRollback[E]) Union(a, b int32, x E) bool {
	v1, x1 := uf.Find(a)
	v2, x2 := uf.Find(b)
	if v1 == v2 {
		return false
	}
	item1, item2 := uf.data.Get(v1), uf.data.Get(v2)
	s1, s2 := -item1.root, -item2.root
	if s1 < s2 {
		s1, s2 = s2, s1
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		x = uf.inv(x)
	}
	x = uf.op(x1, x)
	x = uf.op(x, uf.inv(x2))
	uf.data.Set(v2, item[E]{root: v1, diff: x})
	uf.data.Set(v1, item[E]{root: -(s1 + s2), diff: uf.e()})
	return true
}

func (uf *PotentializedUnionFindRollback[E]) Find(v int32) (root int32, dist E) {
	dist = uf.e()
	for {
		item := uf.data.Get(v)
		if item.root < 0 {
			break
		}
		dist = uf.op(item.diff, dist)
		v = item.root
	}
	root = v
	return
}

// Diff(a, b) = P[a] - P[b]
func (uf *PotentializedUnionFindRollback[E]) Diff(a, b int32) (E, bool) {
	ru, xu := uf.Find(a)
	rv, xv := uf.Find(b)
	if ru != rv {
		return uf.e(), false
	}
	return uf.op(uf.inv(xu), xv), true
}

type HistoryItem[V comparable] struct {
	index int32
	value V
}

type RollbackArray[V comparable] struct {
	n       int32
	data    []V
	history []HistoryItem[V]
}

func NewRollbackArray[V comparable](n int32, f func(i int32) V) *RollbackArray[V] {
	data := make([]V, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray[V]{
		n:    n,
		data: data,
	}
}

func NewRollbackArrayFrom[V comparable](data []V) *RollbackArray[V] {
	return &RollbackArray[V]{n: int32(len(data)), data: data}
}

func (r *RollbackArray[V]) GetTime() int32 {
	return int32(len(r.history))
}

func (r *RollbackArray[V]) Rollback(time int32) {
	for i := int32(len(r.history)) - 1; i >= time; i-- {
		pair := r.history[i]
		r.data[pair.index] = pair.value
	}
	r.history = r.history[:time]
}

func (r *RollbackArray[V]) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray[V]) Get(index int32) V {
	return r.data[index]
}

func (r *RollbackArray[V]) Set(index int32, value V) bool {
	if r.data[index] == value {
		return false
	}
	r.history = append(r.history, HistoryItem[V]{index: index, value: r.data[index]})
	r.data[index] = value
	return true
}

func (r *RollbackArray[V]) GetAll() []V {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray[V]) Len() int32 {
	return r.n
}
