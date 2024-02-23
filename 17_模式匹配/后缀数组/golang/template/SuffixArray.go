// https://github.com/EndlessCheng/codeforces-go/blob/646deb927bbe089f60fc0e9f43d1729a97399e5f/copypasta/strings.go#L556
// https://visualgo.net/zh/suffixarray
// !常用分隔符 #(35) $(36) _(95) |(124)
// SA-IS 与 DC3 的效率对比 https://riteme.site/blog/2016-6-19/sais.html#5
// 注：Go1.13 开始使用 SA-IS 算法
//
// - 支持sa/rank/lcp
// - 比较任意两个子串的字典序
// - 求出任意两个子串的最长公共前缀(lcp)

//  sa : 排第几的后缀是谁.
//  rank : 每个后缀排第几.
//  lcp : 排名相邻的两个后缀的最长公共前缀.
// 	lcp[0] = 0
// 	lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
//
//  "banana" -> sa: [5 3 1 0 4 2], rank: [3 2 5 1 4 0], lcp: [0 1 3 0 0 2]
//
//  !lcp(sa[i],sa[j]) = min(height[i+1..j])
//
// !api:
//  func NewSuffixArray(ords []int) *SuffixArray
//  func NewSuffixArrayWithString(s string) *SuffixArray
//  func (suf *SuffixArray) Lcp(a, b int, c, d int) int
//  func (suf *SuffixArray) CompareSubstr(a, b int, c, d int) int
//  func (suf *SuffixArray) LcpRange(left int, k int) (start, end int)
//  func GetSA(ords []int) (sa []int)
//  func UseSA(ords []int) (sa, rank, lcp []int)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"math/bits"
)

func main() {
	abc141e()
	// abc213f()
	// abc272f()

	// P3804()

	// NumberofSubstrings()

	// testLcpRange()
}

// G3. Good Substrings
// https://codeforces.com/contest/316/submission/218670841
func CF316() {
	// itoa
}

// https://codeforces.com/contest/126/submission/227749650
func CF126() {
}

// P3804 【模板】后缀自动机（SAM）
// https://www.luogu.com.cn/problem/P3804
// 给定一个长度为 n 的只包含小写字母的字符串 s。
// !对于所有 s 的出现次数不为 1 的子串，设其 value值为该 子串出现的次数 × 该子串的长度。
// 请计算，value 的最大值是多少。
// n <= 1e6
//
// !子串出现次数乘以次数的最大值-直方图最大矩形
// 直方图最大矩形
// 子串长度看成高,lcp范围看成宽
// https://www.acwing.com/solution/content/25201/
func P3804() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	if len(s) <= 1 {
		fmt.Fprintln(out, 0)
		return
	}
	S := NewSuffixArrayFromString(s)
	heights := S.Height
	L, R := GetRange(heights, false, false, false) // 求每个元素作为最小值的影响范围(区间)
	res := 0
	for i := 0; i < len(heights); i++ {
		res = max(res, heights[i]*(R[i]-L[i]+2))
	}
	fmt.Fprintln(out, res)
}

// 解法2 O(nlogn):
//
//	二分答案，对 height 分组，判定组内元素个数不小于 k, 类似 不可重叠最长重复子串 的做法
func P2852() {

}

// 长度不小于k的公共子串个数
// https://blog.nowcoder.net/n/0a4cfff0f0bc424c9a29979dc7d8f586

// 重复次数最多的连续重复子串
// https://blog.nowcoder.net/n/47821f2464e146ea86d83b224a91d855
// https://blog.nowcoder.net/n/f9c3bcdf807546bd9c8d8cc43df84079

// 连续的若干个相同子串
// https://oi-wiki.org/string/sa/#%E8%BF%9E%E7%BB%AD%E7%9A%84%E8%8B%A5%E5%B9%B2%E4%B8%AA%E7%9B%B8%E5%90%8C%E5%AD%90%E4%B8%B2

// F - Common Prefixes-每个后缀与所有后缀的LCP长度和
// https://atcoder.jp/contests/abc213/tasks/abc213_f
// 定义LCP(X,Y)为字符串X,Y的公共前缀长度(LCP)。
// 给定长度为N的字符串S，设S表示从第i个字符开始的S的后缀(就是后缀数组里的那些后缀)。
// !计算出:对于k=1,2,...,N,LCP(Sk,S1)+LCP(Sk,S2)+ +...+LCP(Sk,SN)的值。
// !即求每个后缀与所有后缀的公共前缀长度和。
// n<=1e6
//
// https://blog.hamayanhamayan.com/entry/2021/08/09/010405
func abc213f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	S := NewSuffixArrayFromString(s)
	sa, height := S.Sa, S.Height

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[sa[i]] = n - sa[i]
	}

	// !lcp(sa[i],sa[j]) = min(height[i+1..j])
	clampMaxStack := NewClampableStack(false) // 截断最大值的单调栈
	for i := 0; i < n; i++ {
		clampMaxStack.AddAndClamp(height[i])
		res[sa[i]] += clampMaxStack.Sum() // s[i]与左侧所有后缀的lcp和
	}
	clampMaxStack.Clear()
	for i := n - 1; i >= 0; i-- {
		res[sa[i]] += clampMaxStack.Sum() // s[i]与右侧所有后缀的lcp和
		clampMaxStack.AddAndClamp(height[i])
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// F - Two Strings
// https://atcoder.jp/contests/abc272/tasks/abc272_f
// 给定两个长为n的字符串s和t
// 问s和t的所有轮转的子串中 s的轮转子串有多少个字典序 <= t的轮转子串
//
// 技巧:
// 需要一起比较s和t的所有轮转字串的字典序
// !构造一个新的字符串 s+s+'#'+t+t+'|'
// (注意题目要的是小于等于, 这样保证两个字符串在比较完长度为n后S后面的#小于T中任意一个字符。)
// !后缀数组求出每个串的rank
// !然后在t的rank中 用s的每个子串rank二分出t中的pos
func abc272f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s, t string
	fmt.Fscan(in, &s, &t)

	SMALL, BIG := "#", "|"
	sstt := s + s + SMALL + t + t + BIG
	S := NewSuffixArrayFromString(sstt)
	rank := S.Rank
	sRank, tRank := rank[:n], rank[2*n+1:2*n+1+n]
	sort.Ints(tRank)
	res := 0
	for _, r := range sRank {
		res += n - sort.SearchInts(tRank, r)
	}
	fmt.Fprintln(out, res)
}

// !不同子串长度之和
// 枚举每个后缀，计算前缀总数，再减掉重复
func diffSum(s string) int {
	n := len(s)
	ords := make([]int, n)
	for i, c := range s {
		ords[i] = int(c)
	}
	_, _, height := UseSA(ords)
	res := n * (n + 1) * (n + 2) / 6 // 所有子串长度1到n的平方和
	for _, h := range height {
		res -= h * (h + 1) / 2
	}
	return res
}

// 1044. 最长重复子串(可重叠最长重复子串)
// https://leetcode.cn/problems/longest-duplicate-substring/description/
// 给你一个字符串 s ，考虑其所有 重复子串 ：即 s 的（连续）子串，在 s 中出现 2 次或更多次。这些出现之间可能存在重叠。
// 返回 任意一个 可能具有最长长度的重复子串。如果 s 不含重复子串，那么答案为 "" 。
// 子串就是后缀的前缀
// !高度数组中的最大值对应的就是可重叠最长重复子串
func longestDupSubstring(s string) string {
	S := NewSuffixArrayFromString(s)
	sa, height := S.Sa, S.Height
	saIndex, maxHeight := 0, 0
	for i, h := range height {
		if h > maxHeight {
			saIndex = i
			maxHeight = h
		}
	}
	return s[sa[saIndex] : sa[saIndex]+maxHeight]
}

// https://leetcode.cn/problems/largest-merge-of-two-strings/
// 1754. 构造字典序最大的合并字符串
func largestMerge(word1 string, word2 string) string {
	ords1, ords2 := make([]int, len(word1)), make([]int, len(word2))
	for i, c := range word1 {
		ords1[i] = int(c)
	}
	for i, c := range word2 {
		ords2[i] = int(c)
	}
	S := NewSuffixArray2(ords1, ords2)

	n1, n2 := len(word1), len(word2)
	sb := strings.Builder{}

	i, j := 0, 0
	for i < len(word1) && j < len(word2) {
		if S.CompareSubstr(i, n1, j, n2) == 1 {
			sb.WriteByte(word1[i])
			i++
		} else {
			sb.WriteByte(word2[j])
			j++
		}
	}

	sb.WriteString(word1[i:])
	sb.WriteString(word2[j:])

	return sb.String()
}

// 2261. 含最多 K 个可整除元素的子数组
// https://leetcode.cn/problems/k-divisible-elements-subarrays/
// 找出并返回满足要求的不同的子数组数，要求子数组中最多 k 个可被 p 整除的元素。
func countDistinct(nums []int, k int, p int) (res int) {
	n := len(nums)

	mods := make([]int, n)
	for i := range mods {
		mods[i] = nums[i] % p
	}

	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	// 1. 先用双指针O(n)的时间计算出所有满足条件的子数组的数量 注意要枚举后缀(固定left 移动right)
	right, countK := 0, 0
	suffixLen := make([]int, n) // 记录每个后缀取到的长度
	for left := 0; left < n; left++ {
		for right < n && countK+boolToInt((mods[right] == 0)) <= k {
			countK += boolToInt((mods[right] == 0))
			right++
		}
		res += right - left
		suffixLen[left] = right - left
		countK -= boolToInt(mods[left] == 0)
	}

	// 2. height数组去重
	sa, _, height := UseSA(nums)
	// 计算子串重复数量 按后缀排序的顺序枚举后缀 lcp(height)去重
	for i := 0; i < n-1; i++ {
		suffix1, suffix2 := sa[i], sa[i+1]
		subLen1, subLen2 := suffixLen[suffix1], suffixLen[suffix2]
		res -= min(height[i+1], min(subLen1, subLen2))
	}
	return
}

// https://judge.yosupo.jp/problem/number_of_substrings
// 返回 s 的不同子字符串的个数(本质不同子串数)
// 用所有子串的个数，减去相同子串的个数，就可以得到不同子串的个数。
// !子串就是后缀的前缀 按后缀排序的顺序枚举后缀，每次新增的子串就是除了与上一个后缀的 LCP 剩下的前缀
// !计算后缀数组和高度数组。根据高度数组的定义，所有高度之和就是相同子串的个数。(每一对相同子串在高度数组产生1贡献)
func NumberofSubstrings() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	ords := make([]int, n)
	for i, c := range s {
		ords[i] = int(c)
	}
	res := n * (n + 1) / 2
	_, _, height := UseSA(ords)
	for _, h := range height {
		res -= h
	}
	fmt.Fprintln(out, res)
}

func testLcpRange() {
	n := int(1e3)
	ords := make([]int, n)
	for i := 1; i < n; i++ {
		ords[i] = i * i
		ords[i] ^= ords[i-1]
	}

	S := NewSuffixArray(ords)
	S2 := NewSuffixArray(ords)
	LcpRange2 := func(left, k int) (start, end int) {
		curRank := S2.Rank[left]
		for i := curRank; i >= 0; i-- {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				start = i
			} else {
				break
			}
		}
		for i := curRank; i < n; i++ {
			sa := S2.Sa[i]
			if S2.Lcp(sa, n, left, n) >= k {
				end = i + 1
			} else {
				break
			}
		}
		if start == 0 && end == 0 {
			return -1, -1
		}
		return
	}

	fmt.Println(S.LcpRange(4, 0))
	fmt.Println(LcpRange2(4, 0))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			start1, end1 := S.LcpRange(i, j)
			start2, end2 := LcpRange2(i, j)
			if start1 != start2 || end1 != end2 {
				fmt.Println(i, j, start1, end1, start2, end2)
				panic("")
			}
		}
	}
	fmt.Println("pass")
}

func demo() {
	s := "abca"
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	sa, rank, height := UseSA(ords)
	fmt.Println(sa, rank, height)
}

type SuffixArray struct {
	Sa     []int // 排名第i的后缀是谁.
	Rank   []int // 后缀s[i:]的排名是多少.
	Height []int // 排名相邻的两个后缀的最长公共前缀.Height[0] = 0,Height[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	Ords   []int
	n      int
	minSt  *LinearRMQ // 维护lcp的最小值
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray(ords []int) *SuffixArray {
	ords = append(ords[:0:0], ords...)
	res := &SuffixArray{n: len(ords), Ords: ords}
	sa, rank, lcp := res._useSA(ords)
	res.Sa, res.Rank, res.Height = sa, rank, lcp
	return res
}

func NewSuffixArrayFromString(s string) *SuffixArray {
	ords := make([]int, len(s))
	for i, c := range s {
		ords[i] = int(c)
	}
	return NewSuffixArray(ords)
}

// 求任意两个子串s[a,b)和s[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray) Lcp(a, b int, c, d int) int {
	cand := suf._lcp(a, c)
	return min(cand, min(b-a, d-c))
}

// 比较任意两个子串s[a,b)和s[c,d)的字典序.
//
//	s[a,b) < s[c,d) 返回-1.
//	s[a,b) = s[c,d) 返回0.
//	s[a,b) > s[c,d) 返回1.
func (suf *SuffixArray) CompareSubstr(a, b int, c, d int) int {
	len1, len2 := b-a, d-c
	lcp := suf.Lcp(a, b, c, d)
	if len1 == len2 && lcp >= len1 {
		return 0
	}
	if lcp >= len1 || lcp >= len2 { // 一个是另一个的前缀
		if len1 < len2 {
			return -1
		}
		return 1
	}
	if suf.Rank[a] < suf.Rank[c] {
		return -1
	}
	return 1
}

// 与 s[left:] 的 lcp 大于等于 k 的后缀数组(sa)上的区间.
// 如果不存在,返回(-1,-1).
func (suf *SuffixArray) LcpRange(left int, k int) (start, end int) {
	if k > suf.n-left {
		return -1, -1
	}
	if k == 0 {
		return 0, suf.n
	}
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ(suf.Height)
	}
	i := suf.Rank[left] + 1
	start = suf.minSt.MinLeft(i, func(e int) bool { return e >= k }) - 1 // 向左找
	end = suf.minSt.MaxRight(i, func(e int) bool { return e >= k })      // 向右找
	return
}

func (suf *SuffixArray) Print(sa, ords []int) {
	n := len(ords)
	for _, v := range sa {
		s := make([]string, 0, n-v)
		for i := v; i < n; i++ {
			s = append(s, string(ords[i]))
		}
		fmt.Println(strings.Join(s, ""))
	}
}

// 求任意两个后缀s[i:]和s[j:]的最长公共前缀(lcp).
func (suf *SuffixArray) _lcp(i, j int) int {
	if suf.minSt == nil {
		suf.minSt = NewLinearRMQ(suf.Height)
	}
	if i == j {
		return suf.n - i
	}
	r1, r2 := suf.Rank[i], suf.Rank[j]
	if r1 > r2 {
		r1, r2 = r2, r1
	}
	return suf.minSt.Query(r1+1, r2+1)
}

func (suf *SuffixArray) _getSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}
	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = suf._getSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

func (suf *SuffixArray) _useSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = suf._getSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

type LinearRMQ struct {
	n     int
	nums  []int
	small []int
	large [][]int
}

// n: 序列长度.
// less: 入参为两个索引,返回值表示索引i处的值是否小于索引j处的值.
//
//	消除了泛型.
func NewLinearRMQ(nums []int) *LinearRMQ {
	n := len(nums)
	res := &LinearRMQ{n: n}
	stack := make([]int, 0, 64)
	small := make([]int, 0, n)
	var large [][]int
	large = append(large, make([]int, 0, n>>6))
	for i := 0; i < n; i++ {
		for len(stack) > 0 && nums[stack[len(stack)-1]] > nums[i] {
			stack = stack[:len(stack)-1]
		}
		tmp := 0
		if len(stack) > 0 {
			tmp = small[stack[len(stack)-1]]
		}
		small = append(small, tmp|(1<<(i&63)))
		stack = append(stack, i)
		if (i+1)&63 == 0 {
			large[0] = append(large[0], stack[0])
			stack = stack[:0]
		}
	}

	for i := 1; (i << 1) <= n>>6; i <<= 1 {
		csz := n>>6 + 1 - (i << 1)
		v := make([]int, csz)
		for k := 0; k < csz; k++ {
			back := large[len(large)-1]
			v[k] = res._getMin(back[k], back[k+i])
		}
		large = append(large, v)
	}

	res.small = small
	res.large = large
	return res
}

// 查询区间`[start, end)`中的最小值的索引.
func (rmq *LinearRMQ) Query(start, end int) (minIndex int) {
	if start >= end {
		panic(fmt.Sprintf("start(%d) should be less than end(%d)", start, end))
	}
	end--
	left := start>>6 + 1
	right := end >> 6
	if left < right {
		msb := bits.Len64(uint64(right-left)) - 1
		cache := rmq.large[msb]
		i := (left-1)<<6 + bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63))))
		cand1 := rmq._getMin(i, cache[left])
		j := right<<6 + bits.TrailingZeros64(uint64(rmq.small[end]))
		cand2 := rmq._getMin(cache[right-(1<<msb)], j)
		return rmq._getMin(cand1, cand2)
	}
	if left == right {
		i := (left-1)<<6 + bits.TrailingZeros64(uint64(rmq.small[left<<6-1]&(^0<<(start&63))))
		j := left<<6 + bits.TrailingZeros64(uint64(rmq.small[end]))
		return rmq._getMin(i, j)
	}
	return right<<6 + bits.TrailingZeros64(uint64(rmq.small[end]&(^0<<(start&63))))
}

func (rmq *LinearRMQ) _getMin(i, j int) int {
	if rmq.nums[i] < rmq.nums[j] {
		return i
	}
	return j
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ) MaxRight(left int, check func(e int) bool) int {
	if left == st.n {
		return st.n
	}
	ok, ng := left, st.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(st.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (st *LinearRMQ) MinLeft(right int, check func(e int) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(st.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 用于求解`两个字符串s和t`相关性质的后缀数组.
type SuffixArray2 struct {
	SA     *SuffixArray
	offset int
}

// !ord值很大时,需要先离散化.
// !ords[i]>=0.
func NewSuffixArray2(ords1, ords2 []int) *SuffixArray2 {
	newNums := append(ords1, ords2...)
	sa := NewSuffixArray(newNums)
	return &SuffixArray2{SA: sa, offset: len(ords1)}
}

func NewSuffixArray2FromString(s, t string) *SuffixArray2 {
	ords1 := make([]int, len(s))
	for i, c := range s {
		ords1[i] = int(c)
	}
	ords2 := make([]int, len(t))
	for i, c := range t {
		ords2[i] = int(c)
	}
	return NewSuffixArray2(ords1, ords2)
}

// 求任意两个子串s[a,b)和t[c,d)的最长公共前缀(lcp).
func (suf *SuffixArray2) Lcp(a, b int, c, d int) int {
	return suf.SA.Lcp(a, b, c+suf.offset, d+suf.offset)
}

// 比较任意两个子串s[a,b)和t[c,d)的字典序.
//
//	s[a,b) < t[c,d) 返回-1.
//	s[a,b) = t[c,d) 返回0.
//	s[a,b) > t[c,d) 返回1.
func (suf *SuffixArray2) CompareSubstr(a, b int, c, d int) int {
	return suf.SA.CompareSubstr(a, b, c+suf.offset, d+suf.offset)
}

// !注意内部会修改ords.
//
//	 sa : 排第几的后缀是谁.
//	 rank : 每个后缀排第几.
//	 lcp : 排名相邻的两个后缀的最长公共前缀.
//		lcp[0] = 0
//		lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
func UseSA(ords []int) (sa, rank, lcp []int) {
	n := len(ords)
	sa = GetSA(ords)

	rank = make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	// !高度数组 lcp 也就是排名相邻的两个后缀的最长公共前缀。
	// lcp[0] = 0
	// lcp[i] = LCP(s[sa[i]:], s[sa[i-1]:])
	lcp = make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := sa[rk-1]; i+h < n && j+h < n && ords[i+h] == ords[j+h]; h++ {
			}
		}
		lcp[rk] = h
	}

	return
}

// 注意内部会修改ords.
func GetSA(ords []int) (sa []int) {
	if len(ords) == 0 {
		return []int{}
	}

	mn := mins(ords)
	for i, x := range ords {
		ords[i] = x - mn + 1
	}
	ords = append(ords, 0)
	n := len(ords)
	m := maxs(ords) + 1
	isS := make([]bool, n)
	isLms := make([]bool, n)
	lms := make([]int, 0, n)
	for i := 0; i < n; i++ {
		isS[i] = true
	}
	for i := n - 2; i > -1; i-- {
		if ords[i] == ords[i+1] {
			isS[i] = isS[i+1]
		} else {
			isS[i] = ords[i] < ords[i+1]
		}
	}
	for i := 1; i < n; i++ {
		isLms[i] = !isS[i-1] && isS[i]
	}
	for i := 0; i < n; i++ {
		if isLms[i] {
			lms = append(lms, i)
		}
	}
	bin := make([]int, m)
	for _, x := range ords {
		bin[x]++
	}

	induce := func() []int {
		sa := make([]int, n)
		for i := 0; i < n; i++ {
			sa[i] = -1
		}

		saIdx := make([]int, m)
		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := len(lms) - 1; j > -1; j-- {
			i := lms[j]
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		copy(saIdx, bin)
		s := 0
		for i := 0; i < m; i++ {
			s, saIdx[i] = s+saIdx[i], s
		}
		for j := 0; j < n; j++ {
			i := sa[j] - 1
			if i < 0 || isS[i] {
				continue
			}
			x := ords[i]
			sa[saIdx[x]] = i
			saIdx[x]++
		}

		copy(saIdx, bin)
		for i := 0; i < m-1; i++ {
			saIdx[i+1] += saIdx[i]
		}
		for j := n - 1; j > -1; j-- {
			i := sa[j] - 1
			if i < 0 || !isS[i] {
				continue
			}
			x := ords[i]
			saIdx[x]--
			sa[saIdx[x]] = i
		}

		return sa
	}

	sa = induce()

	lmsIdx := make([]int, 0, len(sa))
	for _, i := range sa {
		if isLms[i] {
			lmsIdx = append(lmsIdx, i)
		}
	}
	l := len(lmsIdx)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = -1
	}
	ord := 0
	order[n-1] = ord
	for i := 0; i < l-1; i++ {
		j, k := lmsIdx[i], lmsIdx[i+1]
		for d := 0; d < n; d++ {
			jIsLms, kIsLms := isLms[j+d], isLms[k+d]
			if ords[j+d] != ords[k+d] || jIsLms != kIsLms {
				ord++
				break
			}
			if d > 0 && (jIsLms || kIsLms) {
				break
			}
		}
		order[k] = ord
	}
	b := make([]int, 0, l)
	for _, i := range order {
		if i >= 0 {
			b = append(b, i)
		}
	}
	var lmsOrder []int
	if ord == l-1 {
		lmsOrder = make([]int, l)
		for i, ord := range b {
			lmsOrder[ord] = i
		}
	} else {
		lmsOrder = GetSA(b)
	}
	buf := make([]int, len(lms))
	for i, j := range lmsOrder {
		buf[i] = lms[j]
	}
	lms = buf
	return induce()[1:]
}

type ClampableStackItem = struct {
	value int
	count int32
}

type ClampableStack struct {
	clampMin bool
	total    int
	count    int
	stack    []ClampableStackItem
}

// clampMin：
//
//	为true时，调用AddAndClamp(x)后，容器内所有数最小值被截断(小于x的数变成x)；
//	为false时，调用AddAndClamp(x)后，容器内所有数最大值被截断(大于x的数变成x).
func NewClampableStack(clampMin bool) *ClampableStack {
	return &ClampableStack{clampMin: clampMin}
}

func (h *ClampableStack) AddAndClamp(x int) {
	newCount := 1
	if h.clampMin {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value > x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	} else {
		for len(h.stack) > 0 {
			top := h.stack[len(h.stack)-1]
			if top.value < x {
				break
			}
			h.stack = h.stack[:len(h.stack)-1]
			v, c := top.value, int(top.count)
			h.total -= v * c
			newCount += c
		}
	}
	h.total += x * newCount
	h.count++
	h.stack = append(h.stack, ClampableStackItem{value: x, count: int32(newCount)})
}

func (h *ClampableStack) Sum() int {
	return h.total
}

func (h *ClampableStack) Len() int {
	return h.count
}

func (h *ClampableStack) Clear() {
	h.stack = h.stack[:0]
	h.total = 0
	h.count = 0
}

// 求每个元素作为最值的影响范围(闭区间).
func GetRange(nums []int, isMax, isLeftStrict, isRightStrict bool) (leftMost, rightMost []int) {
	compareLeft := func(stackValue, curValue int) bool {
		if isLeftStrict && isMax {
			return stackValue <= curValue
		} else if isLeftStrict && !isMax {
			return stackValue >= curValue
		} else if !isLeftStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	compareRight := func(stackValue, curValue int) bool {
		if isRightStrict && isMax {
			return stackValue <= curValue
		} else if isRightStrict && !isMax {
			return stackValue >= curValue
		} else if !isRightStrict && isMax {
			return stackValue < curValue
		} else {
			return stackValue > curValue
		}
	}

	n := len(nums)
	leftMost, rightMost = make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		rightMost[i] = n - 1
	}

	stack := []int{}
	for i := 0; i < n; i++ {
		for len(stack) > 0 && compareRight(nums[stack[len(stack)-1]], nums[i]) {
			rightMost[stack[len(stack)-1]] = i - 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	stack = []int{}
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && compareLeft(nums[stack[len(stack)-1]], nums[i]) {
			leftMost[stack[len(stack)-1]] = i + 1
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	return
}

type MonoQueueValue = int
type MonoQueue struct {
	MinQueue       []MonoQueueValue
	_minQueueCount []int32
	_less          func(a, b MonoQueueValue) bool
	_len           int
}

func NewMonoQueue(less func(a, b MonoQueueValue) bool) *MonoQueue {
	return &MonoQueue{
		_less: less,
	}
}

func (q *MonoQueue) Append(value MonoQueueValue) *MonoQueue {
	count := int32(1)
	for len(q.MinQueue) > 0 && q._less(value, q.MinQueue[len(q.MinQueue)-1]) {
		q.MinQueue = q.MinQueue[:len(q.MinQueue)-1]
		count += q._minQueueCount[len(q._minQueueCount)-1]
		q._minQueueCount = q._minQueueCount[:len(q._minQueueCount)-1]
	}
	q.MinQueue = append(q.MinQueue, value)
	q._minQueueCount = append(q._minQueueCount, count)
	q._len++
	return q
}

func (q *MonoQueue) Popleft() {
	q._minQueueCount[0]--
	if q._minQueueCount[0] == 0 {
		q.MinQueue = q.MinQueue[1:]
		q._minQueueCount = q._minQueueCount[1:]
	}
	q._len--
}

func (q *MonoQueue) Head() MonoQueueValue {
	return q.MinQueue[0]
}

func (q *MonoQueue) Min() MonoQueueValue {
	return q.MinQueue[0]
}

func (q *MonoQueue) Len() int {
	return q._len
}

func (q *MonoQueue) String() string {
	sb := []string{}
	for i := 0; i < len(q.MinQueue); i++ {
		sb = append(sb, fmt.Sprintf("%v", pair{q.MinQueue[i], q._minQueueCount[i]}))
	}
	return fmt.Sprintf("MonoQueue{%v}", strings.Join(sb, ", "))
}

type pair struct {
	value MonoQueueValue
	count int32
}

func (p pair) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
}

func mins(a []int) int {
	mn := a[0]
	for _, x := range a {
		if x < mn {
			mn = x
		}
	}
	return mn
}

func maxs(a []int) int {
	mx := a[0]
	for _, x := range a {
		if x > mx {
			mx = x
		}
	}
	return mx
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b

}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
