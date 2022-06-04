# G 公司有 n 个沿铁路运输线环形排列的仓库，每个仓库存储的货物数量不等。
# 如何用最少搬运量可以使 n 个仓库的库存数量相同。
# 搬运货物时，只能在相邻的仓库之间搬运。
# 数据保证一定有解。

# 1≤n≤100,
# 每个仓库的库存量不超过 100。


# 设平均值为 x。

# 使 n 个仓库的库存数量相同 相当于

# 比 x 大的流出 a[i]−x ,
# 比 x 小的流入 x−a[i],
# 其余点 流入流出平衡。
# 故令比 x 大为源点 (限制流出 a[i]−x)，比 x 小的为汇点 (限制流入 x−a[i])。
# 其余为中间点(流入流出平衡)。
# 由于费用流是在最大流条件下的，于是以上条件必然得到满足。
# https://www.acwing.com/activity/content/problem/content/2684/
