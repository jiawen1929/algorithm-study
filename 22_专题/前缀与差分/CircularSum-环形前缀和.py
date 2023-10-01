# 环形数组前缀和/环形前缀和


from itertools import accumulate
from random import randint
from typing import Callable, List

INF = int(1e18)


def circularPresum(nums: List[int]) -> Callable[[int, int], int]:
    """环形数组前缀和."""
    n = len(nums)
    preSum = [0] + list(accumulate(nums))

    def _cal(r: int) -> int:
        return preSum[n] * (r // n) + preSum[r % n]

    def query(start: int, end: int) -> int:
        """[start,end)的和.
        0 <= start < end <= n.
        """
        if start >= end:
            return 0
        return _cal(end) - _cal(start)

    return query


if __name__ == "__main__":
    nums = list(range(1000))
    cs = circularPresum(nums)
    for _ in range(100):
        # check with bf
        left, right = randint(0, 100), randint(0, 100)
        if left > right:
            left, right = right, left
        sum1 = cs(left, right)
        sum2 = sum(nums[i] for i in range(left, right))
        assert sum1 == sum2, (left, right, sum1, sum2)

    # F - More Holidays
    # https://atcoder.jp/contests/abc300/tasks/abc300_f
    # 给定一个01字符串t，它由一个长度为n的串s重复m次拼接得到。
    # 要求将恰好 k个0变成1，问连续1的最大长度。
    # !x->0 o->1
    # !枚举修改的左端点，然后二分找到恰好包含k个0的最右端点
    def moreHolidays(s: str, m: int, k: int) -> int:
        nums = [0 if c == "x" else 1 for c in s]
        query = circularPresum(nums)
        res = 0
        # !枚举变换的0的起点 二分求出右边界
        n = len(nums)
        for i in range(n):
            left, right = i + 1, m * n
            while left <= right:
                mid = (left + right) // 2
                ones = (mid - i) - query(i, mid)
                if ones <= k:
                    left = mid + 1
                else:
                    right = mid - 1
            res = max(res, right - i)
        return res

    n, m, k = map(int, input().split())
    s = input()
    print(moreHolidays(s, m, k))

    # 100076. 无限数组的最短子数组
    # https://leetcode.cn/problems/minimum-size-subarray-in-infinite-array/
    # 求循环数组中和为 target 的最短子数组的长度.不存在则返回 -1.
    # 1 <= nums.length <= 1e5
    # 1 <= nums[i] <= 1e5
    # 1 <= target <= 1e9
    class Solution:
        def minSizeSubarray(self, nums: List[int], target: int) -> int:
            Q = circularPresum(nums)
            res = INF
            for start in range(len(nums)):
                left, right = 0, int(1e9 + 10)
                while left <= right:
                    mid = (left + right) // 2
                    curSum = Q(start, start + mid)
                    if curSum == target:
                        res = min(res, mid)
                        break
                    elif curSum < target:
                        left = mid + 1
                    else:
                        right = mid - 1

            return res if res != INF else -1
