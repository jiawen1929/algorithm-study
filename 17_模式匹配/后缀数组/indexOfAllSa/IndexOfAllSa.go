// IndexOfAllSa/IndexOfAllMultiString/Lookup/LookupAll
// 在t里lookup一个s，
// 就是求在t的sa数组上求lcp(s,t[i:])>=len(s)的一段区间，
// 二分，check比较字典序做的

package main

import (
	"index/suffixarray"
	"reflect"
	"sort"
	"unsafe"
)

// 返回s在原串中所有匹配的位置(无序).
// O(len(s)*log(n))+len(result).
type LookupAllFunc func(m int32, shorter func(i int32) int32) (result []int32)

// sa: 可选参数.
func UseLookupAll(n int32, longer func(i int32) int32, sa []int32) LookupAllFunc {
	if sa == nil {
		s := make([]byte, 0, n*4)
		for i := int32(0); i < n; i++ {
			v := longer(i)
			s = append(s, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
		}
		_sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))
		sa = make([]int32, 0, n)
		for _, v := range _sa {
			if v&3 == 0 {
				sa = append(sa, v>>2)
			}
		}
	}

	compareSlice32 := func(a, b []int32) int8 {
		n1, n2 := len(a), len(b)
		ptr1, ptr2 := 0, 0
		for ptr1 < n1 && ptr2 < n2 {
			if a[ptr1] < b[ptr2] {
				return -1
			}
			if a[ptr1] > b[ptr2] {
				return 1
			}
			ptr1++
			ptr2++
		}
		if ptr1 == n1 && ptr2 == n2 {
			return 0
		}
		if ptr1 == n1 {
			return -1
		}
		return 1
	}

	hasPrefix := func(s []int32, prefix []int32) bool {
		if len(s) < len(prefix) {
			return false
		}
		for i, v := range prefix {
			if s[i] != v {
				return false
			}
		}
		return true
	}

	longerOrds := make([]int32, n)
	for i := int32(0); i < n; i++ {
		longerOrds[i] = longer(i)
	}
	f := func(m int32, shorter func(i int32) int32) []int32 {
		if n == 0 || m == 0 {
			return nil
		}
		target := make([]int32, m)
		for i := int32(0); i < m; i++ {
			target[i] = shorter(i)
		}
		sa, cur := sa, longerOrds
		i := sort.Search(len(sa), func(i int) bool { return compareSlice32(cur[sa[i]:], target) >= 0 })
		j := i + sort.Search(len(sa)-i, func(j int) bool { return !hasPrefix(cur[sa[i+j]:], target) })
		result := make([]int32, j-i)
		for k := range result {
			result[k] = sa[i+k]
		}
		return result
	}

	return f
}

// 面试题 17.17. 多次搜索
// https://leetcode.cn/problems/multi-search-lcci/description/
func multiSearch(big string, smalls []string) [][]int {
	res := make([][]int, len(smalls))
	f := UseLookupAll(int32(len(big)), func(i int32) int32 { return int32(big[i]) }, nil)
	for i, small := range smalls {
		indexes := f(int32(len(small)), func(i int32) int32 { return int32(small[i]) })
		sort.Slice(indexes, func(i, j int) bool { return indexes[i] < indexes[j] })
		res[i] = make([]int, len(indexes))
		for j, idx := range indexes {
			res[i][j] = int(idx)
		}
	}
	return res
}
