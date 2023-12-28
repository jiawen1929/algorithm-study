// 最长公共子序列 LCS模板

package main

import "fmt"

func main() {
	fmt.Println(LongestCommonSubsequence("ba", "ana"))
	fmt.Println(LongestCommonSubsequenceRestore("ba", "ana"))
	fmt.Println(LongestCommonSubsequence("banana", "ana"))
	fmt.Println(LongestCommonSubsequenceRestore("banana", "ana"))
}

type Str = string

// O(n^2) 求两个字符串 s1 和 s2 的最长公共子序列长度.
func LongestCommonSubsequence(s1, s2 Str) int {
	n := len(s2)
	dp := make([]int, n+1)
	for i := 0; i < len(s1); i++ {
		for j := n - 1; j >= 0; j-- {
			if s1[i] == s2[j] {
				dp[j+1] = max(dp[j+1], dp[j]+1)
			}
		}
		for j := 0; j < n; j++ {
			dp[j+1] = max(dp[j+1], dp[j])
		}
	}
	return dp[n]

}

// O(n^2) 求两个字符串 s1 和 s2 的最长公共子序列.
func LongestCommonSubsequenceRestore(s1, s2 Str) [][]int {
	n, m := len(s1), len(s2)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 0; i < n; i++ {
		pre := dp[i]
		cur := append([]int(nil), pre...)
		dp[i+1] = cur
		for j := 0; j < m; j++ {
			cur[j+1] = max(cur[j+1], cur[j])
			if s1[i] == s2[j] {
				cur[j+1] = max(cur[j+1], pre[j]+1)
			}
		}
	}

	var res [][]int
	ptr1, ptr2 := n, m
	for dp[ptr1][ptr2] > 0 {
		if dp[ptr1][ptr2] == dp[ptr1-1][ptr2] {
			ptr1--
		} else if dp[ptr1][ptr2] == dp[ptr1][ptr2-1] {
			ptr2--
		} else {
			ptr1--
			ptr2--
			res = append(res, []int{ptr1, ptr2})
		}
	}

	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
