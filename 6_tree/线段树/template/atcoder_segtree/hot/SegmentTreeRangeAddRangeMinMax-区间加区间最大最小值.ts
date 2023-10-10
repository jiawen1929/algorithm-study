/* eslint-disable no-inner-declarations */

const INF = 2e15 // !超过int32使用2e15

/**
 * 区间加, 区间最大最小值.
 */
class SegmentTreeRangeAddRangeMinMax {
  private readonly _n: number
  private readonly _size: number
  private readonly _height: number
  private readonly _min: Float64Array
  private readonly _max: Float64Array
  private readonly _lazy: Float64Array

  constructor(nOrArr: number | ArrayLike<number>) {
    const n = typeof nOrArr === 'number' ? nOrArr : nOrArr.length
    let size = 1
    let height = 0
    while (size < n) {
      size <<= 1
      height++
    }
    this._n = n
    this._size = size
    this._height = height

    // !0.init data and lazy
    const min = new Float64Array(size << 1).fill(INF)
    const max = new Float64Array(size << 1).fill(-INF)
    const lazy = new Float64Array(size)
    this._min = min
    this._max = max
    this._lazy = lazy

    if (typeof nOrArr !== 'number') this._build(nOrArr)
  }

  set(index: number, value: number): void {
    if (index < 0 || index >= this._n) return
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    // !1. set
    this._min[index] = value
    this._max[index] = value
    for (let i = 1; i <= this._height; i++) this._pushUp(index >> i)
  }

  get(index: number): number {
    if (index < 0 || index >= this._n) {
      throw new RangeError(`index must be in [0, ${this._n})`)
    }
    index += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(index >> i)
    return this._max[index]
  }

  /**
   * 区间`[start,end)`的值与`lazy`进行作用.
   * 0 <= start <= end <= n.
   */
  update(start: number, end: number, lazy: number): void {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let start2 = start
    let end2 = end
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) this._propagate(start++, lazy)
      if (end & 1) this._propagate(--end, lazy)
    }
    start = start2
    end = end2
    for (let i = 1; i <= this._height; i++) {
      if ((start >> i) << i !== start) this._pushUp(start >> i)
      if ((end >> i) << i !== end) this._pushUp((end - 1) >> i)
    }
  }

  /**
   * 查询区间`[start,end)`的聚合值.
   * 0 <= start <= end <= n.
   */
  query(start: number, end: number): { min: number; max: number } {
    if (start < 0) start = 0
    if (end > this._n) end = this._n
    if (start >= end) return { min: INF, max: -INF }
    start += this._size
    end += this._size
    for (let i = this._height; i > 0; i--) {
      if ((start >> i) << i !== start) this._pushDown(start >> i)
      if ((end >> i) << i !== end) this._pushDown((end - 1) >> i)
    }
    let leftMin = INF
    let leftMax = -INF
    let rightMin = INF
    let rightMax = -INF
    for (; start < end; start >>= 1, end >>= 1) {
      if (start & 1) {
        leftMin = Math.min(leftMin, this._min[start])
        leftMax = Math.max(leftMax, this._max[start])
        start++
      }
      if (end & 1) {
        end--
        rightMin = Math.min(rightMin, this._min[end])
        rightMax = Math.max(rightMax, this._max[end])
      }
    }
    return { min: Math.min(leftMin, rightMin), max: Math.max(leftMax, rightMax) }
  }

  queryAll(): { min: number; max: number } {
    return { min: this._min[1], max: this._max[1] }
  }

  /**
   * 树上二分查询最大的`end`使得`[start,end)`内的值满足`predicate`.
   * @alias findFirst
   */
  maxRight(start: number, predicate: (min: number, max: number) => boolean): number {
    if (start < 0) start = 0
    if (start >= this._n) return this._n
    start += this._size
    for (let i = this._height; i > 0; i--) this._pushDown(start >> i)
    let resMin = INF
    let resMax = -INF
    while (true) {
      while (!(start & 1)) start >>= 1
      const tmpMin1 = Math.min(resMin, this._min[start])
      const tmpMax1 = Math.max(resMax, this._max[start])
      if (!predicate(tmpMin1, tmpMax1)) {
        while (start < this._size) {
          this._pushDown(start)
          start <<= 1
          const tmpMin2 = Math.min(resMin, this._min[start])
          const tmpMax2 = Math.max(resMax, this._max[start])
          if (predicate(tmpMin2, tmpMax2)) {
            resMin = tmpMin2
            resMax = tmpMax2
            start++
          }
        }
        return start - this._size
      }
      resMin = Math.min(resMin, this._min[start])
      resMax = Math.max(resMax, this._max[start])
      start++
      if ((start & -start) === start) break
    }
    return this._n
  }

  /**
   * 树上二分查询最小的`start`使得`[start,end)`内的值满足`predicate`
   * @alias findLast
   */
  minLeft(end: number, predicate: (min: number, max: number) => boolean): number {
    if (end > this._n) end = this._n
    if (end <= 0) return 0
    end += this._size
    for (let i = this._height; i > 0; i--) this._pushDown((end - 1) >> i)
    let resMin = INF
    let resMax = -INF
    while (true) {
      end--
      while (end > 1 && end & 1) end >>= 1
      const tmpMin1 = Math.min(resMin, this._min[end])
      const tmpMax1 = Math.max(resMax, this._max[end])
      if (!predicate(tmpMin1, tmpMax1)) {
        while (end < this._size) {
          this._pushDown(end)
          end = (end << 1) | 1
          const tmpMin2 = Math.min(resMin, this._min[end])
          const tmpMax2 = Math.max(resMax, this._max[end])
          if (predicate(tmpMin2, tmpMax2)) {
            resMin = tmpMin2
            resMax = tmpMax2
            end--
          }
        }
        return end + 1 - this._size
      }
      resMin = Math.min(resMin, this._min[end])
      resMax = Math.max(resMax, this._max[end])
      if ((end & -end) === end) break
    }
    return 0
  }

  toString(): string {
    const sb: string[] = []
    sb.push('SegmentTreeRangeUpdateRangeQuery(')
    for (let i = 0; i < this._n; i++) {
      if (i) sb.push(', ')
      sb.push(JSON.stringify(this.get(i)))
    }
    sb.push(')')
    return sb.join('')
  }

  private _build(leaves: ArrayLike<number>): void {
    if (leaves.length !== this._n) throw new RangeError(`length must be equal to ${this._n}`)
    for (let i = 0; i < this._n; i++) {
      this._min[this._size + i] = leaves[i]
      this._max[this._size + i] = leaves[i]
    }
    for (let i = this._size - 1; i > 0; i--) this._pushUp(i)
  }

  private _pushUp(index: number): void {
    this._min[index] = Math.min(this._min[index << 1], this._min[(index << 1) | 1])
    this._max[index] = Math.max(this._max[index << 1], this._max[(index << 1) | 1])
  }

  private _pushDown(index: number): void {
    const lazy = this._lazy[index]
    if (!lazy) return
    this._propagate(index << 1, lazy)
    this._propagate((index << 1) | 1, lazy)
    this._lazy[index] = 0
  }

  private _propagate(index: number, lazy: number): void {
    this._min[index] += lazy
    this._max[index] += lazy
    if (index < this._size) this._lazy[index] += lazy
  }
}

export { SegmentTreeRangeAddRangeMinMax }

if (require.main === module) {
  // checkWithBruteForce()
  timeit()

  function checkWithBruteForce(): void {
    class Mocker {
      readonly _n: number
      private readonly _a: number[]
      constructor(nums: number[]) {
        this._n = nums.length
        this._a = nums.slice()
      }

      set(index: number, value: number): void {
        this._a[index] = value
      }

      get(index: number): number {
        return this._a[index]
      }

      update(start: number, end: number, lazy: number): void {
        for (let i = start; i < end; i++) this._a[i] += lazy
      }

      query(start: number, end: number): { min: number; max: number } {
        let min = INF
        let max = -INF

        for (let i = start; i < end; i++) {
          min = Math.min(min, this._a[i])
          max = Math.max(max, this._a[i])
        }
        return { min, max }
      }

      queryAll(): { min: number; max: number } {
        return this.query(0, this._n)
      }

      maxRight(start: number, predicate: (min: number, max: number) => boolean): number {
        let min = INF
        let max = -INF

        for (let i = start; i < this._n; i++) {
          min = Math.min(min, this._a[i])
          max = Math.max(max, this._a[i])

          if (!predicate(min, max)) return i
        }
        return this._n
      }

      minLeft(end: number, predicate: (min: number, max: number) => boolean): number {
        let min = INF
        let max = -INF

        for (let i = end - 1; i >= 0; i--) {
          min = Math.min(min, this._a[i])
          max = Math.max(max, this._a[i])

          if (!predicate(min, max)) return i + 1
        }
        return 0
      }

      build(leaves: ArrayLike<number>): void {
        for (let i = 0; i < this._a.length; i++) this._a[i] = leaves[i]
      }

      toString(): string {
        return `Mocker(${this._a})`
      }
    }
    function assertSame(obj1: unknown, obj2: unknown) {
      if (JSON.stringify(obj1) !== JSON.stringify(obj2)) {
        throw new Error(`expect ${JSON.stringify(obj2)}, got ${JSON.stringify(obj1)}`)
      }
    }

    const randint = (min: number, max: number) => Math.floor(Math.random() * (max - min + 1)) + min
    const N = 5e4
    const real = new SegmentTreeRangeAddRangeMinMax(Array(N).fill(0))
    const mock = new Mocker(Array(N).fill(0))
    for (let i = 0; i < N; i++) {
      const op = randint(0, 5)
      if (op === 0) {
        // set
        const index = randint(0, N - 1)
        const value = randint(0, 10)
        real.set(index, value)
        mock.set(index, value)
        // console.log('set', index, value)
      } else if (op === 1) {
        // get
        const index = randint(0, N - 1)
        const realValue = real.get(index)
        const mockValue = mock.get(index)
        // console.log(realValue, mockValue, index)
        // console.log('get', index, realValue, mockValue)
        assertSame(realValue, mockValue)
      } else if (op === 2) {
        // update
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const lazy = randint(0, 2)
        real.update(start, end, lazy)
        mock.update(start, end, lazy)
        // console.log('update', start, end, lazy)
      } else if (op === 3) {
        // query
        const start = randint(0, N - 1)
        const end = randint(start, N)
        const realValue = real.query(start, end)
        const mockValue = mock.query(start, end)
        // console.log('query', start, end, realValue, mockValue)
        assertSame(realValue, mockValue)
      } else if (op === 4) {
        // queryAll
        const realValue = real.queryAll()
        const mockValue = mock.queryAll()
        assertSame(realValue, mockValue)
      } else if (op === 5) {
        // maxRight
        const start = randint(0, N - 1)
        const target = randint(0, N)
        const realValue = real.maxRight(start, min => min >= target)
        const mockValue = mock.maxRight(start, min => min >= target)
        assertSame(realValue, mockValue)
      } else if (op === 6) {
        // minLeft
        const end = randint(0, N)
        const target = randint(0, N)
        const realValue = real.minLeft(end, min => min >= target)
        const mockValue = mock.minLeft(end, min => min >= target)
        assertSame(realValue, mockValue)
      }
    }
    console.log('test passed')
  }

  function timeit(): void {
    const n = 2e5
    const arr = Array(n)
    for (let i = 0; i < n; i++) arr[i] = Math.floor(Math.random() * 10)
    const seg = new SegmentTreeRangeAddRangeMinMax(arr)
    console.time('SegmentTreeRangeAddRangeMinMax')
    for (let i = 0; i < n; i++) {
      seg.query(i, n)
      seg.update(i, n, 1)
      seg.set(i, 1)
      seg.maxRight(i, min => min >= i)
      seg.minLeft(i, min => min >= i)
    }
    console.timeEnd('SegmentTreeRangeAddRangeMinMax') // SegmentTreeRangeAddRangeMinMax: 227.276ms
  }
}
