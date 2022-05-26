# 牛牛有一个长度为的数组，牛妹给出个询问，询问有种类型：
# :询问区间内有多少子序列的乘积为奇数
# :询问区间内有多少子序列的乘积为偶数
# 某个序列的子序列是从最初序列通过去除某些元素但不破坏余下元素的相对位置（在前或在后）而形成的新序列。


# 乘积为奇数的子序列一定是全都是奇数，所以是 2^(奇数数量) - 1
# 求 [l,r] 内有`多少个奇数可以将所有数字模2，然后求前缀和`
# 如果就是求奇数，直接输出
# 求偶数就再求一个全部子序列的数量：2^(r - l + 1) - 1，再减去乘积为奇数的子序列数量
# 至于求2的多少次幂可以O(n)预处理，也可以快速幂

# 区间内有多少子序列的乘积为奇数偶数
