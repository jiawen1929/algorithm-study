package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/lyndon_factorization
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	L := NewIncrementalLyndonFactorization()
	for _, c := range s {
		L.Add(c)
	}
	res := L.Factorize()
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

type IncrementalLyndonFactorization struct {
	MinSuffixLen []int32
	i, j, k      int32
	nums         []int32
}

func NewIncrementalLyndonFactorization() *IncrementalLyndonFactorization {
	return &IncrementalLyndonFactorization{MinSuffixLen: []int32{0}}
}

func (ilf *IncrementalLyndonFactorization) Add(c int32) int32 {
	ilf.nums = append(ilf.nums, c)
	for ilf.i < int32(len(ilf.nums)) {
		if ilf.k == ilf.i {
			ilf.i++
		} else if ilf.nums[ilf.k] == ilf.nums[ilf.i] {
			ilf.k++
			ilf.i++
		} else if ilf.nums[ilf.k] < ilf.nums[ilf.i] {
			ilf.k = ilf.j
			ilf.i++
		} else {
			ilf.j += (ilf.i - ilf.j) / (ilf.i - ilf.k) * (ilf.i - ilf.k)
			ilf.i, ilf.k = ilf.k, ilf.j
		}
	}
	if (ilf.i-ilf.j)%(ilf.i-ilf.k) == 0 {
		ilf.MinSuffixLen = append(ilf.MinSuffixLen, ilf.i-ilf.k)
	} else {
		ilf.MinSuffixLen = append(ilf.MinSuffixLen, ilf.MinSuffixLen[ilf.k])
	}
	return ilf.MinSuffixLen[ilf.i]
}

func (ilf *IncrementalLyndonFactorization) Factorize() []int32 {
	i := int32(len(ilf.nums))
	var res []int32
	for i > 0 {
		res = append(res, i)
		i -= ilf.MinSuffixLen[i]
	}
	res = append(res, 0)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
