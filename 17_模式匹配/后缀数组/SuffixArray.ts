/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */

// - 支持sa/rank/lcp
// - 比较任意两个子串的字典序
// - 求出任意两个子串的最长公共前缀(lcp)
//
// sa : 排第几的后缀是谁.
// rank : 每个后缀排第几.
// lcp : 排名相邻的两个后缀的最长公共前缀.
// lcp[0] = 0
// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
//
// "banana" -> sa: [5 3 1 0 4 2], rank: [3 2 5 1 4 0], lcp: [0 1 3 0 0 2]

import { SparseTableUint32 } from '../../22_专题/RMQ问题/SparseTable'

/**
 * 后缀数组.
 */
class SuffixArray {
  readonly sa: number[]
  readonly rank: number[]
  readonly height: number[]
  private _queryMin?: { query: (start: number, end: number) => number }
  private readonly _n: number

  /**
   * @param sOrOrds 字符串或表示字符码的`非负整数`数组.
   * 当`ord`很大时(>1e7),需要对数组进行离散化,减少内存占用.
   */
  constructor(sOrOrds: string | ArrayLike<number>) {
    if (typeof sOrOrds === 'string') {
      const ords = Array(sOrOrds.length)
      for (let i = 0; i < ords.length; i++) ords[i] = sOrOrds.charCodeAt(i)
      sOrOrds = ords
    }
    const sa = SuffixArray._saIs(sOrOrds)
    const [rank, lcp] = SuffixArray._rankLcp(sOrOrds, sa)
    this.sa = sa
    this.rank = rank
    this.height = lcp
    this._n = sOrOrds.length
  }

  /**
   * 求任意两个子串`s[a,b)`和`s[c,d)`的最长公共前缀(lcp).
   * 0 <= a <= b <= n, 0 <= c <= d <= n.
   */
  lcp(a: number, b: number, c: number, d: number): number {
    if (a >= b || c >= d) return 0
    return Math.min(b - a, d - c, this._lcp(a, c))
  }

  /**
   * 比较任意两个子串`s[a,b)`和`s[c,d)`的字典序.
   * 0 <= a <= b <= n, 0 <= c <= d <= n.
   * ```
   * s[a,b) < s[c,d) => -1
   * s[a,b) = s[c,d) => 0
   * s[a,b) > s[c,d) => 1
   * ```
   */
  compareSubstr(a: number, b: number, c: number, d: number): -1 | 0 | 1 {
    const len1 = b - a
    const len2 = d - c
    const lcp = this._lcp(a, c)
    if (len1 === len2 && lcp >= len1) return 0
    if (lcp >= len1 || lcp >= len2) return len1 < len2 ? -1 : 1
    return this.rank[a] < this.rank[c] ? -1 : 1
  }

  private _lcp(i: number, j: number): number {
    if (!this._queryMin) this._queryMin = new SparseTableUint32(this.height, () => 0, Math.min)
    if (i === j) return this._n - i
    let r1 = this.rank[i]
    let r2 = this.rank[j]
    if (r1 > r2) {
      r1 ^= r2
      r2 ^= r1
      r1 ^= r2
    }
    return this._queryMin.query(r1 + 1, r2 + 1)
  }

  /** 基于sais诱导排序算法线性时间构建sa数组. */
  private static _saIs(ords: ArrayLike<number>): number[] {
    const n = ords.length
    let max = 0
    for (let i = 0; i < n; i++) max = Math.max(max, ords[i])
    const buckets = new Uint32Array(max + 2)
    for (let i = 0; i < n; i++) buckets[ords[i] + 1]++
    for (let i = 1; i < buckets.length; i++) buckets[i] += buckets[i - 1]
    const isL = new Uint8Array(n)
    isL[n - 1] = 1
    for (let i = n - 2; ~i; i--) {
      if (ords[i] === ords[i + 1]) {
        isL[i] = isL[i + 1]
      } else {
        isL[i] = +(ords[i] > ords[i + 1])
      }
    }

    const isLMS = new Uint8Array(n + 1)
    for (let i = 1; i < n; i++) isLMS[i] = +(isL[i - 1] && !isL[i])
    isLMS[n] = 1
    let lms1: number[] = []
    for (let i = 0; i < n; i++) if (isLMS[i]) lms1.push(i)
    if (lms1.length > 1) {
      const sa = inducedSort(lms1)
      const lms2: number[] = []
      for (let i = 0; i < sa.length; i++) {
        const v = sa[i]
        if (isLMS[v]) lms2.push(v)
      }
      let pre = -1
      let j = 0
      for (let k = 0; k < lms2.length; k++) {
        const v = lms2[k]
        let i1 = pre
        let i2 = v
        while (pre >= 0 && ords[i1] === ords[i2]) {
          i1++
          i2++
          if (isLMS[i1] || isLMS[i2]) {
            j -= isLMS[i1] && isLMS[i2]
            break
          }
        }
        j++
        pre = v
        sa[v] = j
      }
      const tmp = Array(lms1.length)
      for (let i = 0; i < tmp.length; i++) tmp[i] = sa[lms1[i]]
      const newSa = this._saIs(tmp)
      const nextLms1 = Array(newSa.length)
      for (let i = 0; i < newSa.length; i++) nextLms1[i] = lms1[newSa[i]]
      lms1 = nextLms1
    }

    return inducedSort(lms1)

    function inducedSort(lms: ArrayLike<number>): number[] {
      const sa = Array(n + 1)
      for (let i = 0; i < n; i++) sa[i] = -1
      sa[n] = n
      let endpoint = buckets.slice(1)
      for (let i = lms.length - 1; ~i; i--) {
        const v = lms[i]
        endpoint[ords[v]]--
        sa[endpoint[ords[v]]] = v
      }
      const startpoint = buckets.slice(0, -1)
      const v = sa[n] - 1
      if (v >= 0 && isL[v]) {
        sa[startpoint[ords[v]]] = v
        startpoint[ords[v]]++
      }
      for (let i = 0; i < n; i++) {
        const v = sa[i] - 1
        if (v >= 0 && isL[v]) {
          sa[startpoint[ords[v]]] = v
          startpoint[ords[v]]++
        }
      }
      sa.pop()
      endpoint = buckets.slice(1)
      for (let i = n - 1; ~i; i--) {
        const v = sa[i] - 1
        if (v >= 0 && !isL[v]) {
          endpoint[ords[v]]--
          sa[endpoint[ords[v]]] = v
        }
      }
      return sa
    }
  }

  private static _rankLcp(
    ords: ArrayLike<number>,
    sa: ArrayLike<number>
  ): [rank: number[], lcp: number[]] {
    const n = ords.length
    const rank = Array(n)
    const lcp = Array(n)
    for (let i = 0; i < n; i++) {
      rank[sa[i]] = i
      lcp[i] = 0
    }
    let h = 0
    for (let i = 0; i < n; i++) {
      if (h > 0) h--
      if (!rank[i]) continue
      const j = sa[rank[i] - 1]
      while (j + h < n && i + h < n) {
        if (ords[j + h] !== ords[i + h]) break
        h++
      }
      lcp[rank[i]] = h
    }
    return [rank, lcp]
  }
}

export { SuffixArray }

if (require.main === module) {
  const sa = new SuffixArray('banana')
  console.log(sa.sa, sa.rank, sa.height)
  // https://leetcode.cn/problems/sum-of-scores-of-built-strings/
  function sumScores(s: string): number {
    const sa = new SuffixArray(s)
    const n = s.length
    let res = 0
    for (let i = 0; i < n; i++) {
      res += sa.lcp(0, n, i, n)
    }
    return res
  }

  const n = 100
  const ords = Array.from({ length: n }, () => 100 + Math.floor(Math.random() * 26))
  console.time('test')
  let count = 0
  const test = new SuffixArray(ords)
  // a,b,c,d
  for (let a = 0; a < n; a++) {
    for (let b = a; b < n; b++) {
      for (let c = 0; c < n; c++) {
        for (let d = c; d < n; d++) {
          count++
          const lcp = test.lcp(a, b, c, d)
          const lcp2 = lcpNaive(ords, a, b, c, d)
          if (lcp !== lcp2) {
            // console.log(a, b, c, d, lcp, lcp2)
            throw new Error(`${a} ${b} ${c} ${d} ${lcp} ${lcp2}`)
          }

          const cmp = test.compareSubstr(a, b, c, d)
          const cmp2 = compareSubstrNaive(ords, a, b, c, d)
          if (cmp !== cmp2) {
            throw new Error(`${a} ${b} ${c} ${d} ${cmp} ${cmp2}`)
          }
        }
      }
    }
  }

  console.log('ok')
  console.timeEnd('test')
  console.log(count)

  function lcpNaive(s: ArrayLike<number>, a: number, b: number, c: number, d: number): number {
    let res = 0
    while (a < b && c < d && s[a] === s[c]) {
      a++
      c++
      res++
    }
    return res
  }

  function compareSubstrNaive(
    s: ArrayLike<number>,
    a: number,
    b: number,
    c: number,
    d: number
  ): -1 | 0 | 1 {
    while (a < b && c < d && s[a] === s[c]) {
      a++
      c++
    }
    if (a === b) return c === d ? 0 : -1
    if (c === d) return 1
    return s[a] < s[c] ? -1 : 1
  }
}
