// https://nyaannyaan.github.io/library/tree/tree-hash.hpp

/**
 * @param tree 临接表表示的树 0-n-1.根节点为0
 * @param seed 随机种子
 * @returns 每个结点的子树哈希值(子树包括自己).
 * 树的子树顺序不同也会被认为是不同的树.
 */
function treeHash(n: number, tree: number[][], seed: number): BigUint64Array {
  const random = useRandom(seed)
  const bases = new BigUint64Array(n).map(() => BigInt(random.next()))
  const depths = new Uint32Array(n)
  const hashes = new BigUint64Array(n).fill(1n)
  dfs(0, -1)
  return hashes

  function dfs(cur: number, pre: number): number {
    let dep = 0
    for (const next of tree[cur]) {
      if (next === pre) continue
      dep = Math.max(dep, dfs(next, cur) + 1)
    }

    const base = bases[dep]
    for (const next of tree[cur]) {
      if (next === pre) continue
      hashes[cur] *= base + hashes[next]
    }

    depths[cur] = dep
    return dep
  }
}

function useRandom(seed: number) {
  function fastRandom(): number {
    seed ^= seed << 13
    seed ^= seed >>> 17
    seed ^= seed << 5
    return seed >>> 0
  }

  return {
    next: fastRandom
  }
}

export { treeHash }

if (require.main === module) {
  const seed = (Math.floor(Date.now() / 2) + 1) >>> 0
  const tree1 = [[1, 2, 3], [0, 4, 5], [0, 6, 7], [0], [1], [1], [2], [2]]
  const tree2 = [[1, 2, 3], [0, 4, 5], [0, 6, 7], [0], [1], [1], [2], [2]]
  console.log(treeHash(8, tree1, seed))
  console.log(treeHash(8, tree2, seed))
}
