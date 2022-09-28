from typing import List
from collections import deque

# !给你一个二维整数数组 relations ，其中 relations[j] = [prevCoursej, nextCoursej] ，
# !表示课程 prevCoursej 必须在课程 nextCoursej 之前 完成（先修课的关系）
# !其中 time[i] 表示完成第 (i+1) 门课程需要花费的 月份 数。
# 请你根据以下规则算出完成所有课程所需要的 最少 月份数：
# 如果一门课的所有先修课都已经完成，你可以在 任意 时间开始这门课程。
# 你可以 同时 上 任意门课程 。


# !拓扑排序最长路


def minimumTime(n: int, relations: List[List[int]], time: List[int]) -> int:
    adjList = [[] for _ in range(n + 1)]
    deg = [0] * (n + 1)
    for u, v in relations:
        adjList[u].append(v)
        deg[v] += 1

    # !记录虚拟原点0到达每个点处所需要的距离
    dp = [0] * (n + 1)
    queue = deque()
    for i in range(1, n + 1):
        if deg[i] == 0:
            queue.append(i)
            dp[i] = time[i - 1]

    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            weight = time[next - 1]
            dp[next] = max(dp[next], dp[cur] + weight)
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)

    return max(dp)


print(minimumTime(n=3, relations=[[1, 3], [2, 3]], time=[3, 2, 5]))
# 输出：8
# 解释：上图展示了输入数据所表示的先修关系图，以及完成每门课程需要花费的时间。
# 你可以在月份 0 同时开始课程 1 和 2 。
# 课程 1 花费 3 个月，课程 2 花费 2 个月。
# 所以，最早开始课程 3 的时间是月份 3 ，完成所有课程所需时间为 3 + 5 = 8 个月。
