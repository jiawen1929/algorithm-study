# 629. K个逆序对数组
# https://leetcode.cn/problems/k-inverse-pairs-array/
# 插入dp

from itertools import accumulate


MOD = int(1e9 + 7)


class Solution:
    def kInversePairs(self, n: int, k: int) -> int:
        """
        包含从 1 到 n 的数字，且恰好拥有 k 个逆序对的不同的数组的个数

        n,k<=1000
        https://leetcode-cn.com/problems/k-inverse-pairs-array/solution/dong-tai-gui-hua-qian-zhui-he-you-hua-by-28eb/

        dp[i][j]=dp[i-1][j]+dp[i-1][j-1]+dp[i-1][j-2]+...+dp[i-1][j-i]
        滚动数组dp,后面这一串用 `前缀和` O(1) 求出 preSum[j+1]-preSum[j-i-1]
        """
        dp = [0] * (k + 1)
        dp[0] = 1
        for i in range(1, n):
            ndp = dp[:]
            dpSum = [0] + list(accumulate(dp))
            for j in range(1, k + 1):
                ndp[j] = (dpSum[j + 1] - dpSum[max(0, j - i)]) % MOD  # !一共(i+1)项
            dp = ndp
        return dp[k]


print(Solution().kInversePairs(n=3, k=1))  # 2
print(Solution().kInversePairs(n=1000, k=1000))  # 2
print(Solution().kInversePairs(n=5, k=3))  # 15
# 663677020
