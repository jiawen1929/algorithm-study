## 查询两个树结点的 LCA

[golang 实现](https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta/graph_tree.go#L815)

1.  离线查询

    - **tarjan** O(n+q) dfs+并查集

2.  在线查询

    1. **暴力** 单次O(n)查询

       1. 上跳记录访问的路径，第一个交点即为 LCA
       2. 相交链表找交点(空间复杂度 O(1))
       3. 预处理 depth 和 parent 上跳，两个先跳到一样的高度再同时往上跳直到相等
       4. 从根节点自顶向下递归查询

    2. **倍增** O(nlogn)预处理 O(logn)查询，两个先跳到一样的高度再同时往上跳直到相等
       `倍增的最大优点是可以动态添加叶子节点(挂叶子)`
       这种方法将线性上跳变成了二进制上跳
       还有树分块的方法，上跳可以变成 O(sqrt(n))
    3. **重链剖分**，O(n)预处理 O(logn)查询，沿着最多 logn 段链上跳 `(最快,常数小,甚至比st表快)`
    4. **欧拉序 + RMQ**， O(nlogn)/O(n)预处理 O(1)查询，这两点之间的区间中，深度最小点就是 LCA。这可以用 RMQ 解决；利用加一减 1RMQ 优化的线性时间 LCA。
       https://www.cnblogs.com/pealicx/p/6859901.html
    5. **Schieber-Vishkin** O(n),O(1)
       https://codeforces.com/blog/entry/125371
       https://github.com/pranjalssh/CP_codes/blob/master/anta/!LCA.cpp

    6. **线性空间倍增(CompressedBinaryLift)**

       https://taodaling.github.io/blog/2020/03/18/binary-lifting/
       https://codeforces.com/blog/entry/74847
       https://codeforces.com/blog/entry/100826

3.  支持在线加点删点
    - **link-cut tree** O(nlogn)预处理 O(logn)查询 ，树链剖分 + splay 维护链的信息
    - **toptree** O(nlogn)预处理 O(logn)查询

---

https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/_note/lca_compare.cc#L10
//
// Least Common Ancestor (Heavy-light decomposition)
//
// 1) Heavy-Light decomposition
// 2) Doubling
// 3) Euler-tour + RMQ
//
// `Use 1 or 3.`
//

---

更快的 LA 算法：长链剖分 O(n)预处理，O(1)查询
https://37zigen.com/level-ancestor-problem/
https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/solution/on-o1mo-ban-by-hqztrue-ota3/

---

树上多个点的 LCA，就是 DFS 序最小的和 DFS 序最大的这两个点的 LCA
https://www.zhihu.com/question/46440863/answer/165741734

---

简单树论
https://www.cnblogs.com/alex-wei/p/Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree_Tree.html
