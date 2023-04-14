package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main2() {
	assert := func(a bool) {
		if !a {
			panic("assert fail")
		}
	}
	queue := NewPhysicistQueue()
	assert(queue.Empty() == true)
	queue1 := queue.Push(1)
	assert(queue1.Empty() == false)
	assert(queue1.Front() == 1)
	queue2 := queue1.Push(2)
	assert(queue2.Empty() == false)
	assert(queue2.Front() == 1)
	queue3 := queue2.Shift()
	assert(queue3.Empty() == false)
	assert(queue3.Front() == 2)
	queue4 := queue3.Shift()
	assert(queue4.Empty() == true)

	time1 := time.Now()
	for i := 0; i < 1e5; i++ {
		queue = queue.Push(i)
		queue = queue.Shift()
		queue.Empty()
	}

	fmt.Println(time.Since(time1)) // 343ms
}

func main() {
	// https://judge.yosupo.jp/problem/persistent_queue
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	gits := make([]*PhysicistQueue, q+1)
	gits[0] = NewPhysicistQueue()

	for i := 0; i < q; i++ {
		var op, time, x int
		fmt.Fscan(in, &op, &time)
		time++
		if op == 0 {
			fmt.Fscan(in, &x)
			gits[i+1] = gits[time].Push(x)
		} else {
			res := gits[time].Front()
			fmt.Fprintln(out, res)
			gits[i+1] = gits[time].Shift()
		}
	}
}

type S = int

type PhysicistQueue struct {
	working      *_PersistentStack
	front        *_Suspension
	back         *_PersistentStack
	fSize, bSize int
}

func NewPhysicistQueue() *PhysicistQueue {
	return &PhysicistQueue{
		working: _NewPersistentStack(),
		front:   _NewSuspensionWith(func() interface{} { return _NewPersistentStack() }),
		back:    _NewPersistentStack(),
	}
}

func (q *PhysicistQueue) Empty() bool {
	return q.fSize == 0
}

func (q *PhysicistQueue) Front() S {
	if q.Empty() {
		panic("queue is empty")
	}
	return q.working.Top()
}

func (q *PhysicistQueue) Push(x S) *PhysicistQueue {
	res := &PhysicistQueue{
		working: q.working,
		front:   q.front,
		fSize:   q.fSize,
		back:    q.back.Push(x),
		bSize:   q.bSize + 1,
	}
	return res._check()
}

func (q *PhysicistQueue) Shift() *PhysicistQueue {
	if q.Empty() {
		panic("queue is empty")
	}
	res := &PhysicistQueue{
		working: q.working.Pop(),
		front:   _NewSuspensionWith(func() interface{} { return q.front.Resolve().(*_PersistentStack).Pop() }),
		fSize:   q.fSize - 1,
		back:    q.back,
		bSize:   q.bSize,
	}
	return res._check()
}

func (q *PhysicistQueue) _check() *PhysicistQueue {
	return q._checkR()._checkW()
}
func (q *PhysicistQueue) _checkR() *PhysicistQueue {
	if q.fSize >= q.bSize {
		return q
	}
	tmp := q.front.Resolve().(*_PersistentStack)

	f := func() interface{} {
		r := q.back.Reverse()
		l := tmp.Reverse()
		for !l.Empty() {
			r = r.Push(l.Top())
			l = l.Pop()
		}
		return r
	}

	return &PhysicistQueue{
		working: tmp,
		front:   _NewSuspensionWith(f),
		fSize:   q.fSize + q.bSize,
		back:    _NewPersistentStack(),
	}
}

func (q *PhysicistQueue) _checkW() *PhysicistQueue {
	if !q.working.Empty() {
		return q
	}
	return &PhysicistQueue{
		working: q.front.Resolve().(*_PersistentStack),
		front:   q.front,
		fSize:   q.fSize,
		back:    q.back,
		bSize:   q.bSize,
	}
}

type _Cell struct {
	resolved S
	next     *_Stream
}

type _Stream struct {
	*_Suspension
}

// 惰性求值的流.
func _NewStream() *_Stream {
	return &_Stream{_Suspension: _NewSuspension()}
}

func _Concat(x, y *_Stream) *_Stream {
	return &_Stream{_Suspension: _NewSuspensionWith(func() interface{} {
		if x.Empty() {
			return y.Resolve()
		}
		return &_Cell{resolved: x.Top(), next: _Concat(x.Pop(), y)}
	})}
}

func (s *_Stream) Empty() bool {
	return s.Resolve() == nil
}

func (s *_Stream) Top() S {
	return s.Resolve().(*_Cell).resolved
}

func (s *_Stream) Pop() *_Stream {
	return s.Resolve().(*_Cell).next
}

func (s *_Stream) Push(x S) *_Stream {
	return &_Stream{_Suspension: _NewSuspensionWith(&_Cell{resolved: x, next: s})}
}

func (s *_Stream) Reverse() *_Stream {
	return &_Stream{_Suspension: _NewSuspensionWith(func() interface{} {
		x := s
		res := _NewStream()
		for !x.Empty() {
			res = res.Push(x.Top())
			x = x.Pop()
		}
		return res.Resolve()
	})}
}

func (s *_Stream) String() string {
	x := s
	res := []S{}
	for !x.Empty() {
		res = append(res, x.Top())
		x = x.Pop()
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return fmt.Sprintf("Stream%v", res)
}

// 惰性求值.
type _Suspension struct {
	x        interface{}
	resolved interface{}
}

func _NewSuspension() *_Suspension {
	return &_Suspension{}
}

func _NewSuspensionWith(x interface{}) *_Suspension {
	return &_Suspension{x: x}
}

func (s *_Suspension) Resolve() interface{} {
	if s.resolved == nil {
		if f, ok := s.x.(func() interface{}); ok {
			s.resolved = f()
		} else {
			s.resolved = s.x
		}
	}
	return s.resolved
}

type _PersistentStack struct {
	root *_StackNode
}

type _StackNode struct {
	value S
	pre   *_StackNode
}

// 创建一个新的可持久化栈.
func _NewPersistentStack() *_PersistentStack {
	return &_PersistentStack{}
}

func (s *_PersistentStack) Push(value S) *_PersistentStack {
	return &_PersistentStack{root: &_StackNode{value: value, pre: s.root}}
}

func (s *_PersistentStack) Pop() *_PersistentStack {
	if s.root == nil {
		panic("stack is empty")
	}
	return &_PersistentStack{root: s.root.pre}
}

func (s *_PersistentStack) Top() S {
	if s.root == nil {
		panic("stack is empty")
	}
	return s.root.value
}

func (s *_PersistentStack) Empty() bool {
	return s.root == nil
}

func (s *_PersistentStack) Reverse() *_PersistentStack {
	res := _NewPersistentStack()
	x := s
	for !x.Empty() {
		res = res.Push(x.Top())
		x = x.Pop()
	}
	return res
}

func (s *_PersistentStack) String() string {
	sb := []string{}
	x := s
	for !x.Empty() {
		sb = append(sb, fmt.Sprintf("%v", x.Top()))
		x = x.Pop()
	}
	for i, j := 0, len(sb)-1; i < j; i, j = i+1, j-1 {
		sb[i], sb[j] = sb[j], sb[i]
	}
	return fmt.Sprintf("Stack%v", sb)
}
