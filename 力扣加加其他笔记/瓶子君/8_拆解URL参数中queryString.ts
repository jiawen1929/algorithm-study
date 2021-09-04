import { URL } from 'url'

const url = 'http://sample.com/?a=1&b=2&c=xx&d=2#hash'
const queryString1 = (str: string) => {
  const obj = Object.create(null)
  str.replace(/([^?&=]+)=([^&]+)/g, (_, k, v) => (obj[k] = v))
  return obj
}
const queryString2 = (str: string) => {
  const s = new URL(str).searchParams
  return Object.fromEntries(s.entries())
}

console.log(queryString1(url))
console.log(queryString2(url))

export {}
