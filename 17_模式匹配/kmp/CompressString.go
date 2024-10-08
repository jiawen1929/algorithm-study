package main

import "fmt"

func main() {
	fmt.Println(CompressString("abababa", "ababab"))
	fmt.Println(CompressNums([]int{1, 2, 3}, []int{3, 2, 1}))
	fmt.Println(CompressStringNaive("abababa", "ababab"))
	fmt.Println(CompressNumsNaive([]int{1, 2, 3}, []int{3, 2, 1}))
}

const INF int = 1e18

// pre的后缀和post的前缀的最大公共长度
func CompressString(pre, post string) int {
	cat := post + "#" + pre
	next_ := GetNext(len(cat), func(i int) int { return int(cat[i]) })
	return next_[len(cat)-1]
}

// pre的后缀和post的前缀的最大公共长度
func CompressNums(pre, post []int) int {
	newNums := make([]int, len(pre)+len(post)+1)
	copy(newNums, post)
	newNums[len(post)] = INF
	copy(newNums[len(post)+1:], pre)
	next_ := GetNext(len(newNums), func(i int) int { return newNums[i] })
	return next_[len(newNums)-1]
}

// `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度.
func GetNext(n int, f func(i int) int) []int {
	next := make([]int, n)
	j := 0
	for i := 1; i < n; i++ {
		for j > 0 && f(i) != f(j) {
			j = next[j-1]
		}
		if f(i) == f(j) {
			j++
		}
		next[i] = j
	}
	return next
}

func CompressStringNaive(pre, post string) int {
	m, n := len(pre), len(post)
	for res := min(m, n); res >= 0; res-- {
		if pre[m-res:] == post[:res] {
			return res
		}
	}
	return 0
}

func CompressNumsNaive(pre, post []int) int {
	m, n := len(pre), len(post)
	for res := min(m, n); res >= 0; res-- {
		ok := true
		for i := 0; i < res; i++ {
			if pre[m-res+i] != post[i] {
				ok = false
				break
			}
		}
		if ok {
			return res
		}
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
