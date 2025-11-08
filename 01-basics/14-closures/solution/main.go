package main

func Counter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	c := Counter()
	println(c()) // 1
	println(c()) // 2
}
