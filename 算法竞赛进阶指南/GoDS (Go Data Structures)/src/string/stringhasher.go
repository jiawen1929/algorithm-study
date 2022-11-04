package stringhasher

import "fmt"

// StringHasher returns a function that can be used to hash a slice of the string.
// The returned function takes two indices, left and right,
// and returns the hash of the slice [left, right).
//
// It is based on the Rabin-Karp algorithm.
// The hash function is:
//   hash(s[left:right]) = ((s[left]-offset)*base^(right-left-1) +
//   (s[left+1]-offset)*base^(right-left-2) + ... + (s[right-1]-offset)) % mod
// where base is a prime number and mod is a prime number larger than the maximum value of a rune.
// offset is a constant that is subtracted from each rune to make it non-negative.
func StringHasher(s string, mod int, base int, offset int) func(left int, right int) int {
	prePow := make([]int, len(s)+1)
	prePow[0] = 1
	preHash := make([]int, len(s)+1)
	for i, v := range s {
		prePow[i+1] = (prePow[i] * base) % mod
		preHash[i+1] = (preHash[i]*base + int(v) - offset) % mod
	}

	sliceHash := func(left, right int) int {
		if left >= right {
			return 0
		}
		return (preHash[right] - preHash[left]*prePow[right-left]%mod + mod) % mod
	}

	return sliceHash
}

// In order to avoid hash collision, we can use two hash functions.
// Two strings are equal if and only if two hashes are equal.
func BiStringHasher(s string, mod1, mod2, base1, base2, offset1, offset2 int) func(left int, right int) (hash1, hash2 int) {
	hasher1 := StringHasher(s, mod1, base1, offset1)
	hasher2 := StringHasher(s, mod2, base2, offset2)

	sliceHash := func(left, right int) (hash1, hash2 int) {
		if left >= right {
			return 0, 0
		}

		hash1 = hasher1(left, right)
		hash2 = hasher2(left, right)
		return
	}

	return sliceHash
}

func demo() {
	s := "abcabc"
	const MOD1, MOD2 int = 1e8 + 7, 1e9 + 7
	const BASE1, BASE2 int = 131, 13131
	const OFFSET1, OFFSET2 int = 0, 0
	hasher := BiStringHasher(s, MOD1, MOD2, BASE1, BASE2, OFFSET1, OFFSET2)
	fmt.Println(hasher(0, 3))
	fmt.Println(hasher(3, 6))
}
