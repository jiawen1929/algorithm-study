/* eslint-disable no-inner-declarations */
/* eslint-disable no-cond-assign */
/* eslint-disable no-param-reassign */

// !单点修改+区间查询

class SegmentTreePointUpdateRangeQuery<E = number> {
  private readonly _n: number
  private readonly _size: number
  private readonly _data: E[]
  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E

  /**
   * 单点更新,区间查询的线段树.
   * @param nOrLeaves 大小或叶子节点的值.
   * @param e 幺元.
   * @param op 结合律.
   */
  constructor(nOrLeaves: number | ArrayLike<E>, e: () => E, op: (a: E, b: E) => E) {
    const n = typeof nOrLeaves === 'number' ? nOrLeaves : nOrLeaves.length
    let size = 1
    while (size < n) size <<= 1
    const data = Array(size << 1)
    for (let i = 0; i < data.length; i++) data[i] = e()

    this._n = n
    this._size = size
    this._data = data
    this._e = e
    this._op = op

    if (typeof nOrLeaves !== 'number') this.build(nOrLeaves)
  }

  set(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._data[index] = value
    while ((index >>= 1)) {
      this._data[index] = this._op(this._data[index << 1], this._data[(index << 1) | 1])
    }
  }

  get(index: number): E {
    if (index < 0 || index >= this._n) return this._e()
    return this._data[index + this._size]
  }

  /**
   * 将`index`处的值与作用素`value`结合.
   */
  update(index: number, value: E): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    this._data[index] = this._op(this._data[index], value)
    while ((index >>= 1)) {
      this._data[index] = this._op(this._data[index << 1], this._data[(index << 1) | 1])
    }
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): E {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return this._e()

    let leftRes = this._e()
    let rightRes = this._e()
    for (start += this._size, end += this._size; start < end; start >>= 1, end >>= 1) {
      if (start & 1) leftRes = this._op(leftRes, this._data[start++])
      if (end & 1) rightRes = this._op(this._data[--end], rightRes)
    }
    return this._op(leftRes, rightRes)
  }

  queryAll(): E {
    return this._data[1]
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (value: E) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    let res = this._e()
    while (true) {
      while (!(start & 1)) start >>= 1
      if (!predicate(this._op(res, this._data[start]))) {
        while (start < this._size) {
          start <<= 1
          if (predicate(this._op(res, this._data[start]))) {
            res = this._op(res, this._data[start])
            start++
          }
        }
        return start - this._size
      }
      res = this._op(res, this._data[start])
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (value: E) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    let res = this._e()
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      if (!predicate(this._op(this._data[end], res))) {
        while (end < this._size) {
          end = (end << 1) | 1
          if (predicate(this._op(this._data[end], res))) {
            res = this._op(this._data[end], res)
            end--
          }
        }
        return end + 1 - this._size
      }
      res = this._op(this._data[end], res)
      if ((end & -end) === end) break
    }
    return 0
  }

  build(arr: ArrayLike<E>): void {
    if (arr.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < arr.length; i++) {
      this._data[i + this._size] = arr[i] // 叶子结点
    }
    for (let i = this._size - 1; i > 0; i--) {
      this._data[i] = this._op(this._data[i << 1], this._data[(i << 1) | 1])
    }
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreePointUpdateRangeQuery(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(String(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }
}

export { SegmentTreePointUpdateRangeQuery }

if (require.main === module) {
  const seg = new SegmentTreePointUpdateRangeQuery(
    10,
    () => 0,
    (a, b) => a + b
  )
  console.log(seg.toString())
  seg.set(0, 1)
  seg.set(1, 2)
  console.log(seg.toString())
  seg.update(3, 4)
  console.log(seg.toString())
  console.log(seg.query(0, 4))
  seg.build([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  console.log(seg.toString())
  console.log(seg.minLeft(10, x => x < 15))
  console.log(seg.maxRight(0, x => x <= 15))
  console.log(seg.queryAll())

  benchMark()
  function benchMark(): void {
    const n = 2e5
    const seg = new SegmentTreePointUpdateRangeQuery<number>(
      n,
      () => 0,
      (parent, child) => parent + child
    )
    console.time('update')
    for (let i = 0; i < n; i++) {
      seg.update(i, i)
      seg.query(0, i)
    }
    console.timeEnd('update')
  }

  // https://leetcode.cn/problems/maximum-sum-queries/
  // 2736. 最大和查询
  function maximumSumQueries(nums1: number[], nums2: number[], queries: number[][]): number[] {}
}
