## useUpdateEffect

useUpdateEffect 用法等同于 useEffect，但是会忽略首次执行，只在依赖更新时执行。

```ts
declare const _default: typeof useEffect | typeof import('react').useLayoutEffect
```

## useUpdateLayoutEffect

## useAsyncEffect

useEffect 支持异步函数。
组件加载时进行`异步的检查`、`中断执行`。

```tsx
import React, { useState } from 'react'

declare function useAsyncEffect(effect: () => AsyncGenerator<void, void, void> | Promise<void>, deps?: DependencyList): void

function mockCheck(): Promise<boolean> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(true)
    }, 3000)
  })
}

export default () => {
  const [pass, setPass] = useState<boolean>()

  useAsyncEffect(async () => {
    setPass(await mockCheck())
  }, [])

  return (
    <div>
      {pass === undefined && 'Checking...'}
      {pass === true && 'Check passed.'}
    </div>
  )
}
```

```tsx
import React, { useState } from 'react'
import { useAsyncEffect } from 'ahooks'

function mockCheck(val: string): Promise<boolean> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(val.length > 0)
    }, 1000)
  })
}

export default () => {
  const [value, setValue] = useState('')
  const [pass, setPass] = useState<boolean>()

  useAsyncEffect(
    async function* () {
      setPass(undefined)
      const result = await mockCheck(value)
      yield // 检查当前副作用（effect）是否还有效，如果已经被清理（比如依赖变化或组件卸载），就停止执行后面的代码。
      setPass(result)
    },
    [value]
  )

  return (
    <div>
      <input
        value={value}
        onChange={e => {
          setValue(e.target.value)
        }}
      />
      <p>
        {pass === null && 'Checking...'}
        {pass === false && 'Check failed.'}
        {pass === true && 'Check passed.'}
      </p>
    </div>
  )
}
```

## useDebounceEffect

为 useEffect 增加防抖的能力。

```tsx
import { useDebounceEffect } from 'ahooks'
import React, { useState } from 'react'

export default () => {
  const [value, setValue] = useState('hello')
  const [records, setRecords] = useState<string[]>([])
  useDebounceEffect(
    () => {
      setRecords(val => [...val, value])
    },
    [value],
    {
      wait: 1000
    }
  )

  return (
    <div>
      <input value={value} onChange={e => setValue(e.target.value)} placeholder="Typed value" style={{ width: 280 }} />
      <p style={{ marginTop: 16 }}>
        <ul>
          {records.map((record, index) => (
            <li key={index}>{record}</li>
          ))}
        </ul>
      </p>
    </div>
  )
}
```

## useDebounceFn

用来处理防抖函数的 Hook。
run 触发执行，cancel 取消执行，flush 立即执行当前的防抖函数。

```ts
type noop = (...args: any[]) => any
declare function useDebounceFn<T extends noop>(
  fn: T,
  options?: DebounceOptions
): {
  run: import('lodash').DebouncedFunc<(...args: Parameters<T>) => ReturnType<T>>
  cancel: () => void
  flush: () => ReturnType<T> | undefined
}
```

```tsx
import { useDebounceFn } from 'ahooks'
import React, { useState } from 'react'

export default () => {
  const [value, setValue] = useState(0)
  const { run } = useDebounceFn(
    () => {
      setValue(value + 1)
    },
    {
      wait: 500
    }
  )

  return (
    <div>
      <p style={{ marginTop: 16 }}> Clicked count: {value} </p>
      <button type="button" onClick={run}>
        Click fast!
      </button>
    </div>
  )
}
```

## useThrottleFn

同上。

## useThrottleEffect

## useDeepCompareEffect

## useDeepCompareLayoutEffect

## useInterval

## useRafInterval

## useTimeout

## useRafTimeout

## useLockFn

## useUpdate
