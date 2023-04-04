# 下一个排列/上一个排列


from typing import Any, List, Tuple


def nextPermutation(nums: List[Any], inPlace=False) -> Tuple[bool, List[Any]]:
    """返回下一个字典序的排列，如果不存在，返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] >= nums[left]:  # 1. 找到最后一个递增位置
        left -= 1
    if left == 0:  # 全部递减
        return False, nums
    last = left - 1  # 最后一个递增位置

    while nums[right] <= nums[last]:  # 2. 找到最小的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递减数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


def prePermutation(nums: List[Any], inPlace=False) -> Tuple[bool, List[Any]]:
    """返回前一个字典序的排列,如果不存在,返回本身;时间复杂度O(n)"""
    if not inPlace:
        nums = nums[:]

    left = right = len(nums) - 1

    while left > 0 and nums[left - 1] <= nums[left]:  # 1. 找到最后一个递减位置
        left -= 1
    if left == 0:  # 全部递增
        return False, nums
    last = left - 1  # 最后一个递减位置

    while nums[right] >= nums[last]:  # 2. 找到最大的可交换的right，交换这两个数
        right -= 1
    nums[last], nums[right] = nums[right], nums[last]

    left, right = last + 1, len(nums) - 1  # 3. 翻转后面间这段递增数列
    while left < right:
        nums[left], nums[right] = nums[right], nums[left]
        left += 1
        right -= 1
    return True, nums


from sortedcontainers import SortedSet


# https://maspypy.github.io/library/seq/kth_next_permutation.hpp
# 有重复的情况：先把重复的元素贴上不同标签,最后换回来.
def kthNextPermutation(unique: List[Any], k: int, inPlace=False) -> Tuple[bool, List[Any], int]:
    """下k个字典序的排列

    Args:
        unique (List[Any]): `无重复元素`的数组
        k (int): 后续第k个(`本身算第0个`)
        inPlace (bool, optional): 是否原地修改. 默认为False

    Returns:
        Tuple[bool, List[Any], int]: `是否存在, 下k个排列, 需要移动的元素个数`
    """
    if not inPlace:
        unique = unique[:]
    rank, q = [], []
    ss = SortedSet()
    while k and unique:
        n = len(rank) + 1
        p = unique[-1]
        now = ss.bisect_left(p)
        k += now
        r = k % n
        k //= n
        rank.append(r)
        q.append(unique[-1])
        ss.add(unique[-1])
        unique.pop()

    if k:
        return False, unique, len(rank)

    move = len(rank)
    while len(rank):
        r = rank.pop()
        it = ss[r]
        unique.append(it)
        ss.remove(it)
    return True, unique, move


if __name__ == "__main__":
    isOk, nextP = nextPermutation([1, 2, 3])
    if isOk:
        print("nextP:", nextP)
    isOk, nextP = nextPermutation(list(map(int, "12345")))
    if isOk:
        print("nextP:", nextP)

    isOk, preP = prePermutation([3, 2, 1])
    if isOk:
        print("preP:", preP)
    isOk, preP = prePermutation(list(map(int, "12345")))
    if isOk:
        print("preP:", preP)

    isOk, kthP, move = kthNextPermutation(list(map(int, "12345")), 5)
    if isOk:
        print("kthP:", kthP, "move:", move)

    a1 = [1, 2, 3, 4, 5, 6, 7, 9, 23, 14, 56, 99, 876, 222, 444, 555]
    a2 = [1, 2, 3, 4, 5, 6, 7, 9, 23, 14, 56, 99, 876, 222, 444, 555]

    k = 1000000
    for _ in range(k):
        isOk, a1 = nextPermutation(a1, inPlace=True)
        if not isOk:
            break
    assert a1 == kthNextPermutation(a2, k)[1]
