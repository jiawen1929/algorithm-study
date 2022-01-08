class RBTreeNode<T = number> {
  data: T
  color: number
  left?: RBTreeNode<T>
  right?: RBTreeNode<T>
  parent?: RBTreeNode<T>

  constructor(data: T) {
    this.data = data
    this.left = this.right = this.parent = undefined
    this.color = 0
  }

  sibling(): RBTreeNode<T> | undefined {
    if (!this.parent) return undefined // sibling null if no parent
    return this.isOnLeft() ? this.parent.right : this.parent.left
  }

  hasRedChild(): boolean {
    return Boolean(this.left?.color === 0) || Boolean(this.right?.color === 0)
  }

  isOnLeft(): boolean {
    return this === this.parent?.left
  }
}

type CompareFunction<T = number> = (a: T, b: T) => boolean

class RBTree<T = number> {
  root?: RBTreeNode<T>
  private compare: CompareFunction<T>
  static defaultCompare = (a: any, b: any) => a < b

  constructor(compare = RBTree.defaultCompare) {
    this.root = undefined
    this.compare = compare
  }

  insert(data: T): boolean {
    const node = new RBTreeNode(data)
    const parent = this.search(data)
    if (!parent) this.root = node
    else if (this.compare(node.data, parent.data)) parent.left = node
    else if (this.compare(parent.data, node.data)) parent.right = node
    else return false
    node.parent = parent
    this.fixAfterInsert(node)
    return true
  }

  find(data: T): RBTreeNode<T> | undefined {
    const node = this.search(data)
    return node && node.data === data ? node : undefined
  }

  deleteByValue(val: T): boolean {
    const node = this.search(val)
    if (node?.data !== val) return false
    this.deleteNode(node)
    return true
  }

  *inOrder(root = this.root): Generator<T, unknown, unknown> {
    if (!root) return
    yield* this.inOrder(root.left)
    yield root.data
    yield* this.inOrder(root.right)
  }

  private rotateLeft(pt: RBTreeNode<T>): void {
    const right = pt.right
    pt.right = right?.left
    if (pt.right) pt.right.parent = pt
    right!.parent = pt.parent
    if (!pt.parent) this.root = right
    else if (pt === pt.parent.left) pt.parent.left = right
    else pt.parent.right = right
    right!.left = pt
    pt.parent = right
  }

  private rotateRight(pt: RBTreeNode<T>): void {
    const left = pt.left
    pt.left = left?.right
    if (pt.left) pt.left.parent = pt
    left!.parent = pt.parent
    if (!pt.parent) this.root = left
    else if (pt === pt.parent.left) pt.parent.left = left
    else pt.parent.right = left
    left!.right = pt
    pt.parent = left
  }

  private swapColor(p1: RBTreeNode<T>, p2: RBTreeNode<T>): void {
    const tmp = p1.color
    p1.color = p2.color
    p2.color = tmp
  }

  private swapData(p1: RBTreeNode<T>, p2: RBTreeNode<T>): void {
    const tmp = p1.data
    p1.data = p2.data
    p2.data = tmp
  }

  private fixAfterInsert(pt: RBTreeNode<T>): void {
    let parent: RBTreeNode<T> | undefined = undefined
    let grandParent: RBTreeNode<T> | undefined = undefined
    while (pt !== this.root && pt.color !== 1 && pt.parent?.color === 0) {
      parent = pt.parent
      grandParent = pt.parent.parent
      /*  Case : A
                    Parent of pt is left child of Grand-parent of pt */
      if (parent === grandParent?.left) {
        const uncle = grandParent.right
        /* Case : 1
                         The uncle of pt is also red
                         Only Recoloring required */
        if (uncle && uncle.color === 0) {
          grandParent.color = 0
          parent.color = 1
          uncle.color = 1
          pt = grandParent
        } else {
          /* Case : 2
                               pt is right child of its parent
                               Left-rotation required */
          if (pt === parent.right) {
            this.rotateLeft(parent)
            pt = parent
            parent = pt.parent
          }
          /* Case : 3
                               pt is left child of its parent
                               Right-rotation required */
          this.rotateRight(grandParent)
          this.swapColor(parent!, grandParent)
          pt = parent!
        }
      } else {
        /* Case : B
                   Parent of pt is right child of Grand-parent of pt */
        const uncle = grandParent!.left
        /*  Case : 1
                          The uncle of pt is also red
                          Only Recoloring required */
        if (uncle != undefined && uncle.color === 0) {
          grandParent!.color = 0
          parent.color = 1
          uncle.color = 1
          pt = grandParent!
        } else {
          /* Case : 2
                               pt is left child of its parent
                               Right-rotation required */
          if (pt === parent.left) {
            this.rotateRight(parent)
            pt = parent
            parent = pt.parent
          }
          /* Case : 3
                               pt is right child of its parent
                               Left-rotation required */
          this.rotateLeft(grandParent!)
          this.swapColor(parent!, grandParent!)
          pt = parent!
        }
      }
    }
    this.root!.color = 1
  }

  // searches for given value
  // if found returns the node (used for delete)
  // else returns the last node while traversing (used in insert)
  private search(val: T): RBTreeNode<T> | undefined {
    let p = this.root
    while (p) {
      if (this.compare(val, p.data)) {
        if (!p.left) break
        else p = p.left
      } else if (this.compare(p.data, val)) {
        if (!p.right) break
        else p = p.right
      } else break
    }
    return p
  }

  private deleteNode(v: RBTreeNode<T>): void {
    const u = BSTreplace(v)
    // True when u and v are both black
    const uvBlack = (u == undefined || u.color === 1) && v.color === 1
    const parent = v.parent
    if (!u) {
      // u is null therefore v is leaf
      if (v === this.root) this.root = undefined
      // v is root, making root null
      else {
        if (uvBlack) {
          // u and v both black
          // v is leaf, fix double black at v
          this.fixDoubleBlack(v)
        } else {
          // u or v is red
          if (v.sibling()) {
            // sibling is not null, make it red"
            v.sibling()!.color = 0
          }
        }
        // delete v from the tree
        if (v.isOnLeft()) parent!.left = undefined
        else parent!.right = undefined
      }
      return
    }
    if (!v.left || !v.right) {
      // v has 1 child
      if (v === this.root) {
        // v is root, assign the value of u to v, and delete u
        v.data = u.data
        v.left = v.right = undefined
      } else {
        // Detach v from tree and move u up
        if (v.isOnLeft()) parent!.left = u
        else parent!.right = u
        u.parent = parent
        if (uvBlack) this.fixDoubleBlack(u)
        // u and v both black, fix double black at u
        else u.color = 1 // u or v red, color u black
      }
      return
    }
    // v has 2 children, swap data with successor and recurse
    this.swapData(u, v)
    this.deleteNode(u)
    // find node that replaces a deleted node in BST
    function BSTreplace(x: RBTreeNode<T>) {
      // when node have 2 children
      if (x.left && x.right) return successor(x.right)
      // when leaf
      if (!x.left && !x.right) return null
      // when single child
      return x.left ?? x.right
    }
    // find node that do not have a left child
    // in the subtree of the given node
    function successor(x: RBTreeNode<T>) {
      let temp = x
      while (temp.left) temp = temp.left
      return temp
    }
  }

  private fixDoubleBlack(x: RBTreeNode<T>): void {
    if (x === this.root) return // Reached root
    const sibling = x.sibling()
    const parent = x.parent as RBTreeNode<T>
    if (!sibling) {
      // No sibiling, double black pushed up
      this.fixDoubleBlack(parent)
    } else {
      if (sibling.color === 0) {
        // Sibling red
        parent!.color = 0
        sibling.color = 1
        if (sibling.isOnLeft()) this.rotateRight(parent)
        // left case
        else this.rotateLeft(parent) // right case
        this.fixDoubleBlack(x)
      } else {
        // Sibling black
        if (sibling.hasRedChild()) {
          // at least 1 red children
          if (sibling.left && sibling.left.color === 0) {
            if (sibling.isOnLeft()) {
              // left left
              sibling.left.color = sibling.color
              sibling.color = parent.color
              this.rotateRight(parent)
            } else {
              // right left
              sibling.left.color = parent.color
              this.rotateRight(sibling)
              this.rotateLeft(parent)
            }
          } else {
            if (sibling.isOnLeft()) {
              // left right
              sibling.right!.color = parent.color
              this.rotateLeft(sibling)
              this.rotateRight(parent)
            } else {
              // right right
              sibling.right!.color = sibling.color
              sibling.color = parent.color
              this.rotateLeft(parent)
            }
          }
          parent.color = 1
        } else {
          // 2 black children
          sibling.color = 0
          if (parent.color === 1) this.fixDoubleBlack(parent)
          else parent.color = 1
        }
      }
    }
  }
}

/**
 * @description C++ 里的set
 */
class TreeSet<T = number> {
  private _size: number
  private tree: RBTree<T>
  private compare: CompareFunction<T>

  constructor(collection: Iterable<T> = [], compare = RBTree.defaultCompare) {
    this._size = 0
    this.tree = new RBTree(compare)
    this.compare = compare
    for (const val of collection) this.add(val)
  }

  get size() {
    return this._size
  }

  has(val: T): boolean {
    return !!this.tree.find(val)
  }

  add(val: T): boolean {
    const added = this.tree.insert(val)
    this._size += added ? 1 : 0
    return added
  }

  delete(val: T): boolean {
    const deleted = this.tree.deleteByValue(val)
    this._size -= deleted ? 1 : 0
    return deleted
  }

  ceiling(val: T): T | undefined {
    let p = this.tree.root
    let higher = undefined
    while (p) {
      if (!this.compare(p.data, val)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  floor(val: T): T | undefined {
    let p = this.tree.root
    let lower = undefined
    while (p) {
      if (!this.compare(val, p.data)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  higher(val: T): T | undefined {
    let p = this.tree.root
    let higher = undefined
    while (p) {
      if (this.compare(val, p.data)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  lower(val: T): T | undefined {
    let p = this.tree.root
    let lower = undefined
    while (p) {
      if (this.compare(p.data, val)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  *[Symbol.iterator](): Generator<T, void, unknown> {
    yield* this.values()
  }

  *keys(): Generator<T, void, unknown> {
    yield* this.values()
  }

  *values(): Generator<T, void, unknown> {
    yield* this.tree.inOrder()
  }
}

/**
 * @description C++里的multiset
 */
class TreeMultiSet<T = number> {
  private _size: number
  private tree: RBTree<T>
  private counts: Map<T, number>
  private compare: CompareFunction<T>

  constructor(collection: Iterable<T> = [], compare = RBTree.defaultCompare) {
    this._size = 0
    this.tree = new RBTree(compare)
    this.counts = new Map()
    this.compare = compare
    for (const val of collection) this.add(val)
  }

  get size() {
    return this._size
  }

  has(val: T): boolean {
    return !!this.tree.find(val)
  }

  add(val: T): void {
    this.tree.insert(val)
    this.increase(val)
    this._size++
  }

  delete(val: T): void {
    this.decrease(val)
    if (this.count(val) === 0) {
      this.tree.deleteByValue(val)
    }
    this._size--
  }

  count(val: T): number {
    return this.counts.get(val) ?? 0
  }

  ceiling(val: T): T | undefined {
    let p = this.tree.root
    let higher = null
    while (p) {
      if (!this.compare(p.data, val)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  floor(val: T): T | undefined {
    let p = this.tree.root
    let lower = null
    while (p) {
      if (!this.compare(val, p.data)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  higher(val: T): T | undefined {
    let p = this.tree.root
    let higher = null
    while (p) {
      if (this.compare(val, p.data)) {
        higher = p
        p = p.left
      } else {
        p = p.right
      }
    }
    return higher?.data
  }

  lower(val: T): T | undefined {
    let p = this.tree.root
    let lower = null
    while (p) {
      if (this.compare(p.data, val)) {
        lower = p
        p = p.right
      } else {
        p = p.left
      }
    }
    return lower?.data
  }

  *keys(): Generator<T, void, unknown> {
    yield* this.values()
  }

  *values(): Generator<T, void, unknown> {
    for (const val of this.tree.inOrder()) {
      let count = this.count(val)
      while (count--) yield val
    }
  }

  private decrease(val: T): void {
    this.counts.set(val, this.count(val) - 1)
  }

  private increase(val: T): void {
    this.counts.set(val, this.count(val) + 1)
  }
}

if (require.main === module) {
  const treeSet = new TreeSet()
  treeSet.add(1)
  treeSet.add(2)
  console.log(treeSet.size)
  console.log(treeSet.has(1))
  console.log(treeSet.has(2))
}

export { TreeSet, TreeMultiSet }
