package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1000000007

func mul(a, b int) int {
	ret := a * b
	if ret > MOD {
		ret %= MOD
	}
	return ret
}

func mulMod(a, b, c int) int {
	ret := a * b
	if ret > MOD {
		ret %= MOD
	}
	return ret
}

func add(a, b int) int {
	ret := a + b
	if ret > MOD {
		ret %= MOD
	}
	return ret
}

func power(a, b, c int) int {
	ret := 1
	for b > 0 {
		if b&1 == 1 {
			ret = mulMod(ret, a, c)
		}
		b >>= 1
		a = mulMod(a, a, c)
	}
	return ret
}

func solver() {
	var n, m, k int
	fmt.Scan(&n, &m, &k)

	p := n * (n - 1) / 2
	den := power(power(p%MOD, MOD-2, MOD), k, MOD)
	var res int

	pk := power(p%MOD, k, MOD)
	pk1 := power(p%MOD, k-1, MOD)
	pk2 := power(p%MOD, k-2, MOD)

	s1 := mulMod(pk, k, MOD)
	s1 = add(s1, -mulMod(mul(k, pk1), p-1, MOD))

	inv2 := power(2, MOD-2, MOD)
	s2 := mulMod(mulMod(k, p-1, MOD), pk1, MOD)
	s2 = mulMod(s2, inv2, MOD)
	s2 = add(s2, -mulMod(mul(k, pk), inv2, MOD))

	s3 := mulMod(mulMod(k, k, MOD), pk, MOD)
	s3 = add(s3, mulMod(mulMod(k, p-1, MOD), pk1, MOD))
	s3 = add(s3, mulMod(mulMod(k, k-1, MOD), mul(pk2, mulMod(p-1, p-1, MOD)), MOD))
	s3 = add(s3, -mulMod(mulMod(2*k, k, MOD), mulMod(pk1, p-1, MOD), MOD))
	s3 = mulMod(s3, inv2, MOD)

	fsum := 0
	for j := 1; j <= m; j++ {
		var a, b, f int
		fmt.Scan(&a, &b, &f)
		fsum += f
	}
	fsum %= MOD

	res = add(add(mulMod(fsum, s1, MOD), mulMod(m, s2, MOD)), mulMod(m, s3, MOD))
	res = mulMod(res, den, MOD)

	fmt.Println(res)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	// Uncomment the following line if you are reading from an input file
	// file, _ := os.Open("inp.txt")
	// defer file.Close()
	// scanner := bufio.NewScanner(file)

	var testCases int
	fmt.Scan(&testCases)

	for t := 0; t < testCases; t++ {
		solver()
	}
}
