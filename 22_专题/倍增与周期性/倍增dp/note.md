https://taodaling.github.io/blog/2020/03/18/binary-lifting/

## 定义

倍增技术的强大是基于一个很简单的倍增结构。
给定一个有向图 G=(V,E)，每个结点只有唯一的出边。

- `next(u)`: u 沿着唯一的出边移动一步后抵达的位置
- `link(u,i)`：从 u 出发，通过少于 i 次 next 的转移能抵达的所有结点的集合。
- `jump(u,i)`：从 u 出发沿着出边移动 2^i 步所在的位置，
  jump(u,0)=next(u), jump(u,i)=jump(jump(u,i-1),i-1)
  在图上进行倍增实际上就是维护 jump.
  **将 jump(u,i)视为一个结点，它覆盖了所有 link(u,2^i)上的结点，称 i 为这个结点的高度。**

## 问题：

1. 我们需要处理若干个请求，每个请求要求修改路径 link(u,l)上的所有结点。在所有请求完成后，要求输出所有结点的权值。
   实际上可以发现 u,v 对应的区间可以截断为 O(log2n)个倍增结构上的结点，我们只需要在这些结点上**打上标记**就可以了。并且考虑到标记只需要从高度较大的结点下推到高度较小的结点，因此在最后阶段我们可以**从高到低处理结点**。
   倍增结构可以处理存在环的情况
2. 有 n 个人，以及一颗大小为 m 的树。第 i 个人可以居住在 ui 和 vi 之间的路径上的任意一个顶点中，且一个顶点最多居住一个人。现在希望让尽可能多人居住在树上，问最多有多少人可以居住在树上。
   每个结点只有一个父亲，因此我们把 next(u)设置为 u 的父亲。之后我们建立网络流，将第 i 个人向路径 ui 到 vi 上的所有结点连一条边，借助倍增结构，我们可以**只需要向 O(log2m)个结点连边**。这样我们得到了一个包含 O(mlog2m+n)个结点和 O((n+m)log2m)条边的网络。之后求最大流即可。

## 倍增拆点技巧：

- 拆序列上的区间： [RangeUnionFindTreeOffline](../%E5%80%8D%E5%A2%9E%E4%BC%98%E5%8C%96%E5%BB%BA%E5%9B%BE/RangeUnionFindTreeOffline.go)
  倍增拆区间和线段树拆区间的区别：
  倍增拆分出的区间形态总是相同的，但是线段树拆分出的区间形态不同。
- 树(倍增结构)上的一段路径拆成 logn 个点：

---

题目特点：进入一个状态，下一个状态一定确定
