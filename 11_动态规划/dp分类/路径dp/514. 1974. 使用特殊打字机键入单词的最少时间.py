# 每一秒钟，你可以执行以下操作之一：

# 将指针 顺时针 或者 逆时针 移动一个字符。
# 键入指针 当前 指向的字符。
# 给你一个字符串 word ，请你返回键入 word 所表示单词的 最少 秒数 。


from itertools import pairwise


class Solution:
    def minTimeToType(self, word: str) -> int:
        word = 'a' + word
        nums = list(map(ord, word))
        res = len(word)
        for pre, cur in pairwise(nums):
            res += min(abs(pre - cur), 26 - abs(pre - cur))
        return res - 1


print(Solution().minTimeToType("bza"))
# 输出：7
# 解释：
# 单词按如下操作键入：
# - 花 1 秒将指针顺时针移到 'b' 。
# - 花 1 秒键入字符 'b' 。
# - 花 2 秒将指针逆时针移到 'z' 。
# - 花 1 秒键入字符 'z' 。
# - 花 1 秒将指针顺时针移到 'a' 。
# - 花 1 秒键入字符 'a' 。

