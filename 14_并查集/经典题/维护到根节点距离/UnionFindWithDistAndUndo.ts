/* eslint-disable max-len */

import { RollbackArray } from '../../../0_数组/RollbackArray'

/**
 * 维护到根节点距离的可撤销并查集.
 * 用于维护环的权值，树上的距离等.
 */
class UnionFindWithDistAndUndo<D> {
  private readonly _n: number
  private readonly _e: () => D
  private readonly _op: (x: D, y: D) => D
  private readonly _inv: (x: D) => D
  private readonly _data: RollbackArray<{ parent: number; dist: D }>

  constructor(n: number, monoid: { e: () => D; op: (x: D, y: D) => D; inv: (x: D) => D } & ThisType<void>) {
    this._n = n
    this._e = monoid.e
    this._op = monoid.op
    this._inv = monoid.inv
    this._data = new RollbackArray(n, () => ({ parent: -1, dist: this._e() }))
  }

  /**
   * `distToRoot(parent) + dist = distToRoot(child)`.
   * @returns 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
   */
  union(parent: number, child: number, dist: D): boolean {
    let { groupRoot: v1, distToRoot: x1 } = this.find(parent)
    let { groupRoot: v2, distToRoot: x2 } = this.find(child)
    if (v1 === v2) {
      return dist === this._e()
    }
    let s1 = -this._data.get(v1).parent
    let s2 = -this._data.get(v2).parent
    if (s1 < s2) {
      const tmp = v1
      v1 = v2
      v2 = tmp
      const tmp2 = x1
      x1 = x2
      x2 = tmp2
      dist = this._inv(dist)
    }
    // v1 <- v2
    dist = this._op(x1, dist)
    dist = this._op(dist, this._inv(x2))
    this._data.set(v2, { parent: v1, dist })
    this._data.set(v1, { parent: -(s1 + s2), dist: this._e() })
    return true
  }

  find(x: number): { groupRoot: number; distToRoot: D } {
    let root = x
    let distToRoot = this._e()
    while (true) {
      const { parent, dist } = this._data.get(root)
      if (parent < 0) {
        break
      }
      distToRoot = this._op(distToRoot, dist)
      root = parent
    }
    return { groupRoot: root, distToRoot }
  }

  /**
   * 返回x到y的距离`f(x) - f(y)`.
   * @throws 如果x和y不在同一个集合,抛出错误.
   */
  dist(x: number, y: number): D {
    const { groupRoot: vx, distToRoot: dx } = this.find(x)
    const { groupRoot: vy, distToRoot: dy } = this.find(y)
    if (vx !== vy) {
      throw new Error('x and y are not in the same set')
    }
    return this._op(dx, this._inv(dy))
  }

  distToRoot(x: number): D {
    return this.find(x).distToRoot
  }

  getTime(): number {
    return this._data.getTime()
  }

  rollback(time: number): void {
    this._data.rollback(time)
  }

  undo(): boolean {
    return this._data.undo()
  }

  getSize(x: number): number {
    return -this._data.get(this.find(x).groupRoot).parent
  }

  getGroups(): Map<number, number[]> {
    const res = new Map<number, number[]>()
    for (let i = 0; i < this._n; i++) {
      const { groupRoot } = this.find(i)
      if (!res.has(groupRoot)) res.set(groupRoot, [])
      res.get(groupRoot)!.push(i)
    }
    return res
  }
}

export { UnionFindWithDistAndUndo }
