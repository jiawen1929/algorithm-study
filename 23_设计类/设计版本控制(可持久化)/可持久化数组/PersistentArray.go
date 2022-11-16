// !可持久化数组
// https://www.luogu.com.cn/problem/P3919
// 如题，你需要维护这样的一个长度为N的数组，支持如下几种操作
// 1.在某个历史版本上修改某一个位置上的值
// 2.访问某个历史版本上的某一位置的值
// 此外，每进行一次操作（对于操作2，即为生成一个完全一样的版本，不作任何改动)，
// 就会生成一个新的版本。版本编号即为当前操作的编号(从1开始编号，版本0表示初始状态数组)

package main

import (
	"fmt"
)

func main() {
	// // time
	// nums1e5 := make([]E, 1e5)
	// for i := range nums1e5 {
	// 	nums1e5[i] = E(i)
	// }
	// root3 := Build(0, len(nums1e5)-1, nums1e5)
	// time1 := time.Now()
	// for i := 0; i < 1e5; i++ {
	// 	root4 := root3.Set(0, 10)
	// 	root4.Get(i)
	// }
	// fmt.Println(time.Since(time1))

	array1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	root1 := Build(0, len(array1)-1, array1)
	fmt.Println(root1)
	root2 := root1.Set(0, 10)
	fmt.Println(root1, root2)

}

type E = int

type Node struct {
	left, right           int
	size                  int
	value                 E
	leftChild, rightChild *Node
}

func Build(left, right int, nums []E) *Node {
	node := &Node{left: left, right: right}
	if left == right {
		node.value = nums[left]
		node.size = 1
		return node
	}
	mid := (left + right) >> 1
	node.leftChild = Build(left, mid, nums)
	node.rightChild = Build(mid+1, right, nums)
	node.pushUp()
	return node
}

func (o *Node) Get(index int) int {
	if o.left == o.right {
		return o.value
	}
	mid := (o.left + o.right) >> 1
	if index <= mid {
		return o.leftChild.Get(index)
	}
	return o.rightChild.Get(index)
}

func (o Node) Set(index int, value E) *Node {
	// !修改时拷贝一个新节点(这里用值作为接收者，已经隐式拷贝了一份root结点o)
	if o.left == o.right {
		o.value = value
		return &o
	}

	mid := (o.left + o.right) >> 1
	if index <= mid {
		o.leftChild = o.leftChild.Set(index, value)
	} else {
		o.rightChild = o.rightChild.Set(index, value)
	}

	return &o
}

func (o *Node) pushUp() {
	o.size = o.leftChild.size + o.rightChild.size
}

func (o *Node) String() string {
	res := o.dfs()
	return fmt.Sprintf("PersistentArray %v", res)
}

func (o *Node) dfs() []E {
	res := make([]E, 0, o.size)
	if o.left == o.right {
		res = append(res, o.value)
		return res
	}
	res = append(res, o.leftChild.dfs()...)
	res = append(res, o.rightChild.dfs()...)
	return res
}
