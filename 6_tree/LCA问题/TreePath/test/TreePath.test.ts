import { CompressedBinaryLiftWithSum } from '../../CompressedBinaryLiftWithSum/CompressedBinaryLiftWithSum'
import { ITreePath, TreePath } from '../TreePath'

describe('TreePath.ts', () => {
  let n: number
  let edges: [number, number][]
  let tree: number[][]
  let values: number[]
  let bl: CompressedBinaryLiftWithSum
  const createPath = (from: number, to: number): ITreePath => {
    return new TreePath(from, to, { depth: bl.depth, kthAncestorFn: bl.kthAncestor, lcaFn: bl.lca })
  }

  beforeEach(() => {
    //          0
    //        /   \
    //       1     2
    //      / \     \
    //     3   4     5
    //         /
    //        6
    n = 7
    edges = [
      [0, 1],
      [0, 2],
      [1, 3],
      [1, 4],
      [2, 5],
      [4, 6]
    ]
    tree = Array(n)
    for (let i = 0; i < n; i++) tree[i] = []
    edges.forEach(([u, v]) => {
      tree[u].push(v)
      tree[v].push(u)
    })
    values = [1, 1, 2, 3, 4, 5, 6]
    bl = new CompressedBinaryLiftWithSum(tree, i => values[i], {
      e: () => 0,
      op: (a, b) => a + b
    })
  })

  it('should support kthNodeOnPath', () => {
    const path1 = createPath(3, 6)
    expect(path1.kthNodeOnPath(0)).toBe(3)
    expect(path1.kthNodeOnPath(1)).toBe(1)
    expect(path1.kthNodeOnPath(2)).toBe(4)
    expect(path1.kthNodeOnPath(3)).toBe(6)
    expect(path1.kthNodeOnPath(4)).toBe(-1)
    const path2 = createPath(6, 3)
    expect(path2.kthNodeOnPath(0)).toBe(6)
    expect(path2.kthNodeOnPath(1)).toBe(4)
    expect(path2.kthNodeOnPath(2)).toBe(1)
    expect(path2.kthNodeOnPath(3)).toBe(3)
    expect(path2.kthNodeOnPath(4)).toBe(-1)
    const path3 = createPath(3, 3)
    expect(path3.kthNodeOnPath(0)).toBe(3)
    expect(path3.kthNodeOnPath(1)).toBe(-1)
    const path4 = createPath(5, 6)
    expect(path4.kthNodeOnPath(0)).toBe(5)
    expect(path4.kthNodeOnPath(1)).toBe(2)
    expect(path4.kthNodeOnPath(2)).toBe(0)
    expect(path4.kthNodeOnPath(3)).toBe(1)
    expect(path4.kthNodeOnPath(4)).toBe(4)
    expect(path4.kthNodeOnPath(5)).toBe(6)
    expect(path4.kthNodeOnPath(6)).toBe(-1)
  })

  it('should support onPath', () => {
    const path1 = createPath(3, 6)
    expect(path1.onPath(3)).toBeTruthy()
    expect(path1.onPath(1)).toBeTruthy()
    expect(path1.onPath(4)).toBeTruthy()
    expect(path1.onPath(6)).toBeTruthy()
    expect(path1.onPath(0)).toBeFalsy()
    expect(path1.onPath(2)).toBeFalsy()
    expect(path1.onPath(5)).toBeFalsy()

    const path2 = createPath(6, 3)
    expect(path2.onPath(6)).toBeTruthy()
    expect(path2.onPath(4)).toBeTruthy()
    expect(path2.onPath(1)).toBeTruthy()
    expect(path2.onPath(3)).toBeTruthy()
    expect(path2.onPath(0)).toBeFalsy()
    expect(path2.onPath(2)).toBeFalsy()
    expect(path2.onPath(5)).toBeFalsy()

    const path3 = createPath(3, 3)
    expect(path3.onPath(3)).toBeTruthy()

    const path4 = createPath(5, 6)
    expect(path4.onPath(5)).toBeTruthy()
    expect(path4.onPath(2)).toBeTruthy()
    expect(path4.onPath(0)).toBeTruthy()
    expect(path4.onPath(1)).toBeTruthy()
    expect(path4.onPath(4)).toBeTruthy()
    expect(path4.onPath(6)).toBeTruthy()
    expect(path4.onPath(3)).toBeFalsy()
  })

  it('should support hasIntersection', () => {
    expect(createPath(3, 5).hasIntersection(createPath(1, 6))).toBeTruthy()
    expect(createPath(0, 5).hasIntersection(createPath(1, 6))).toBeFalsy()
  })

  it('should support getIntersection', () => {
    const res1 = createPath(3, 5).getIntersection(createPath(1, 6))
    expect(res1).toEqual({ p1: 1, p2: 1 })
    const res2 = createPath(3, 6).getIntersection(createPath(1, 4))
    expect(res2).toEqual({ p1: 4, p2: 1 })
    const res3 = createPath(0, 5).getIntersection(createPath(1, 6))
    expect(res3).toBeUndefined()
  })

  it('should support countIntersection', () => {
    expect(createPath(3, 5).countIntersection(createPath(1, 6))).toBe(1)
    expect(createPath(0, 5).countIntersection(createPath(1, 6))).toBe(0)
    expect(createPath(3, 6).countIntersection(createPath(1, 4))).toBe(2)
    expect(createPath(3, 3).countIntersection(createPath(3, 3))).toBe(1)
    expect(createPath(5, 6).countIntersection(createPath(4, 2))).toBe(4)
  })

  // it('debug', () => {
  //   const path3_6 = createPath(3, 6)
  //   console.log(path3_6.split(4))
  // })

  // 6种情况:
  // down和to在一条链上，此时separator为:down/to/非down非to
  // down和to不在一条链上，此时separator为:down/to/非down非to
  it('should support split', () => {
    const p3_6 = createPath(3, 6)
    const { path1: pathx_x_3_6_3, path2: path1_6 } = p3_6.split(3)
    expect(pathx_x_3_6_3).toBeUndefined()
    expect(path1_6?.from).toBe(1)
    expect(path1_6?.to).toBe(6)
    const { path1: path3_3, path2: path4_6 } = p3_6.split(1)
    expect(path3_3?.from).toBe(3)
    expect(path3_3?.to).toBe(3)
    expect(path4_6?.from).toBe(4)
    expect(path4_6?.to).toBe(6)
    const { path1: path3_1, path2: path6_6 } = p3_6.split(4)
    expect(path3_1?.from).toBe(3)
    expect(path3_1?.to).toBe(1)
    expect(path6_6?.from).toBe(6)
    expect(path6_6?.to).toBe(6)
    const { path1: path3_4, path2: pathx_x_3_6_6 } = p3_6.split(6)
    expect(path3_4?.from).toBe(3)
    expect(path3_4?.to).toBe(4)
    expect(pathx_x_3_6_6).toBeUndefined()

    const p5_3 = createPath(5, 3)
    const { path1: path5_3_5, path2: path2_3 } = p5_3.split(5)
    expect(path5_3_5).toBeUndefined()
    expect(path2_3?.from).toBe(2)
    expect(path2_3?.to).toBe(3)
    const { path1: path5_5, path2: path0_3 } = p5_3.split(2)
    expect(path5_5?.from).toBe(5)
    expect(path5_5?.to).toBe(5)
    expect(path0_3?.from).toBe(0)
    expect(path0_3?.to).toBe(3)
    const { path1: path5_2, path2: path1_3 } = p5_3.split(0)
    expect(path5_2?.from).toBe(5)
    expect(path5_2?.to).toBe(2)
    expect(path1_3?.from).toBe(1)
    expect(path1_3?.to).toBe(3)
    const { path1: path5_0, path2: path_5_3_3_3 } = p5_3.split(1)
    expect(path5_0?.from).toBe(5)
    expect(path5_0?.to).toBe(0)
    expect(path_5_3_3_3?.from).toBe(3)
    expect(path_5_3_3_3?.to).toBe(3)
    const { path1: path5_1, path2: path_5_3_x_x } = p5_3.split(3)
    expect(path5_1?.from).toBe(5)
    expect(path5_1?.to).toBe(1)
    expect(path_5_3_x_x).toBeUndefined()

    //! path命名不够用了，命名空间限制
    const check = (
      path: ITreePath,
      separator: number,
      from1: number | undefined,
      to1: number | undefined,
      from2: number | undefined,
      to2: number | undefined
    ): boolean => {
      const { path1, path2 } = path.split(separator)
      const autual: (number | undefined)[] = [path1?.from, path1?.to, path2?.from, path2?.to]
      const expected: (number | undefined)[] = [from1, to1, from2, to2]
      return autual.every((v, i) => v === expected[i])
    }

    const p3_3 = createPath(3, 3)
    expect(check(p3_3, 3, undefined, undefined, undefined, undefined)).toBeTruthy()

    const p6_0 = createPath(6, 0)
    expect(check(p6_0, 6, undefined, undefined, 4, 0)).toBeTruthy()
    expect(check(p6_0, 4, 6, 6, 1, 0)).toBeTruthy()
    expect(check(p6_0, 1, 6, 4, 0, 0)).toBeTruthy()
    expect(check(p6_0, 0, 6, 1, undefined, undefined)).toBeTruthy()

    const p1_6 = createPath(1, 6)
    expect(check(p1_6, 1, undefined, undefined, 4, 6)).toBeTruthy()
    expect(check(p1_6, 4, 1, 1, 6, 6)).toBeTruthy()
    expect(check(p1_6, 6, 1, 4, undefined, undefined)).toBeTruthy()
  })
})
