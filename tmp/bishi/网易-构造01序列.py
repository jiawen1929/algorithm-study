# TODO
# !小红希望你构造一个长度为n的01串，其中恰好有k个1,且恰好有t对相邻字符串对1。
# n<=1e5
# !m个连续1有m*(m-1)//2 对相邻的1
# 相当于k个1之间插入n-k个0
# !先放1，再插0，多余的0全部放最后面

# 额外规定加一个0放在最后，那么就是k-t个10，
# 加上n+1-k-k+t个0随意排列是C(n+1-k, k-t)，
# 最后剩下的t个1 扔到k-t个10的前面位置，
# 属于小球相同盒子不同且允许为空C(k-1,t-1)，结果乘起来就好了
