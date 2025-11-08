package main
import "fmt"

func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	fmt.Println("Sum(1,2,3,4):", Sum(1, 2, 3, 4))
	slice := []int{5, 6, 7}
	fmt.Println("Sum(5,6,7):", Sum(slice...))
}
