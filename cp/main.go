package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

// Time Complexity O(N * N)
func solve() {
	var n, k, d, ans int
	fmt.Fscan(reader, &n, &k, &d)
	var A = make([]int, n)
	var V = make([]int, k)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i])
		if A[i] == i+1 {
			ans++
		}
	}
	ans += (d - 1) / 2
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &V[i])
	}
	for idx, cnt := 0, 0; cnt < 2*n && d > 1; cnt, d, idx = cnt+1, d-1, (idx+1)%k {
		var val int
		for i := 0; i < V[idx]; i++ {
			A[i]++
		}
		for i := 0; i < n; i++ {
			if A[i] == i+1 {
				val++
			}
		}
		val += (d - 2) / 2
		if val > ans {
			ans = val
		}
	}
	fmt.Fprintln(writer, ans)
}

func main() {
	defer writer.Flush()
	var tests int
	fmt.Fscan(reader, &tests)
	for ; tests > 0; tests-- {
		solve()
	}
}
