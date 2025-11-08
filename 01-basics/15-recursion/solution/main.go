package main

func SumRecursive(n int) int {
	if n <= 0 {
		return 0
	}
	return n + SumRecursive(n-1)
}

func main() {}
