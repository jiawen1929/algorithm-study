## 重心分解+点分治 处理大规模的树上`路径`问题

1. 暴力遍历路径 `O(n^2)`
   枚举顶点，对每个点出发 dfs
2. 分治+递归 `O(nlogn*logn)`
   对每个点，考虑包含这个点的路径和不包含这个点的路径

   - **包含这个点的路径**，可以用 dfs 求出
   - **不包含这个点的路径**，删除这个点，然后对每个子树递归求解

   在单链时，直接枚举点会退化成 O(n^2)
   **重心分解**：如果选取一个 mid 值(树的重心)作为子树的根，那么每个递归问题中子树的大小都不超过整棵树的一半，所以每次递归问题规模可以下降至一半或更低，从而将时间复杂度降到 O(nlogn)

   > 应用:
   > https://zhuanlan.zhihu.com/p/359209926
   > 求无根树中长度为 k 的路径数目
   > 树上距离不超过 upper 的点对数
   > 求最长的 gcd 大于 1 的路径。
   > ...

## 在使用重心分解前，想想树形 dp/换根 dp 能不能做

## 点分治的核心其实就是树的重心，如果了解了树的重心的做法其实就知道怎么做了

「部分木->オイラーツアー, パス->HL 分解, 同心円状->重心分解」