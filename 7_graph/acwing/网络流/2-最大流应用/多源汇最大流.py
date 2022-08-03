# 多源汇添加虚拟点后就变成了单源汇
"""
添加一个虚拟源点SS和虚拟汇点TT，
SS 和 所有原图源点间连接容量无穷大的边
原图汇点和TT 连接容量无穷大的边
求SS到TT的最大流
"""
# 2≤n≤10000,
# 1≤m≤105,

from collections import defaultdict, deque
from typing import Set


INF = int(1e18)


class MaxFlow:
    def __init__(self, start: int, end: int) -> None:
        self.graph = defaultdict(lambda: defaultdict(int))  # 原图
        self._start = start
        self._end = end

    def calMaxFlow(self) -> int:
        self._updateRedisualGraph()
        start, end = self._start, self._end
        flow = 0

        while self._bfs():
            delta = INF
            while delta:
                delta = self._dfs(start, end, INF)
                flow += delta
        return flow

    def addEdge(self, v1: int, v2: int, w: int, *, cover=False) -> None:
        """添加边 v1->v2, 容量为w

        Args:
            v1: 边的起点
            v2: 边的终点
            w: 边的容量
            cover: 是否覆盖原有边
        """
        if cover:
            self.graph[v1][v2] = w
        else:
            self.graph[v1][v2] += w

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """边的流量=容量-残量"""
        assert v1 in self.graph and v2 in self.graph[v1]
        return self.graph[v1][v2] - self._reGraph[v1][v2]

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        """边的残量(剩余的容量)"""
        assert v1 in self.graph and v2 in self.graph[v1]
        return self._reGraph[v1][v2]

    def getPath(self) -> Set[int]:
        """最大流经过了哪些点"""
        visited = set()
        stack = [self._start]
        reGraph = self._reGraph
        while stack:
            cur = stack.pop()
            visited.add(cur)
            for next, remain in reGraph[cur].items():
                if next not in visited and remain > 0:
                    visited.add(next)
                    stack.append(next)
        return visited

    def _updateRedisualGraph(self) -> None:
        """残量图 存储每条边的剩余流量"""
        self._reGraph = defaultdict(lambda: defaultdict(int))
        for cur in self.graph:
            for next, cap in self.graph[cur].items():
                self._reGraph[cur][next] = cap
                self._reGraph[next].setdefault(cur, 0)  # 注意自环边

    def _bfs(self) -> bool:
        self._depth = depth = defaultdict(lambda: -1, {self._start: 0})
        reGraph, start, end = self._reGraph, self._start, self._end
        queue = deque([start])
        self._iters = {cur: iter(reGraph[cur].keys()) for cur in reGraph.keys()}
        while queue:
            cur = queue.popleft()
            nextDist = depth[cur] + 1
            for next, remain in reGraph[cur].items():
                if depth[next] == -1 and remain > 0:
                    depth[next] = nextDist
                    queue.append(next)

        return depth[end] != -1

    def _dfs(self, cur: int, end: int, flow: int) -> int:
        if cur == end:
            return flow
        reGraph, depth, iters = self._reGraph, self._depth, self._iters
        for next in iters[cur]:
            remain = reGraph[cur][next]
            if remain and depth[cur] < depth[next]:
                nextFlow = self._dfs(next, end, min(flow, remain))
                if nextFlow:
                    reGraph[cur][next] -= nextFlow
                    reGraph[next][cur] += nextFlow
                    return nextFlow
        return 0


n, m, startCount, endCount = map(int, input().split())
starts = list(map(int, input().split()))
ends = list(map(int, input().split()))
START, END, OFFSET = -1, -2, int(1e6)
maxFlow = MaxFlow(START, END)
for _ in range(m):
    u, v, w = map(int, input().split())
    maxFlow.addEdge(u, v, w)

for num in starts:
    maxFlow.addEdge(START, num, int(1e9))
for num in ends:
    maxFlow.addEdge(num, END, int(1e9))


print(maxFlow.calMaxFlow())
