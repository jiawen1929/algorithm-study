// https://nyaannyaan.github.io/library/math/grundy-number.hpp

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var E = newEratosthenesSieve(1e5 + 10)

func main() {
	// https://yukicoder.me/problems/no/103
	// 给定 n 个数nums (nums[i] <= 1e5,n<=100)
	// 两个人交互变化数字,可以将每个数除以他的素因子p或者p^2(如果有的话)
	// 不能继续操作就算输,问先手是否必胜
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	dag := make([][]int, 1e5+10) // 状态i定义为当前数字为i
	for num := 0; num <= 1e5; num++ {
		factors := E.GetPrimeFactors(num)
		for p, c := range factors {
			dag[num] = append(dag[num], num/p)
			if c > 1 {
				dag[num] = append(dag[num], num/(p*p))
			}
		}
	}

	// 母状态可以拆分成多个相互独立的子状态,则母状态的SG数等于各个子状态的SG数的异或
	grundy := GrundyNumber(dag)
	xor := 0
	for _, num := range nums {
		xor ^= grundy[num]
	}
	if xor == 0 {
		fmt.Println("Bob")
	} else {
		fmt.Println("Alice")
	}
}

// dag: 博弈的每个状态组成的有向无环图.
//  返回值: 每个状态的Grundy数.
//  grundy[i] = mex{grundy[j] | j in dag[i]}.
//  - 如果grundy为0,则先手必败,否则先手必胜.
//  - 若一个母状态可以拆分成多个相互独立的子状态，`则母状态的 SG 数等于各个子状态的 SG 数的异或。`
func GrundyNumber(dag [][]int) (grundy []int) {
	order, ok := topoSort(dag)
	if !ok {
		return
	}

	grundy = make([]int, len(dag))
	memo := make([]int, len(dag)+1)
	for j := len(order) - 1; j >= 0; j-- {
		i := order[j]
		if len(dag[i]) == 0 {
			continue
		}
		for _, v := range dag[i] {
			memo[grundy[v]]++
		}
		for memo[grundy[i]] > 0 {
			grundy[i]++
		}
		for _, v := range dag[i] {
			memo[grundy[v]]--
		}
	}

	return
}

func topoSort(dag [][]int) (order []int, ok bool) {
	n := len(dag)
	visited, temp := make([]bool, n), make([]bool, n)
	var dfs func(int) bool
	dfs = func(i int) bool {
		if temp[i] {
			return false
		}
		if !visited[i] {
			temp[i] = true
			for _, v := range dag[i] {
				if !dfs(v) {
					return false
				}
			}
			visited[i] = true
			order = append(order, i)
			temp[i] = false
		}
		return true
	}

	for i := 0; i < n; i++ {
		if !visited[i] {
			if !dfs(i) {
				return nil, false
			}
		}
	}

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return order, true
}

//
//
//
//
// 埃氏筛
type eratosthenesSieve struct {
	minPrime []int
}

func newEratosthenesSieve(maxN int) *eratosthenesSieve {
	minPrime := make([]int, maxN+1)
	for i := range minPrime {
		minPrime[i] = i
	}
	upper := int(math.Sqrt(float64(maxN))) + 1
	for i := 2; i < upper; i++ {
		if minPrime[i] < i {
			continue
		}
		for j := i * i; j <= maxN; j += i {
			if minPrime[j] == j {
				minPrime[j] = i
			}
		}
	}
	return &eratosthenesSieve{minPrime}
}

func (es *eratosthenesSieve) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	return es.minPrime[n] == n
}

func (es *eratosthenesSieve) GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		m := es.minPrime[n]
		res[m]++
		n /= m
	}
	return res
}

func (es *eratosthenesSieve) GetPrimes() []int {
	res := []int{}
	for i, x := range es.minPrime {
		if i >= 2 && i == x {
			res = append(res, x)
		}
	}
	return res
}
