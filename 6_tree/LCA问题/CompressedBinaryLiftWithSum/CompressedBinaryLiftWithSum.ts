interface MutableArrayLike<T> {
  readonly length: number
  [n: number]: T
}

/**
 * 空间复杂度`O(n)`的树上倍增，用于倍增结构优化建图、查询路径聚合值.
 *
 * @see
 * - https://taodaling.github.io/blog/2020/03/18/binary-lifting/
 * - https://codeforces.com/blog/entry/74847
 * - https://codeforces.com/blog/entry/100826
 */
class CompressedBinaryLiftWithSum<S = number> {
  readonly depth: MutableArrayLike<number>
  readonly parent: MutableArrayLike<number>
  private readonly _jump: MutableArrayLike<number>

  /**
   * 从当前结点到`jump`结点的路径上的聚合值(不包含`jump`结点).
   */
  private readonly _attachments: S[]

  /**
   * 当前结点的聚合值.
   */
  private readonly _singles: S[]
  private readonly _e: () => S
  private readonly _op: (e1: S, e2: S) => S

  /**
   * @param values
   * 每个点的`点权`.
   * !如果需要查询边权，则每个点的`点权`设为`该点与其父亲结点的边权`，根节点的`点权`设为`幺元`.
   */
  constructor(
    tree: MutableArrayLike<MutableArrayLike<number>>,
    values: (index: number) => S,
    monoid: {
      e: () => S
      op: (e1: S, e2: S) => S
    },
    root?: number
  )
  constructor(
    n: number,
    depthOnTree: MutableArrayLike<number>,
    parentOnTree: MutableArrayLike<number>,
    values: (index: number) => S,
    monoid: {
      e: () => S
      op: (e1: S, e2: S) => S
    }
  )
  constructor(arg0: any, arg1: any, arg2: any, arg3?: any, arg4?: any) {
    if (arguments.length === 5) {
      const n = arg0
      this.depth = arg1
      this.parent = arg2
      this._jump = new Int32Array(n).fill(-1)
      this._attachments = Array(n)
      this._singles = Array(n)
      this._e = arg4.e
      this._op = arg4.op
      for (let i = 0; i < n; i++) {
        this._attachments[i] = this._e()
        this._singles[i] = arg3(i)
      }

      for (let i = 0; i < n; i++) this._consider(i)
    } else {
      const n = arg0.length
      if (arg3 == undefined) arg3 = 0
      this.depth = new Int32Array(n)
      this.parent = new Int32Array(n)
      this.parent[arg3] = -1
      this._jump = new Int32Array(n)
      this._jump[arg3] = arg3
      this._attachments = Array(n)
      this._singles = Array(n)
      this._e = arg2.e
      this._op = arg2.op
      for (let i = 0; i < n; i++) {
        this._attachments[i] = this._e()
        this._singles[i] = arg1(i)
      }

      this._setUp(arg0, arg3)
    }
  }

  firstTrue = (start: number, predicate: (end: number) => boolean): number => {
    while (!predicate(start)) {
      if (predicate(this._jump[start])) {
        start = this.parent[start]
      } else {
        if (start === this._jump[start]) return -1
        start = this._jump[start]
      }
    }
    return start
  }

  firstTrueWithSum = (
    start: number,
    predicate: (end: number, sum: S) => boolean,
    isEdge: boolean
  ): { node: number; sum: S } => {
    if (isEdge) {
      let sum = this._e() // 不包含_singles[start]
      while (true) {
        if (predicate(start, sum)) {
          return { node: start, sum }
        }

        const jumpStart = this._jump[start]
        const jumpSum = this._op(sum, this._attachments[start])
        if (predicate(jumpStart, jumpSum)) {
          sum = this._op(sum, this._singles[start])
          start = this.parent[start]
        } else {
          if (start === jumpStart) {
            return { node: -1, sum: jumpSum }
          }
          sum = jumpSum
          start = jumpStart
        }
      }
    } else {
      let sum = this._e() // 不包含_singles[start]
      while (true) {
        const sumWithSingle = this._op(sum, this._singles[start])
        if (predicate(start, sumWithSingle)) {
          return { node: start, sum: sumWithSingle }
        }

        const jumpStart = this._jump[start]
        const jumpSum1 = this._op(sum, this._attachments[start])
        const jumpSum2 = this._op(jumpSum1, this._singles[jumpStart])
        if (predicate(jumpStart, jumpSum2)) {
          sum = sumWithSingle
          start = this.parent[start]
        } else {
          if (start === jumpStart) {
            return { node: -1, sum: jumpSum2 }
          }
          sum = jumpSum1
          start = jumpStart
        }
      }
    }
  }

  lastTrue = (start: number, predicate: (end: number) => boolean): number => {
    if (!predicate(start)) return -1
    while (true) {
      if (predicate(this._jump[start])) {
        if (start === this._jump[start]) return start
        start = this._jump[start]
      } else if (predicate(this.parent[start])) {
        start = this.parent[start]
      } else {
        return start
      }
    }
  }

  lastTrueWithSum = (
    start: number,
    predicate: (end: number, sum: S) => boolean,
    isEdge: boolean
  ): { node: number; sum: S } => {
    if (isEdge) {
      let sum = this._e() // 不包含_singles[start]
      if (!predicate(start, sum)) {
        return { node: -1, sum }
      }

      while (true) {
        const jumpStart = this._jump[start]
        const jumpSum = this._op(sum, this._attachments[start])
        if (predicate(jumpStart, jumpSum)) {
          if (start === jumpStart) {
            return { node: start, sum }
          }

          sum = jumpSum
          start = jumpStart
        } else {
          const parentStart = this.parent[start]
          const parentSum = this._op(sum, this._singles[start])
          if (predicate(parentStart, parentSum)) {
            sum = parentSum
            start = parentStart
          } else {
            return { node: start, sum }
          }
        }
      }
    } else {
      if (!predicate(start, this._singles[start])) {
        return { node: -1, sum: this._singles[start] }
      }

      let sum = this._e() // 不包含_singles[start]
      while (true) {
        const jumpStart = this._jump[start]
        const jumpSum1 = this._op(sum, this._attachments[start])
        const jumpSum2 = this._op(jumpSum1, this._singles[jumpStart])
        if (predicate(jumpStart, jumpSum2)) {
          if (start === jumpStart) {
            return { node: start, sum: jumpSum2 }
          }

          sum = jumpSum1
          start = jumpStart
        } else {
          const parentStart = this.parent[start]
          const parentSum1 = this._op(sum, this._singles[start])
          const parentSum2 = this._op(parentSum1, this._singles[parentStart])
          if (predicate(parentStart, parentSum2)) {
            sum = parentSum1
            start = parentStart
          } else {
            return { node: start, sum: parentSum1 }
          }
        }
      }
    }
  }

  upToDepth = (root: number, toDepth: number): number => {
    if (!(toDepth >= 0 && toDepth <= this.depth[root])) return -1
    if (this.depth[root] < toDepth) return -1
    while (this.depth[root] > toDepth) {
      if (this.depth[this._jump[root]] < toDepth) {
        root = this.parent[root]
      } else {
        root = this._jump[root]
      }
    }
    return root
  }

  upToDepthWithSum = (root: number, toDepth: number, isEdge: boolean): { node: number; sum: S } => {
    let sum = this._e() // 不包含_singles[root]
    if (!(toDepth >= 0 && toDepth <= this.depth[root])) return { node: -1, sum }
    while (this.depth[root] > toDepth) {
      if (this.depth[this._jump[root]] < toDepth) {
        sum = this._op(sum, this._singles[root])
        root = this.parent[root]
      } else {
        sum = this._op(sum, this._attachments[root])
        root = this._jump[root]
      }
    }
    if (!isEdge) {
      sum = this._op(sum, this._singles[root])
    }
    return { node: root, sum }
  }

  kthAncestor = (node: number, k: number): number => {
    const targetDepth = this.depth[node] - k
    return this.upToDepth(node, targetDepth)
  }

  kthAncestorWithSum = (node: number, k: number, isEdge: boolean): { node: number; sum: S } => {
    const targetDepth = this.depth[node] - k
    return this.upToDepthWithSum(node, targetDepth, isEdge)
  }

  lca = (a: number, b: number): number => {
    if (this.depth[a] > this.depth[b]) {
      a = this.kthAncestor(a, this.depth[a] - this.depth[b])
    } else if (this.depth[a] < this.depth[b]) {
      b = this.kthAncestor(b, this.depth[b] - this.depth[a])
    }
    while (a !== b) {
      if (this._jump[a] === this._jump[b]) {
        a = this.parent[a]
        b = this.parent[b]
      } else {
        a = this._jump[a]
        b = this._jump[b]
      }
    }
    return a
  }

  /**
   * 查询路径`a`到`b`的聚合值.
   * @param isEdge 是否是边权.
   */
  lcaWithSum = (a: number, b: number, isEdge: boolean): { node: number; sum: S } => {
    let e: S // 不包含_singles[a]和_singles[b]
    if (this.depth[a] > this.depth[b]) {
      const { node: end, sum } = this.upToDepthWithSum(a, this.depth[b], true)
      a = end
      e = sum
    } else if (this.depth[a] < this.depth[b]) {
      const { node: end, sum } = this.upToDepthWithSum(b, this.depth[a], true)
      b = end
      e = sum
    } else {
      e = this._e()
    }

    while (a !== b) {
      if (this._jump[a] === this._jump[b]) {
        e = this._op(e, this._singles[a])
        e = this._op(e, this._singles[b])
        a = this.parent[a]
        b = this.parent[b]
      } else {
        e = this._op(e, this._attachments[a])
        e = this._op(e, this._attachments[b])
        a = this._jump[a]
        b = this._jump[b]
      }
    }

    if (!isEdge) {
      e = this._op(e, this._singles[a])
    }
    return { node: a, sum: e }
  }

  jump = (start: number, target: number, step: number): number => {
    const lca = this.lca(start, target)
    const dep1 = this.depth[start]
    const dep2 = this.depth[target]
    const deplca = this.depth[lca]
    const dist = dep1 + dep2 - 2 * deplca
    if (step > dist) return -1
    if (step <= dep1 - deplca) return this.kthAncestor(start, step)
    return this.kthAncestor(target, dist - step)
  }

  dist = (a: number, b: number): number => {
    return this.depth[a] + this.depth[b] - 2 * this.depth[this.lca(a, b)]
  }

  inSubtree = (maybeChild: number, maybeAncestor: number): boolean => {
    return (
      this.depth[maybeChild] >= this.depth[maybeAncestor] &&
      this.kthAncestor(maybeChild, this.depth[maybeChild] - this.depth[maybeAncestor]) ===
        maybeAncestor
    )
  }

  private _consider = (root: number): void => {
    if (root === -1 || this._jump[root] !== -1) return
    const p = this.parent[root]
    this._consider(p)
    this._addLeaf(root, p)
  }

  private _addLeaf = (leaf: number, parent: number): void => {
    if (parent == -1) {
      this._jump[leaf] = leaf
    } else {
      const tmp = this._jump[parent]
      if (this.depth[parent] - this.depth[tmp] === this.depth[tmp] - this.depth[this._jump[tmp]]) {
        this._jump[leaf] = this._jump[tmp]
        this._attachments[leaf] = this._op(this._singles[leaf], this._attachments[parent])
        this._attachments[leaf] = this._op(this._attachments[leaf], this._attachments[tmp])
      } else {
        this._jump[leaf] = parent
        this._attachments[leaf] = this._singles[leaf] // TODO: copy
      }
    }
  }

  private _setUp = (tree: ArrayLike<ArrayLike<number>>, root: number): void => {
    const queue: number[] = [root]
    let head = 0
    while (head < queue.length) {
      const cur = queue[head++]
      const nexts = tree[cur]
      for (let i = 0; i < nexts.length; i++) {
        const next = nexts[i]
        if (next === this.parent[cur]) continue
        this.depth[next] = this.depth[cur] + 1
        this.parent[next] = cur
        queue.push(next)
        this._addLeaf(next, cur)
      }
    }
  }
}

export { CompressedBinaryLiftWithSum }

if (require.main === module) {
  const n = 7
  const edges = [
    [0, 1],
    [0, 2],
    [1, 3],
    [1, 4],
    [2, 5],
    [4, 6]
  ]

  //          0
  //        /   \
  //       1     2
  //      / \     \
  //     3   4     5
  //         /
  //        6

  const tree: number[][] = Array(n)
  for (let i = 0; i < n; i++) tree[i] = []
  edges.forEach(([u, v]) => {
    tree[u].push(v)
    tree[v].push(u)
  })

  const values: number[] = [1, 1, 2, 3, 4, 5, 6]
  const bl = new CompressedBinaryLiftWithSum(tree, i => values[i], {
    e: () => 0,
    op: (a, b) => a + b
  })

  // https://leetcode.cn/problems/kth-ancestor-of-a-tree-node/
  class TreeAncestor {
    private readonly _lca: CompressedBinaryLiftWithSum

    constructor(n: number, parent: number[]) {
      const adjList: number[][] = Array(n)
      for (let i = 0; i < n; i++) adjList[i] = []
      parent.forEach((p, i) => {
        if (p !== -1) adjList[p].push(i)
      })
      this._lca = new CompressedBinaryLiftWithSum(adjList, i => 0, {
        e: () => 0,
        op: (a, b) => a + b
      })
    }

    getKthAncestor(node: number, k: number): number {
      return this._lca.kthAncestorWithSum(node, k, true).node
    }
  }

  // test with sum api

  // console.log(
  //   bl.firstTrueWithSum(
  //     6,
  //     (i, sum) => {
  //       console.log(i, sum, 'test')
  //       return sum >= 16
  //     },
  //     true
  //   )
  // )

  // console.log(
  //   bl.lastTrueWithSum(
  //     6,
  //     (i, sum) => {
  //       console.log(i, sum, 'test')
  //       return sum <= 10
  //     },
  //     false
  //   )
  // )

  // console.log(bl.upToDepthWithSum(6, 0, false))
  // console.log(bl.upToDepthWithSum(6, 1, false))
  // console.log(bl.upToDepthWithSum(6, 2, false))
  // console.log(bl.upToDepthWithSum(6, 3, false))
  // console.log(bl.upToDepthWithSum(6, 0, true))
  // console.log(bl.upToDepthWithSum(6, 1, true))
  // console.log(bl.upToDepthWithSum(6, 2, true))
  // console.log(bl.upToDepthWithSum(6, 3, true))

  // console.log(bl.kthAncestorWithSum(6, 0, true))
  // console.log(bl.kthAncestorWithSum(6, 1, true))
  // console.log(bl.kthAncestorWithSum(6, 2, true))
  // console.log(bl.kthAncestorWithSum(6, 3, true))
  // console.log(bl.kthAncestorWithSum(6, 4, true))

  // console.log(bl.kthAncestorWithSum(6, 0, false))
  // console.log(bl.kthAncestorWithSum(6, 1, false))
  // console.log(bl.kthAncestorWithSum(6, 2, false))
  // console.log(bl.kthAncestorWithSum(6, 3, false))
  // console.log(bl.kthAncestorWithSum(6, 4, false))

  // console.log(bl.lcaWithSum(3, 5, false))
  // console.log(bl.lcaWithSum(3, 6, false))
  // console.log(bl.lcaWithSum(4, 5, false))
  // console.log(bl.lcaWithSum(4, 6, false))

  // console.log(bl.lcaWithSum(3, 5, true))
  // console.log(bl.lcaWithSum(3, 6, true))
  // console.log(bl.lcaWithSum(4, 5, true))
  // console.log(bl.lcaWithSum(4, 6, true))
}
