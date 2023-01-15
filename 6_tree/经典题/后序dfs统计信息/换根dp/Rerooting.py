"""
换根dp框架

op是相邻结点转移时,fromRes如何变化
merge是如何合并两个子节点的res
e是每个节点res的初始值

框架传入op和merge看似只求根节点0的值,实际上求出了每个点的dp值

https://atcoder.jp/contests/dp/submissions/22766939
https://nyaannyaan.github.io/library/tree/rerooting.hpp
"""


from typing import Callable, Generic, List, TypeVar

T = TypeVar("T")


class Rerooting(Generic[T]):

    __slots__ = ("adjList", "_n", "_decrement")

    def __init__(self, n: int, decrement: int = 0):
        self.adjList = [[] for _ in range(n)]
        self._n = n
        self._decrement = decrement

    def addEdge(self, u: int, v: int) -> None:
        u -= self._decrement
        v -= self._decrement
        self.adjList[u].append(v)
        self.adjList[v].append(u)

    def rerooting(
        self,
        e: Callable[[int], T],
        op: Callable[[T, T], T],
        composition: Callable[[T, int, int, int], T],
        root=0,
    ) -> List["T"]:
        root -= self._decrement
        assert 0 <= root < self._n
        parents = [-1] * self._n
        order = [root]
        stack = [root]
        while stack:
            cur = stack.pop()
            for next in self.adjList[cur]:
                if next == parents[cur]:
                    continue
                parents[next] = cur
                order.append(next)
                stack.append(next)

        dp1 = [e(i) for i in range(self._n)]
        dp2 = [e(i) for i in range(self._n)]
        for cur in order[::-1]:
            res = e(cur)
            for next in self.adjList[cur]:
                if parents[cur] == next:
                    continue
                dp2[next] = res
                res = op(res, composition(dp1[next], cur, next, 0))
            res = e(cur)
            for next in self.adjList[cur][::-1]:
                if parents[cur] == next:
                    continue
                dp2[next] = op(res, dp2[next])
                res = op(res, composition(dp1[next], cur, next, 0))
            dp1[cur] = res

        for newRoot in order[1:]:
            parent = parents[newRoot]
            dp2[newRoot] = composition(op(dp2[newRoot], dp2[parent]), parent, newRoot, 1)
            dp1[newRoot] = op(dp1[newRoot], dp2[newRoot])
        return dp1


# 310-求树上每个节点到其他节点的最远距离
# 310. 最小高度树
# 在所有可能的树中，具有最小高度的树（即，min(h)）被称为 最小高度树 。
class Solution:
    def findMinHeightTrees(self, n: int, edges: List[List[int]]) -> List[int]:
        def e(root: int) -> int:
            return 0

        def op(childRes1: int, childRes2: int) -> int:
            return max(childRes1, childRes2)

        def composition(fromRes: int, parent: int, cur: int, direction: int) -> int:
            if direction == 0:  # cur -> parent
                return fromRes + 1
            return fromRes + 1  # parent -> cur

        R = Rerooting(n)
        for u, v in edges:
            R.addEdge(u, v)
        maxDists = R.rerooting(e=e, op=op, composition=composition, root=0)
        min_ = min(maxDists)
        return [i for i in range(n) if maxDists[i] == min_]
