# 2936. 包含相等值数字块的数量
# https://leetcode.cn/problems/number-of-equal-numbers-blocks/

# 给定一个整数数组 nums，其 下标从 0 开始。对于 nums，有以下性质：

# 所有相同值的元素都是相邻的。换句话说，如果存在两个下标 i < j，使得 nums[i] == nums[j]，
# 那么对于所有下标 k，满足 i < k < j，都有 nums[k] == nums[i]。
# 由于 nums 是一个非常大的数组，这里提供了一个 BigArray 类的实例，该实例具有以下函数：

# int at(long long index): 返回 nums[i] 的值。
# void size(): 返回 nums.length。
# 让我们把数组分成 最大 的块，使得每个块包含 相等的值。返回这些块的数量。

# 输入：nums = [1,1,1,3,9,9,9,2,10,10]
# 输出：5
# 解释：这里有 5 个块：
# 块号 1: [1,1,1,3,9,9,9,2,10,10]
# 块号 2: [1,1,1,3,9,9,9,2,10,10]
# 块号 3: [1,1,1,3,9,9,9,2,10,10]
# 块号 4: [1,1,1,3,9,9,9,2,10,10]
# 块号 5: [1,1,1,3,9,9,9,2,10,10]
# 因此答案是 5。

# 1 <= nums.length <= 1e15
# 1 <= nums[i] <= 1e9
# 在生成的输入中所有相同值的元素是相邻的。
# nums 的所有元素之和最多为 1e15。


class BigArray:
    def at(self, index: int) -> int:
        ...

    def size(self) -> int:
        ...


class Solution:
    def countBlocks2(self, nums: "BigArray") -> int:
        def dfs(start: int, end: int) -> int:
            if nums.at(start) == nums.at(end - 1):  # 整个区间相等
                return 1
            mid = (start + end) // 2
            leftRes = dfs(start, mid)
            rightRes = dfs(mid, end)
            return leftRes + rightRes - (nums.at(mid - 1) == nums.at(mid))

        n = nums.size()
        return dfs(0, n)

    def countBlocks(self, nums: "BigArray") -> int:
        def maxRight(left: int) -> int:
            right = n
            cur = nums.at(left)
            while left <= right:
                mid = (left + right) // 2
                if mid < n and nums.at(mid) == cur:
                    left = mid + 1
                else:
                    right = mid - 1
            return left

        n = nums.size()
        res = 0
        pos = 0
        while pos < n:
            pos = maxRight(pos)
            res += 1
        return res
