# 求组合数
from functools import lru_cache

MOD = int(1e9 + 7)
# 1. 逆元最快


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


def C(n: int, k: int) -> int:
    if n < k:
        return 0
    return (fac(n) * ifac(k) * ifac(n - k)) % MOD


if __name__ == '__main__':
    print(C(n=3, k=3))
    print(C(n=4, k=3))
    print(C(n=5, k=3))

#########################################################
# 预处理组合数 C(n,k)=C(n-1,k)+C(n-1,k-1)
# 不太快
comb = [[0] * 36 for _ in range(36)]
for i in range(36):
    comb[i][0] = 1
    for j in range(1, i + 1):
        comb[i][j] = comb[i - 1][j - 1] + comb[i - 1][j]

print(comb[10][2])

#########################################################
# 不太快
@lru_cache(None)
def C1(n: int, k: int) -> int:
    if n < k:
        return 0
    if n == 1 or k == 0:
        return 1
    return (C1(n - 1, k) + C1(n - 1, k - 1)) % MOD

