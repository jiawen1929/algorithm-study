from typing import Optional, Tuple


def exgcd(a: int, b: int) -> Tuple[int, int, int]:
    """
    求a, b最大公约数,同时求出裴蜀定理中的一组系数x, y,
    满足 x*a + y*b = gcd(a, b)

    ax + by = gcd_ 返回 `(gcd_, x, y)`

    """
    if b == 0:
        return a, 1, 0
    gcd_, x, y = exgcd(b, a % b)
    return gcd_, y, x - a // b * y


def modInv(a: int, mod: int) -> Optional[int]:
    """
    扩展gcd求a在mod下的逆元
    即求出逆元 `inv` 满足 `a*inv ≡ 1 (mod m)`
    """
    gcd_, x, _ = exgcd(a, mod)
    if gcd_ != 1:
        return None
    return x % mod


def rationalMod(a: int, b: int, mod: int) -> Optional[int]:
    """
    有理数取模(有理数取余)
    求 a/b 模 mod 的值
    """
    inv = modInv(b, mod)
    if inv is None:
        return None
    return a * inv % mod


def exgcdFarey(a: int, b: int) -> Tuple[int, int, int]:
    """ax + by = gcd_ 返回 `(gcd_, x, y)`

    满足解最小, 且 (abs(x)+abs(y), x) 字典序最小
    """
    x1, y1, x2, y2 = farey(a, b)
    x1, y1 = y1, -x1
    x2, y2 = -y2, x2
    g = a * x1 + b * y1
    key1 = (abs(x1) + abs(y1), x1)
    key2 = (abs(x2) + abs(y2), x2)
    if key1 < key2:
        return g, x1, y1
    return g, x2, y2


# Farey 数列 中 a/b 第一次出现的位置的前驱和后继
# a/b = 19/12 → (x1/y1, x2/y2) = (11/7, 8/5) → 返回 (11,7,8,5)
def farey(a: int, b: int) -> Tuple[int, int, int, int]:
    """
    求法雷数列中某一项的的前驱和后继
    https://zhuanlan.zhihu.com/p/323538981
    """
    assert a > 0 and b > 0
    if a == b:
        return 0, 1, 1, 0
    q = (a - 1) // b
    x1, y1, x2, y2 = farey(b, a - q * b)
    return q * x2 + y2, x2, q * x1 + y1, x1


if __name__ == "__main__":
    assert exgcd(2, 3) == (1, -1, 1)
    assert modInv(2, 998244353) == (998244353 + 1) // 2

    # Rational Approximation
    # https://yukicoder.me/problems/no/1936
    p, q = map(int, input().split())
    a, b, c, d = farey(p, q)
    print(a + b + c + d)
