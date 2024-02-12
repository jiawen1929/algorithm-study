from collections import Counter
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# AtCoder 諸島は
# N 個の島からなり、これらの島々は
# N 本の橋によって結ばれています。 島には
# 1 から
# N までの番号が付けられていて、
# i (1≤i≤N−1) 本目の橋は島
# i と島
# i+1 を、
# N 本目の橋は島
# N と島
# 1 を双方向に結んでいます。 橋を渡る以外に島の間を行き来する方法は存在しません。

# AtCoder 諸島では、島
# X
# 1
# ​
#   から始めて島
# X
# 2
# ​
#  ,X
# 3
# ​
#  ,…,X
# M
# ​
#   を順に訪れるツアーが定期的に催行されています。 移動の過程で訪れる島とは別の島を経由することもあり、ツアー中に橋を通る回数の合計がツアーの長さと定義されます。

# 厳密には、ツアーとは以下の条件を全て満たす
# l+1 個の島の列
# a
# 0
# ​
#  ,a
# 1
# ​
#  ,…,a
# l
# ​
#   のことであり、その長さ は
# l として定義されます。

# 全ての
# j (0≤j≤l−1) について、島
# a
# j
# ​
#   と島
# a
# j+1
# ​
#   は橋で直接結ばれている
# ある
# 0＝y
# 1
# ​
#  <y
# 2
# ​
#  <⋯<y
# M
# ​
#  =l が存在して、全ての
# k (1≤k≤M) について
# a
# y
# k
# ​

# ​
#  =X
# k
# ​

# 財政難に苦しむ AtCoder 諸島では、維持費削減のため橋を
# 1 本封鎖することになりました。 封鎖する橋をうまく選んだとき、ツアーの長さの最小値がいくつになるか求めてください。
if __name__ == "__main__":
    N, M = map(int, input().split())
    points = list(map(int, input().split()))
    for i in range(N):
        points[i] -= 1
