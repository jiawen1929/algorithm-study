# 考试中有 n 种类型的题目。给你一个整数 target 和一个下标从 0 开始的二维整数数组 types ，
# 其中 types[i] = [counti, marksi] 表示第 i 种类型的题目有 counti 道，每道题目对应 marksi 分。
# !返回你在考试中恰好得到 target 分的方法数。由于答案可能很大，结果需要对 1e9 +7 取余。
# !注意，同类型题目无法区分。
# target<=1000
# n<=50
# counti<=50

# !O(n*target) 按模分组前缀和优化dp
# dp[i][j]表示前i种题目恰好得到j分的方法数
# !ndp[j] = sum(dp[j-k*mark] for k in range(count+1) if j-k*mark>=0
# 这是一个按模分组的前缀和


from typing import List

MOD = int(1e9 + 7)


class Solution:
    def waysToReachTarget(self, target: int, types: List[List[int]]) -> int:
        dp = [0] * (target + 1)
        dp[0] = 1
        for count, mark in types:
            ndp = dp[:]
            preSum = dp[:]  # 按模分组前缀和
            for v in range(mark, target + 1):
                preSum[v] = (preSum[v] + preSum[v - mark]) % MOD
            for v in range(target + 1):
                ndp[v] = preSum[v]
                # !最多选count个,所以要减去得到v-(count+1)*mark的前缀和(超过count个的方法数都是不合法的)
                if (v - (count + 1) * mark) >= 0:
                    ndp[v] = (ndp[v] - preSum[v - (count + 1) * mark]) % MOD
            dp = ndp
        return dp[target]
