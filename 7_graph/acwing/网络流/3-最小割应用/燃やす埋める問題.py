# https://zenn.dev/kiwamachan/articles/37a2c646f82c7d
# https://shindannin.hatenadiary.com/entry/2017/11/15/043009
# 最小カットを使って「燃やす埋める問題」を解く
# `最小割考虑的是不连接时的花费`

# a1,...,an の　N個の廃棄物がある。
# 全て燃やすか埋めるかして処分する。それぞれにコストがかかる。
# 燃やすコスト：[m1,m2,...,mn]
# 埋めるコスト：[u1,u2,...,un]
# !さらに a1を燃やしてa2を埋めると罰金100円　2つの処分方法の組合せによる追加条件がいくつかある
# かかるコストを最小化せよ。
# !即求最小割(最少的罚款)

# 源点为燃やす
# 汇点为埋める
# !a1を燃やしてa2を埋めると罰金100円:a1=>a2有一条罚款100的流
# !求图中的最小割 (即每个点都有唯一的归属)
