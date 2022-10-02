"""
字符串哈希
注意字符串较短时,使用切片更快,时间复杂度约为`O(S/250)`
"""

from typing import Sequence


def useStringHasher(s: Sequence[str], mod=10**11 + 7, base=1313131):
    n = len(s)
    prePow = [1] * (n + 2)
    preHash = [0] * (n + 1)
    for i in range(n):
        prePow[i + 1] = (prePow[i] * base) % mod
        preHash[i] = (preHash[i - 1] * base + ord(s[i])) % mod

    def sliceHash(left: int, right: int):
        """切片 `s[left:right]` 的哈希值"""
        right -= 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


def useArrayHasher(nums: Sequence[int], mod=10**11 + 7, base=1313131):
    n = len(nums)
    prePow = [1] * (n + 2)
    preHash = [0] * (n + 1)
    for i in range(n):
        prePow[i + 1] = (prePow[i] * base) % mod
        preHash[i] = (preHash[i - 1] * base + nums[i]) % mod

    def sliceHash(left: int, right: int):
        """切片 `nums[left:right]` 的哈希值"""
        right -= 1
        return (preHash[right] - preHash[left - 1] * prePow[right - left + 1]) % mod

    return sliceHash


if __name__ == "__main__":
    stringHasher = useStringHasher("abc")
    print(stringHasher(1, 2))
    print(stringHasher(1, 3))
    print(stringHasher(0, 0))
    print(stringHasher(1, 1))
    print(stringHasher(0, 2))

    arrayHasher = useArrayHasher([1, 2, 3])
    print(arrayHasher(1, 2))
    print(arrayHasher(1, 3))
    print(arrayHasher(0, 0))

    from functools import lru_cache

    class Solution:
        def deleteString(self, s: str) -> int:
            @lru_cache(None)
            def dfs(index: int) -> int:
                if index == n:
                    return 0

                remain = n - index
                res = 1
                for i in range(1, remain // 2 + 1):
                    if s[index : i + index] == s[i + index : i + i + index]:  # 字符串短时切片更快
                        # if hasher(index, i + index) == hasher(i + index, i + i + index):
                        res = max(res, dfs(i + index) + 1)
                return res

            n = len(s)
            # hasher = useStringHasher(s)
            res = dfs(0)
            dfs.cache_clear()
            return res
