from collections import defaultdict, deque
from typing import List, Set


class MaxFlowMap:
    INF = int(1e18)

    def __init__(self, start: int, end: int) -> None:
        self._graph = defaultdict(lambda: defaultdict(int))
        self._start = start
        self._end = end

    def calMaxFlow(self) -> int:
        self._updateRedisualGraph()
        start, end, INF = self._start, self._end, self.INF
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
            self._graph[v1][v2] = w
        else:
            self._graph[v1][v2] += w

    def getFlowOfEdge(self, v1: int, v2: int) -> int:
        """边的流量=容量-残量"""
        assert v1 in self._graph and v2 in self._graph[v1]
        return self._graph[v1][v2] - self._reGraph[v1][v2]

    def getRemainOfEdge(self, v1: int, v2: int) -> int:
        """边的残量(剩余的容量)"""
        assert v1 in self._graph and v2 in self._graph[v1]
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
        for cur in self._graph:
            for next, cap in self._graph[cur].items():
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


# 相邻两个1组成一条边，每条边都要去掉一个端点，其实是找最小点覆盖，即求二分图的最大匹配，跑匈牙利算法
# !起点 => 奇数黑格 => 偶数黑格 => 终点
DIR2 = [(0, 1), (1, 0)]


class Solution:
    def minimumOperations(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        START, END = -1, -2
        maxFlow = MaxFlowMap(start=START, end=END)

        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 1:
                    cur = r * COL + c
                    for dr, dc in DIR2:
                        nr, nc = r + dr, c + dc
                        if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 1:
                            next = nr * COL + nc
                            v1, v2 = (next, cur) if (r + c) & 1 else (cur, next)
                            maxFlow.addEdge(v1, v2, 1, cover=True)
                            maxFlow.addEdge(START, v1, 1, cover=True)
                            maxFlow.addEdge(v2, END, 1, cover=True)

        return maxFlow.calMaxFlow()


print(Solution().minimumOperations(grid=[[1, 1, 0], [0, 1, 1], [1, 1, 1]]))
