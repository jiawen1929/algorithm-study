// https://www.cnblogs.com/wenruo/p/17050995.html
// !可调度，指的是当 trigger 动作触发副作用函数重新执行时，有能力决定副作用函数的执行时机、次数以及方式。
// 在 effect 函数增加选项，可以指定执行副作用函数的调度器。

interface IEffectFn {
  (): void
  deps: Set<() => void>[]
  options?: IEffectOptions
}

interface IEffectOptions {
  scheduler?: (fn: () => void) => void
}

/** 全局变量，用于存储被注册的副作用函数. */
let activeEffect: IEffectFn | undefined

const effectStack: IEffectFn[] = []

/** 用于注册副作用函数. */
function effect(fn: () => void, options: IEffectOptions = {}) {
  const effectFn = () => {
    /** 每次执行副作用函数之前，先清理依赖. */
    cleanup(effectFn)
    activeEffect = effectFn
    effectStack.push(effectFn)
    fn()
    effectStack.pop()
    activeEffect = effectStack.length > 0 ? effectStack[effectStack.length - 1] : undefined
  }

  /** 把 options 挂在effectFn上. */
  effectFn.options = options
  /** 用来存储与该副作用函数关联的依赖集合. */
  effectFn.deps = [] as IEffectFn['deps']

  effectFn()
}

function cleanup(effectFn: IEffectFn): void {
  /** 在每个依赖集合中把该函数删除. */
  for (let i = 0; i < effectFn.deps.length; i++) {
    const depsSet = effectFn.deps[i]
    depsSet.delete(effectFn)
  }
  effectFn.deps.length = 0
}

const bucket = new WeakMap<object, Map<PropertyKey, Set<IEffectFn>>>() // target -> target key -> Set<副作用函数>
const data = { age: 1 }

/** reactive. */
const reactiveData = new Proxy(data, {
  get(target, key: string) {
    track(target, key)
    // @ts-ignore
    return target[key]
  },
  set(target, key, newValue) {
    // @ts-ignore
    target[key] = newValue
    trigger(target, key)
    return true
  }
})

/** 在track中记录deps. */
function track(target: object, key: PropertyKey): void {
  if (!activeEffect) return
  let depsMap = bucket.get(target)
  if (!depsMap) {
    depsMap = new Map()
    bucket.set(target, depsMap)
  }
  let depsSet = depsMap.get(key)
  if (!depsSet) {
    depsSet = new Set()
    depsMap.set(key, depsSet)
  }
  depsSet.add(activeEffect)

  /** 当前副作用函数也记录下关联的依赖. */
  activeEffect.deps.push(depsSet)
}

function trigger(target: object, key: PropertyKey): void {
  const depsMap = bucket.get(target)
  if (!depsMap) return
  const depsSet = depsMap.get(key)
  if (!depsSet) return

  const effectsToRun = new Set<IEffectFn>()
  depsSet.forEach(effect => {
    // 在 trigger 中执行副作用函数的时候，不执行当前正在处理的副作用函数，即 activeEffect
    if (effect !== activeEffect) {
      effectsToRun.add(effect)
    }
  })

  /** 如果一个副作用函数存在调度器 就用调度器执行副作用函数. */
  effectsToRun.forEach(effect => {
    if (effect.options?.scheduler) {
      effect.options.scheduler(effect)
    } else {
      effect()
    }
  })
}

export {}

if (require.main === module) {
  function demo1() {
    effect(
      () => {
        console.log({ age: reactiveData.age })
      },
      { scheduler: fn => setTimeout(fn, 1000) }
    )
    reactiveData.age++
    console.log('结束')
  }

  // !批处理更新
  // 我们也可以指定副作用函数的执行次数，比如我们对同一个变量连续操作了多次，我们只需要对最终的结果执行副作用函数，中间值可以被忽略。
  // 基于调度器实现这个功能
  //
  // !通过 Set 实现去重，防止函数执行多次，通过 isflushing 做标记，执行过程中不会再次执行
  function demo2() {
    const jobQueue = new Set<() => void>()
    const p = Promise.resolve()

    let isFlushing = false
    function flushJob(): void {
      if (isFlushing) return
      isFlushing = true
      p.then(() => {
        jobQueue.forEach(job => job())
      }).finally(() => {
        isFlushing = false
      })
    }

    effect(
      () => {
        console.log('age:', reactiveData.age)
      },
      {
        scheduler(fn) {
          jobQueue.add(fn)
          flushJob()
        }
      }
    )

    reactiveData.age++
    reactiveData.age++
    reactiveData.age++
    reactiveData.age++
    reactiveData.age++
    console.log('结束')
  }

  demo2()
}
