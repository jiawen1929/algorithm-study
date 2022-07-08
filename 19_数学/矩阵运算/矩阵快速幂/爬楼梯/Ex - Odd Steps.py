# 爬s阶(s<=1e18)楼梯 每次只能爬`奇数级`
# 有些楼梯断了不能爬
# 爬到最后一个位置有多少种方法

# !dp[i] = dp[i-1] + dp[i-3] + dp[i-5] + ...
# !令 ep[i] = dp[i] + dp[i-2] + dp[i-4] + ...
# !那么 dp[i] = ep[i-1] = dp[i-1] + ep[i-3]
# dp[i] ep[i-1] ep[i-2] 的关系可由矩阵快速幂 logS 求出
# 转移矩阵
# 1 0 1
# 1 0 1
# 0 1 0
# 初始向量 [1 0 0]

# !坏的楼梯处
# !矩阵快速幂中间要打断 每次都直接把dp[i]手动赋值为0 再继续算
# 总时间复杂度 O(nlogS)
import gc
import sys
import os
import numpy as np

gc.disable()


input = sys.stdin.readline
MOD = np.uint64(998244353)


NPArray = np.ndarray


def matqpow2(base: NPArray, exp: int, mod: np.uint64) -> NPArray:
    """矩阵快速幂np版"""

    base = base.copy()
    res = np.eye(*base.shape, dtype=np.uint64)

    while exp:
        if exp & 1:
            res = (res @ base) % mod
        exp >>= 1
        base = (base @ base) % mod
    return res


def main() -> None:
    _, s = map(int, input().split())
    bad = list(map(int, input().split()))
    res = np.array([[1], [0], [0]], np.uint64)  # 3 x 1 答案矩阵
    trans = np.array([[1, 0, 1], [1, 0, 1], [0, 1, 0]], np.uint64)
    pre = 0
    for cur in bad:
        res = (matqpow2(trans, cur - pre, MOD) @ res) % MOD
        res[0][0] = 0  # 坏的楼梯不能走
        pre = cur

    res = (matqpow2(trans, s - bad[-1], MOD) @ res) % MOD
    print(int(res[0][0]))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
