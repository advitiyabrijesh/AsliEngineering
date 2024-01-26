package main

import "fmt"

func hash(x int) int {
	// Placeholder for now.
	return x
}

func rho(x int) int {
	for i := 63; i >= 0; i-- {
		if (x & (1 << i)) != 0 {
			return i
		}
	}
	return 63
}

func flajoletMartin(ms []int) int {
	b := 0
	for _, x := range ms {
		msb := rho(hash(x))
		fmt.Printf("most significant bit: %d for number %d\n", msb, x)
		b = max(b, msb)
	}
	return 1 << (b + 1)
}

func main() {
	MagicNumber := 0.77351
	stream := []int{1, 2, 3, 4, 5, 5, 6, 7}
	fmt.Printf("Approximate number of elements in the stream : %d\n",
		int(float64(flajoletMartin(stream))/MagicNumber))
}
