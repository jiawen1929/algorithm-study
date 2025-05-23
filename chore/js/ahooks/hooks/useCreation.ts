import type { DependencyList } from 'react'
import { useRef } from 'react'

function depsAreSame(oldDeps: DependencyList, deps: DependencyList): boolean {
  if (oldDeps === deps) return true
  for (let i = 0; i < oldDeps.length; i++) {
    if (!Object.is(oldDeps[i], deps[i])) return false
  }
  return true
}

/**
 * useCreation 是 useMemo 或 useRef 的替代品。
 *
 *  https://ahooks.js.org/zh-CN/hooks/use-creation
 */
export function useCreation<T>(factory: () => T, deps: DependencyList) {
  const { current } = useRef({
    deps,
    obj: undefined as undefined | T,
    initialized: false
  })
  if (current.initialized === false || !depsAreSame(current.deps, deps)) {
    current.deps = deps
    current.obj = factory()
    current.initialized = true
  }
  return current.obj as T
}
