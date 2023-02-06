// 带有上下界的最大流
// https://ei1333.github.io/library/graph/flow/maxflow-lower-bound.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, START, END int
	fmt.Fscan(in, &n, &m, &START, &END)
	START--
	END--
	M := NewMaxFlowLowerBound(n)
	for i := 0; i < m; i++ {
		var u, v, lower, upper int
		fmt.Fscan(in, &u, &v, &lower, &upper)
		u--
		v--
		M.AddEdge(u, v, lower, upper)
	}

	res := M.MinFlow(START, END)
	if res == -1 {
		fmt.Fprintln(out, "No Solution")
	} else {
		fmt.Fprintln(out, res)
	}
}

const INF int = 1e18

type MaxFlowLowerBound struct {
	n            int
	sum          int
	X, Y         int
	in, up       []int
	latte, malta *DinicEdge
	flow         *Dinic
}

func NewMaxFlowLowerBound(n int) *MaxFlowLowerBound {
	return &MaxFlowLowerBound{
		n:    n,
		X:    n,
		Y:    n + 1,
		in:   make([]int, n),
		flow: NewDinic(n + 2),
	}
}

// 从from到to的,容量为[lower,upper]的边
//  from != to
func (m *MaxFlowLowerBound) AddEdge(from, to, lower, upper int) {
	if from == to {
		return
	}
	m.flow.AddEdgeWithIndex(from, to, upper-lower, len(m.up))
	m.in[to] += lower
	m.in[from] -= lower
	m.up = append(m.up, upper)
}

// 求解最大流,返回最大流量,如果不存在可行流,则返回-1.
//  alias:有源汇上下界最大流
func (m *MaxFlowLowerBound) MaxFlow(s, t int) int {
	if m.CanFlow(s, t) {
		return m.flow.MaxFlow(s, t)
	}
	return -1
}

// 求解最小流,返回最小流量,如果不存在可行流,则返回-1.
//  alias:有源汇上下界最小流
func (m *MaxFlowLowerBound) MinFlow(s, t int) int {
	if m.CanFlow(s, t) {
		res := INF - m.latte.cap
		m.latte.cap = 0
		m.malta.cap = 0
		return res - m.flow.MaxFlow(t, s)
	}
	return -1
}

// 判断是否存在从s到t的可行流.
func (m *MaxFlowLowerBound) CanFlow(s, t int) bool {
	if s == t {
		return true
	}
	m.flow.AddEdge(t, s, INF)
	m.latte = &m.flow.graph[t][len(m.flow.graph[t])-1]
	m.malta = &m.flow.graph[s][len(m.flow.graph[s])-1]
	return m.CanCycleFlow()
}

// 判断整张图是否存在满足条件的循环流(所有点的出流量等于入流量).
//  alias:无源汇上下界可行流
func (m *MaxFlowLowerBound) CanCycleFlow() bool {
	m.build()
	res := m.flow.MaxFlow(m.X, m.Y)
	return res >= m.sum
}

// 输出每条边的流量.
func (m *MaxFlowLowerBound) GetFlows() []int {
	mp := make(map[int]int)
	for i := 0; i < len(m.flow.graph); i++ {
		for j := 0; j < len(m.flow.graph[i]); j++ {
			e := m.flow.graph[i][j]
			if !e.isRev && e.index != -1 {
				mp[e.index] = m.up[e.index] - e.cap
			}
		}
	}

	res := make([]int, len(mp))
	for i := 0; i < len(res); i++ {
		res[i] = mp[i]
	}
	return res
}

func (m *MaxFlowLowerBound) build() {
	for i := 0; i < m.n; i++ {
		if m.in[i] > 0 {
			m.flow.AddEdge(m.X, i, m.in[i])
			m.sum += m.in[i]
		} else if m.in[i] < 0 {
			m.flow.AddEdge(i, m.Y, -m.in[i])
		}
	}
}

type DinicEdge struct {
	to    int
	cap   int
	rev   int
	isRev bool
	index int
}

type Dinic struct {
	graph   [][]DinicEdge
	minCost []int
	iter    []int
}

func NewDinic(n int) *Dinic {
	return &Dinic{
		graph: make([][]DinicEdge, n),
	}
}

func (d *Dinic) AddEdge(from, to, cap int) {
	d.AddEdgeWithIndex(from, to, cap, -1)
}

func (d *Dinic) AddEdgeWithIndex(from, to, cap, index int) {
	d.graph[from] = append(d.graph[from], DinicEdge{to, cap, len(d.graph[to]), false, index})
	d.graph[to] = append(d.graph[to], DinicEdge{from, 0, len(d.graph[from]) - 1, true, index})
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for d.buildAugmentingPath(s, t) {
		d.iter = make([]int, len(d.graph))
		f := 0
		for {
			f = d.findMinDistAugmentPath(s, t, INF)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func (d *Dinic) findMinDistAugmentPath(idx, t, flow int) int {
	if idx == t {
		return flow
	}

	i := d.iter[idx]
	for i < len(d.graph[idx]) {
		e := d.graph[idx][i]
		if e.cap > 0 && d.minCost[idx] < d.minCost[e.to] {
			f := d.findMinDistAugmentPath(e.to, t, min(flow, e.cap))
			if f > 0 {
				d.graph[idx][i].cap -= f
				d.graph[e.to][e.rev].cap += f
				return f
			}
		}
		i++
		d.iter[idx]++
	}
	return 0
}

func (d *Dinic) buildAugmentingPath(s, t int) bool {
	d.minCost = make([]int, len(d.graph))
	for i := range d.minCost {
		d.minCost[i] = -1
	}
	d.minCost[s] = 0
	queue := []int{s}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range d.graph[v] {
			if e.cap > 0 && d.minCost[e.to] == -1 {
				d.minCost[e.to] = d.minCost[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.minCost[t] != -1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
