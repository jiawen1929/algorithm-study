package main

import (
	"bufio"
	"fmt"
	"os"
)

// 超快读
// 选择 4KB 作为缓存块大小的原因 https://stackoverflow.com/questions/6578394/whats-so-special-about-4kb-for-a-buffer-length
// 4K 是磁盘驱动器上的集群大小(文件的最小分配单位)
// !如果文件仅包含 1 个字节，它将消耗 4K 的物理磁盘空间。而 5K 的文件将导致 8K 分配
func main() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	NextByte := func() byte {
		b := rc()
		for ; '0' > b; b = rc() {
		}
		return b
	}
	_ = NextByte

	// 读一个整数，支持负数
	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt

	// 读一个仅包含小写字母的字符串
	NextString := func() (s []byte) {
		b := rc()
		for ; 'a' > b || b > 'z'; b = rc() { // 'A' 'Z'
		}
		for ; 'a' <= b && b <= 'z'; b = rc() { // 'A' 'Z'
			s = append(s, b)
		}
		return
	}
	_ = NextString

	// 读一个长度为 n 的仅包含小写字母的字符串
	NextStringN := func(n int) []byte {
		b := rc()
		for ; 'a' > b || b > 'z'; b = rc() { // 'A' 'Z'
		}
		s := make([]byte, 0, n)
		s = append(s, b)
		for i := 1; i < n; i++ {
			s = append(s, rc())
		}
		return s
	}
	_ = NextStringN

	n, q := NextByte(), NextByte()
	fmt.Fprintln(out, n, q)
}
