/* eslint-disable @typescript-eslint/no-non-null-assertion */

class OrderedStream {
  private ptr = 1
  private readonly visited = new Map<number, string>()

  // 有 n 个 (id, value) 对，其中 id 是 1 到 n 之间的一个整数，value 是一个字符串。
  // 不存在 id 相同的两个 (id, value) 对。
  insert(idKey: number, value: string): string[] {
    this.visited.set(idKey, value)

    const res: string[] = []
    while (this.visited.has(this.ptr)) {
      res.push(this.visited.get(this.ptr)!)
      this.ptr++
    }

    return res
  }
}

export {}
// 如果流存储有 id = ptr 的 (id, value) 对，则找出从 id = ptr 开始的 最长 id 连续递增序列
// 否则，返回一个空列表
